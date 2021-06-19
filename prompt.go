package main

import (
	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/dop251/goja"
)

func completer(d prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{}
}

func executor(in string) {
	CallInJSLoop(func(r *goja.Runtime) {
		ret, err := vm.RunString(in)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(ret)
		}
	})
}

func PromptLoop() {
	p := prompt.New(executor, completer)
	p.Run()
}
