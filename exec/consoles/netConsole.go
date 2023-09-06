// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package consoles

import (
	"crypto/tls"
	"fmt"
	"khalehla/exec/messages"
	"log"
	"net"
	"sync"
)

// NetConsole implements a net server which provides console functionality to 0 or more
// connected Console handlers (such as is found in kdte).
type NetConsole struct {
	pendingReadReplyMessages map[int]*messages.ReadReplyMessage
	connections              []net.Conn
	mutex                    sync.Mutex
	terminate                bool
}

func NewNetConsole(hostAddress string, port int) Console {
	c := &NetConsole{
		pendingReadReplyMessages: make(map[int]*messages.ReadReplyMessage),
		connections:              make([]net.Conn, 0),
	}

	go c.listener()
	return c
}

func (c *NetConsole) listener() {
	port := 2200

	log.Printf("listening on port %d\n", port)
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("ERROR:%v\n", err.Error())
		return
	}

	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			fmt.Printf("ERROR:%v\n", err.Error())
		}
	}(l)

	for !c.terminate {
		connection, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("accepted connection from %s\n", connection.RemoteAddr())

		c.mutex.Lock()
		c.connections = append(c.connections, connection)
		c.mutex.Unlock()
	}

	c.closeConnections()
}

func (c *NetConsole) secureListener() {
	port := 2200
	certFile := "resources/certificate.pem"
	keyFile := "resources/key.pem"

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	log.Printf("listening on port %d\n", port)
	l, err := tls.Listen("tcp", fmt.Sprintf(":%d", port), config)
	if err != nil {
		fmt.Printf("ERROR:%v\n", err.Error())
		return
	}

	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			fmt.Printf("ERROR:%v\n", err.Error())
		}
	}(l)

	for !c.terminate {
		connection, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("accepted connection from %s\n", connection.RemoteAddr())

		c.mutex.Lock()
		c.connections = append(c.connections, connection)
		c.mutex.Unlock()
	}

	c.closeConnections()
}

func (c *NetConsole) closeConnections() {
	for _, conn := range c.connections {
		log.Printf("shutting down connection from %s\n", conn.RemoteAddr())
		_ = conn.Close()
	}

	c.connections = make([]net.Conn, 0)
}

func (c *NetConsole) removeConnection(index int) {
	conn := c.connections[index]
	fmt.Printf("Removing connection to %s\n", conn.RemoteAddr())
	_ = conn.Close()
	c.connections = append(c.connections[:index], c.connections[index+1:]...)
}

func (c *NetConsole) sendMessage(message []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for cx := 0; cx < len(c.connections); {
		conn := c.connections[cx]
		_, err := conn.Write(message)
		if err != nil {
			fmt.Printf("%s\n", err)
			_ = conn.Close()
			c.removeConnection(cx)
		} else {
			cx++
		}
	}
}

// -----------------------------------------------------------------------------

// ClearReadReplyMessage is invoked by the exec to indicate that the particular solicited
// (i.e., read-reply) message is no longer outstanding. The console should take whatever
// steps are necessary (if any) to visually indicate such to the operator.
func (c *NetConsole) ClearReadReplyMessage(messageId int) (err error) {
	// TODO send net message to clear RR msg
	delete(c.pendingReadReplyMessages, messageId)
	return nil
}

func (c *NetConsole) Close() {
	c.closeConnections()
}

func (c *NetConsole) IsReady() bool {
	return true // TODO return true only if connected
}

// PollSolicitedInput is invoked by the exec to ask whether the console operator has responded
// to a particular solicited (i.e., read-reply) message.
func (c *NetConsole) PollSolicitedInput(messageId int) (response string, hasInput bool, err error) {
	// TODO are we connected? if so, pass this to the client
	return "", false, nil
}

// PollUnsolicitedInput is invoked by the exec to ask whether the console operator has provided
// unsolicited input (i.e., a console key-in)
func (c *NetConsole) PollUnsolicitedInput() (input string, hasInput bool) {
	return "", false // TODO
}

// Reset is invoked by the exec to cause the console to reset itself.
// This might simply result in clearing the console screen.
func (c *NetConsole) Reset() {
	c.pendingReadReplyMessages = make(map[int]*messages.ReadReplyMessage)
	//	TODO screen stuff
}

// SendReadOnlyMessage is invoked by the exec to send a read only message to the console
func (c *NetConsole) SendReadOnlyMessage(message *messages.ReadOnlyMessage) {
	c.sendMessage(message.Serialize())
}

// SendReadReplyMessage is invoked by the exec to send a read-reply message to the console
func (c *NetConsole) SendReadReplyMessage(message *messages.ReadReplyMessage) {
	// TODO
}

// SendStatusMessage is invoked by the exec to send a status message to the console
func (c *NetConsole) SendStatusMessage(message *messages.StatusMessage) {
	// TODO
}
