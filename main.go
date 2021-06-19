package main

/*
#cgo CFLAGS: -Ilib/SDL2/include -Ilib/lvgl
#cgo LDFLAGS: -L./lib/SDL2/lib/x64 -lSDL2 -L./lib/lvgl/build -l lvgl
#include "mcujs.h"
#include "lvgl.h"
*/
import "C"

import (
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"time"
	"unsafe"

	"github.com/dop251/goja"
)

type JSCallback func(*goja.Runtime)

var MainLoopChan chan JSCallback = make(chan JSCallback, 1024)
var ticker *time.Ticker = time.NewTicker(100 * time.Millisecond)
var vm *goja.Runtime

func print(call goja.FunctionCall) goja.Value {
	s := make([]interface{}, len(call.Arguments))
	for i, v := range call.Arguments {
		s[i] = v.String()
	}
	fmt.Println(s...)
	return nil
}

func runFile(path string) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(path, err)
	}
	_, err = vm.RunScript(path, string(buf))
	if err != nil {
		log.Fatal(path, err)
	}
}

func CallInJSLoop(f JSCallback) {
	MainLoopChan <- f
}

func main() {
	runtime.LockOSThread()

	C.mjsInit()

	vm = goja.New()
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())
	vm.Set("print", print)
	HttpInit()
	LvglInit()
	runFile("js/underscore.js")
	runFile("js/devserver.js")
	runFile("js/boot.js")
	go PromptLoop()

	for {
		select {
		case f := <-MainLoopChan:
			f(vm)
		case <-ticker.C:
			//fmt.Println("tick")
			if C.mjsEventLoop() == 1 {
				log.Println("SDL_QUIT received, exiting...")
				return
			}
		}

	}
}

var lvObjMap map[*C.lv_obj_t]LvObjBase = make(map[*C.lv_obj_t]LvObjBase)
var lvEventCodeMap map[string]int = map[string]int{
	"click": C.LV_EVENT_CLICKED,
}

func ptrToLvObj(ptr *C.lv_obj_t) LvObjBase {
	v, ok := lvObjMap[ptr]
	if ok {
		return v
	}
	ret := &LvObj{
		ptr: ptr,
	}
	lvObjMap[ptr] = ret
	return ret
}

type LvObjBase interface {
	getPtr() *C.lv_obj_t
	getEventHandler(int) func()
}

type LvObj struct {
	ptr         *C.lv_obj_t
	evtHandlers map[int]func()
}

func (l *LvObj) getPtr() *C.lv_obj_t {
	return l.ptr
}

func (l *LvObj) getEventHandler(eCode int) func() {
	if l.evtHandlers == nil {
		return nil
	}
	v, _ := l.evtHandlers[eCode]
	return v
}

type LvLabel struct {
	LvObj
}

type LvBtn struct {
	LvObj
}

func (l *LvObj) On(etype string, callback func()) {
	eCode, ok := lvEventCodeMap[etype]
	if !ok {
		log.Println("Invalid etype:", etype)
		return
	}
	if l.evtHandlers == nil {
		l.evtHandlers = make(map[int]func())
	}
	_, ok = l.evtHandlers[eCode]
	if ok {
		log.Println("Already binded handler:", etype)
		return
	}
	l.evtHandlers[eCode] = callback
	C.mjsSetLvCallback(l.ptr, C.int(eCode))
}

func (l *LvObj) Set_pos(x int, y int) {
	C.lv_obj_set_pos(l.ptr, C.short(x), C.short(y))
}

func (l *LvLabel) Set_text(t string) {
	cstr := C.CString(t)
	C.lv_label_set_text(l.ptr, cstr)
	C.free(unsafe.Pointer(cstr))
}

func LvglInit() {
	obj := vm.NewObject()
	obj.Set("obj", func(p LvObjBase) *LvObj {
		ret := &LvObj{
			ptr: C.lv_obj_create(p.getPtr()),
		}
		lvObjMap[ret.ptr] = ret
		return ret
	})
	obj.Set("label", func(p LvObjBase) *LvLabel {
		ret := &LvLabel{
			LvObj: LvObj{
				ptr: C.lv_label_create(p.getPtr()),
			},
		}
		lvObjMap[ret.ptr] = ret
		return ret
	})
	obj.Set("btn", func(p LvObjBase) *LvBtn {
		ret := &LvBtn{
			LvObj: LvObj{
				ptr: C.lv_btn_create(p.getPtr()),
			},
		}
		lvObjMap[ret.ptr] = ret
		return ret
	})
	obj.Set("scr_act", func() LvObjBase {
		return ptrToLvObj(C.lv_scr_act())
	})
	vm.Set("lv", obj)
}

//export lvEventCallback
func lvEventCallback(evt *C.lv_event_t) {
	fp := ptrToLvObj(evt.target).getEventHandler(int(evt.code))
	if fp != nil {
		fp()
	}
}
