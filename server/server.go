package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"rrpc/message"
	"sync"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error: ", err)
			return
		}
		go server.ServeConn(conn)
	}
}

type Option struct {
	AuthType    message.AuthType
	MessageType message.Type
}

var DefaultOption = &Option{
	AuthType:    message.NullAuth,
	MessageType: message.GobType,
}

func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	defer func() { _ = conn.Close() }()
	var option *Option
	if err := json.NewDecoder(conn).Decode(&option); err != nil {
		log.Println("rpc server: options error: ", err)
		return
	}
	// TODO add auth
	f := message.MessageFuncMap[option.MessageType]
	if f == nil {
		log.Printf("rpc server: invalid code type: %s", option.MessageType)
		return
	}
	server.ServerMessage(f(conn))
}

var invalidRequest = struct{}{}

func (server *Server) ServerMessage(msg message.Message) {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		req, err := server.readRequest(msg)
		if err != nil {
			if req == nil {
				break
			}
			req.h.Error = err.Error()
			server.sendResponse(msg, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		go server.handleRequest(msg, req, sending, wg)
	}
	wg.Wait()
	_ = msg.Close()
}
func Accept(lis net.Listener) { DefaultServer.Accept(lis) }

type request struct {
	h            *message.Header
	argv, replyv reflect.Value
}

func (server *Server) readRequestHeader(msg message.Message) (*message.Header, error) {
	var h message.Header
	if err := msg.ReadHeader(&h); err != nil {
		if err != io.EOF && !errors.Is(err, io.ErrUnexpectedEOF) {
			log.Println("rpc server err in readRequestHeader: ", err)
		}
	}
	return &h, nil
}
func (server *Server) readRequest(msg message.Message) (*request, error) {
	h, err := server.readRequestHeader(msg)
	if err != nil {
		return nil, err
	}
	req := &request{h: h}

	req.argv = reflect.New(reflect.TypeOf(""))
	if err = msg.ReadBody(req.argv.Interface()); err != nil {
		log.Println("rpc server read argv err:", err)
	}
	return req, nil
}
func (server *Server) sendResponse(msg message.Message, h *message.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := msg.WriteMessage(h, body); err != nil {
		log.Println("rpc server err in sendResponse: ", err)
	}
}

func (server *Server) handleRequest(msg message.Message, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println(req.h, req.argv.Elem())
	req.replyv = reflect.ValueOf(fmt.Sprintf("rrpc resp %d", req.h.SeqNum))
	server.sendResponse(msg, req.h, req.replyv.Interface(), sending)
}
