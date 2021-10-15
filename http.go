package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dop251/goja"
	"github.com/gorilla/websocket"
)

type WSConn struct {
	c        *websocket.Conn
	Onbinary func(*WSConn, []byte)
	Ontext   func(*WSConn, string)
}

type ServerRequest struct {
	r      *http.Request
	Path   string
	Method string
}

type ServerResponse struct {
	w         http.ResponseWriter
	completed chan int
}

type Server struct {
	port        int
	server      *http.Server
	mux         *http.ServeMux
	upgrader    websocket.Upgrader
	Onwebsocket func(*WSConn)
	Onrequest   func(*ServerRequest, *ServerResponse)
}

func (ws *WSConn) SendText(t string) bool {
	err := ws.c.WriteMessage(websocket.TextMessage, []byte(t))
	if err != nil {
		return false
	}
	return true
}

func (ws *WSConn) SendBinary(m []byte) bool {
	err := ws.c.WriteMessage(websocket.BinaryMessage, m)
	if err != nil {
		return false
	}
	return true
}

func (ws *WSConn) Close() {
	ws.c.Close()
}

func (s *Server) Close() {
	s.server.Close()
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	completed := make(chan int, 1)
	if r.Method == http.MethodGet {
		if r.Header.Get("upgrade") == "websocket" {
			c, err := s.upgrader.Upgrade(w, r, nil)
			if err != nil {
				return
			}
			wsConn := &WSConn{c: c}
			CallInJSLoop(func(r *goja.Runtime) {
				if s.Onwebsocket == nil {
					c.Close()
					return
				}
				s.Onwebsocket(wsConn)
				go func() {
					for {
						msgType, buf, err := c.ReadMessage()
						if err != nil {
							log.Println("WS read: ", err)
							c.Close()
							break
						}
						CallInJSLoop(func(r *goja.Runtime) {
							switch msgType {
							case websocket.TextMessage:
								if wsConn.Ontext != nil {
									wsConn.Ontext(wsConn, string(buf))
								}
							case websocket.BinaryMessage:
								if wsConn.Onbinary != nil {
									wsConn.Onbinary(wsConn, buf)
								}
							}
						})
					}
				}()
			})
			return
		}
	}
	CallInJSLoop(func(vm *goja.Runtime) {
		if s.Onrequest == nil {
			completed <- 1
			return
		}
		s.Onrequest(&ServerRequest{
			r:      r,
			Path:   r.URL.Path,
			Method: r.Method,
		}, &ServerResponse{
			w:         w,
			completed: completed,
		})
	})
	<-completed
}

func (resp *ServerResponse) Write(buf []byte) {
	resp.w.Write(buf)
}

func (resp *ServerResponse) Close() {
	resp.completed <- 1
}

func HttpInit() {
	obj := vm.NewObject()
	obj.Set("createServer", func(port int, staticPath string, apiPath string, flags int) *Server {
		mux := http.NewServeMux()
		server := &Server{
			port:   port,
			server: &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", port), Handler: mux},
			mux:    mux,
			upgrader: websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		}
		if staticPath != "" {
			mux.Handle("/", http.FileServer(http.Dir(staticPath)))
		}
		mux.HandleFunc(apiPath, func(w http.ResponseWriter, r *http.Request) {
			server.handleRequest(w, r)
		})
		go func() {
			err := server.server.ListenAndServe()
			if err != nil {
				log.Println("Httpd ListenAndServe: ", err)
			}
		}()
		return server
	})
	vm.Set("http", obj)
}
