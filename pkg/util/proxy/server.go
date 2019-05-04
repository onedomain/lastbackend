//
// KULADO INC. CONFIDENTIAL
// __________________
//
// [2014] - [2018] KULADO INC.
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of KULADO INC. and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to KULADO INC.
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from KULADO INC..
//

package proxy

import (
	"encoding/binary"
	protoio "github.com/gogo/protobuf/io"
	"github.com/onedomain/lastbackend/pkg/log"
	"net"
	"time"
)

type Server struct {
	Addr        string
	IdleTimeout time.Duration
	inShutdown  bool
	conns       map[Conn]bool
	listener    *net.Listener
}

func (srv Server) Listen(handler Handler) error {
	addr := srv.Addr
	if addr == "" {
		addr = ":2963"
	}
	log.Debugf("starting server on %v", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	for {

		c, err := listener.Accept()
		if err != nil {
			log.Debugf("error accepting connection %v", err)
			continue
		}

		conn := Conn{
			Conn:        c,
			IdleTimeout: srv.IdleTimeout,
			done:        make(chan bool),
			error:       make(chan string),
			writer:      protoio.NewUint32DelimitedWriter(c, binary.BigEndian),
		}

		log.Debugf("accepted connection from %v", conn.RemoteAddr())
		go conn.Handle(handler)

		srv.conns[conn] = true
	}

}
