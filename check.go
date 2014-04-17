//   Copyright 2014 Vastech SA (PTY) LTD
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package mount

import (
	"io"
	"runtime"
)

type runtimeError struct {
	Err error
}

var _ runtime.Error = (*runtimeError)(nil)

func (this *runtimeError) Error() string {
	return this.Err.Error()
}

func (this *runtimeError) RuntimeError() {
}

// The checkClose function calls close on a Closer and panics with a
// runtime error if the Closer returns an error
func checkClose(c io.Closer) {
	if err := c.Close(); err != nil {
		panic(&runtimeError{err})
	}
}