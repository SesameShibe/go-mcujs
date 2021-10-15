package main

/*
#cgo CFLAGS: -Ilib/SDL2/include -Ilib/lvgl
#cgo LDFLAGS: -L./lib/SDL2/lib/x64 -lSDL2 -L./lib/lvgl -l lvgl
#include "mcujs.h"
*/
import "C"

import (
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"time"

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
