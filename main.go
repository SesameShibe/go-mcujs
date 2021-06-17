package main

/*
#cgo CFLAGS: -Ilib/SDL2/include -Ilib/lvgl
#cgo LDFLAGS: -L./lib/SDL2/lib/x64 -lSDL2 -L./lib/lvgl/build -l lvgl
#include "mcujs.h"
*/
import "C"

import (
	"log"
	"time"

	"github.com/dop251/goja"
)

var ticker *time.Ticker = time.NewTicker(100 * time.Millisecond)

var vm *goja.Runtime

func print(call goja.FunctionCall) goja.Value {
	s := make([]interface{}, len(call.Arguments))
	for i, v := range call.Arguments {
		s[i] = v.String()
	}
	log.Println(s...)
	return nil
}

func main() {
	C.mjsInit()

	vm = goja.New()
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())
	vm.Set("print", print)
	vm.RunString("print('hello','world')")

	for {
		select {
		case <-ticker.C:
			if C.mjsEventLoop() == 1 {
				log.Println("SDL_QUIT received, exiting...")
				return
			}
		}

	}
}
