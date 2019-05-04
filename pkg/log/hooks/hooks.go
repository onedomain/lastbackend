//
// KULADO INC. CONFIDENTIAL
// __________________
//
// [2014] - [2019] KULADO INC.
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

package hooks

import (
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
	"strings"
)

type ContextHook struct {
	Skip int
}

func (ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h ContextHook) Fire(entry *logrus.Entry) error {
	if h.Skip == 0 {
		h.Skip = 8
	}

	pc := make([]uintptr, 3, 3)
	cnt := runtime.Callers(h.Skip, pc)

	for i := 0; i < cnt; i++ {
		fu := runtime.FuncForPC(pc[i] - 1)
		name := fu.Name()
		if !strings.Contains(name, "github.com/sirupsen/logrus") {
			file, line := fu.FileLine(pc[i] - 1)
			entry.Data["func"] = path.Base(name)
			entry.Data["file"] = file
			entry.Data["line"] = line
			break
		}
	}
	return nil
}
