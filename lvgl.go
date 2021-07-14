package main

/*
#include "mcujs.h"
*/
import "C"
import (
	"log"
	"unsafe"
)

type LvObjBase interface {
	getPtr() *C.lv_obj_t
	getEventHandler(int) func()
}

type LvObj struct {
	ptr         *C.lv_obj_t
	evtHandlers map[int]func()
}
type LvLabel struct {
	LvObj
}

type LvBtn struct {
	LvObj
}

var lvObjMap map[*C.lv_obj_t]LvObjBase = make(map[*C.lv_obj_t]LvObjBase)
var lvEventCodeMap map[string]int = map[string]int{
	"click": C.LV_EVENT_CLICKED,
}

//export lvEventCallback
func lvEventCallback(evt *C.lv_event_t) {
	fp := ptrToLvObj(evt.target).getEventHandler(int(evt.code))
	if fp != nil {
		fp()
	}
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
