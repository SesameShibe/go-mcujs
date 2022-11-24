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
	"net"
	"runtime"
	"time"

	"github.com/dop251/goja"
	"github.com/google/uuid"
)

const (
	OPCODE_REPORT_ID          = 0x01
	OPCODE_CHALLENGE_RESPONSE = 0x02
	OPCODE_STATE_CHANGE       = 0x03
	OPCODE_RPC_CALL           = 0x04
	OPCODE_RPC_RESPONSE       = 0x05
)

type JSCallback func(*goja.Runtime)

var MainLoopChan chan JSCallback = make(chan JSCallback, 1024)
var ticker *time.Ticker = time.NewTicker(100 * time.Millisecond)
var vm *goja.Runtime
var msgBuilder = MsgBuilder{}
var deviceId = uuid.New()


// The steps of attaching the device to the server are:
// 1. The device sends a OPCODE_REPORT_ID packet to the server.
// 2. The server sends a OPCODE_CHALLENGE_RESPONSE packet to the device.
// 3. The device sends a OPCODE_CHALLENGE_RESPONSE packet to the server.
// 4. If the server accepts the device, it sends a OPCODE_STATE_CHANGE packet to the device.
// 5. The device sends a OPCODE_STATE_CHANGE packet to the server.
// 6. The device enters the main loop to receive and handle OPCODE_RPC_CALL packets.
func attachDevice() {
	// Create a UDP connection to the server.
	conn, err := net.Dial("udp", "127.0.0.1:8801")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer conn.Close()

	// Send the message to the server.
	reportDeviceId(conn)

	// Handle the response.
	handleChallengeResponse(conn)

	// Handle the state change.
	handleStateChange(conn)

	// Enter the main loop.
	for {
		// Handle the RPC call.
		payload, err := udpGetPayload(conn)
		if err != nil {
			log.Fatal("udpGetPayload:", err)
		}

		opcode := payload[0]
		if opcode != OPCODE_RPC_CALL {
			log.Fatal("unexpected opcode:", opcode)
		}

		// Currently not implemented.
		log.Println("RPC call:", payload[1:])

		// Send the response.
		response := make([]byte, 1+32)
		response[0] = OPCODE_RPC_RESPONSE
		copy(response[1:], payload[1:])
		pkt, err := msgBuilder.BuildMsg(response, deviceId)
		if err != nil {
			log.Fatal("BuildMsg:", err)
		}

		_, err = conn.Write(pkt)
		if err != nil {
			log.Fatal("conn.Write:", err)
		}
	}
}

func udpGetPayload(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	
	pkt := buf[:n]
	payload, deviceId, err := msgBuilder.ParseMsg(pkt)
	if err != nil {
		return nil, err
	}
	log.Println("deviceId:", deviceId)

	return payload, nil
}

func reportDeviceId(conn net.Conn) {
	// Payload format:
	// 1 byte: OPCODE_REPORT_ID
	// 32 bytes: Device ID
	payload := make([]byte, 1+32)
	payload[0] = OPCODE_REPORT_ID
	copy(payload[1:], deviceId[:])
	pkt, err := msgBuilder.BuildMsg(payload, deviceId)
	if err != nil {
		log.Fatal("build msg:", err)
	}

	_, err = conn.Write(pkt)
	if err != nil {
		log.Fatal("send msg:", err)
	}
}

func handleChallengeResponse(conn net.Conn) {
	// Payload format:
	// 1 byte: OPCODE_CHALLENGE_RESPONSE
	// 32 bytes: challenge content
	payload, err := udpGetPayload(conn)
	if err != nil {
		log.Fatal("get payload:", err)
	}

	opcode := payload[0]
	if opcode != OPCODE_CHALLENGE_RESPONSE {
		log.Fatal("unexpected opcode:", opcode)
	}

	challenge := payload[1:]
	// The challenge content should starts with "What's your name" and padded with 0.
	if string(challenge[:15]) != "What's your name" {
		log.Fatal("unexpected challenge content:", challenge)
	}

	// Send the response to the server.
	response := make([]byte, 1+32)
	response[0] = OPCODE_CHALLENGE_RESPONSE
	copy(response[1:], challenge)
	pkt, err := msgBuilder.BuildMsg(response, deviceId)
	if err != nil {
		log.Fatal("build msg:", err)
	}

	_, err = conn.Write(pkt)
}

func handleStateChange(conn net.Conn) {
	// Payload format:
	// 1 byte: OPCODE_STATE_CHANGE
	// 1 byte: state
	payload, err := udpGetPayload(conn)
	if err != nil {
		log.Fatal("get payload:", err)
	}

	opcode := payload[0]
	if opcode != OPCODE_STATE_CHANGE {
		log.Fatal("unexpected opcode:", opcode)
	}

	state := payload[1]
	if state != 1 {
		log.Fatal("unexpected state:", state)
	}

	// Respond to the server.
	response := make([]byte, 1+1)
	response[0] = OPCODE_STATE_CHANGE
	response[1] = 1
	pkt, err := msgBuilder.BuildMsg(response, deviceId)
	if err != nil {
		log.Fatal("build msg:", err)
	}

	_, err = conn.Write(pkt)
}


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

	msgBuilder.Init()
	go attachDevice()

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
