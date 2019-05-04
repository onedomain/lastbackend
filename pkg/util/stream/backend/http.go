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

package backend

import (
	"context"
	"github.com/onedomain/lastbackend/pkg/log"
	"net/http"
)

type HttpWriter struct {
	writer http.ResponseWriter
	done   chan bool
}

func (hw *HttpWriter) Disconnect() {
	<-hw.done
}

func (hw *HttpWriter) Write(data []byte) {

	var err error

	_, err = hw.writer.Write(data)
	if err == context.Canceled {
		log.V(logLevel).Debug("Stream is canceled")
		hw.done <- true
		return
	}

	if f, ok := hw.writer.(http.Flusher); ok {
		f.Flush()
	}
}

func (hw *HttpWriter) End() error {
	hw.done <- true
	return nil
}

func NewHttpWriterBackend(writer http.ResponseWriter) StreamBackend {
	hw := new(HttpWriter)
	hw.writer = writer
	hw.done = make(chan bool)
	return hw
}
