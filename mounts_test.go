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
	"fmt"
	"os"
	"syscall"
	"testing"
)

func TestMounts(t *testing.T) {
	mounts, err := Mounts()
	if err != nil {
		t.Fatalf("Mounts failed: %v", err)
	}
	for _, mount := range mounts {
		fmt.Printf("%+v\n", mount)
	}
}

func TestStatfs(t *testing.T) {
	var stat syscall.Statfs_t
	err := syscall.Statfs("/dev/shm", &stat)
	if err == nil {
		t.Logf("/dev/shm: %+v", stat)
	}
	err = syscall.Statfs("/dev/hugepages", &stat)
	if err == nil {
		t.Logf("/dev/hugepages: %+v", stat)
	}
	f, err := os.Open("/dev/hugepages")
	if err != nil {
		panic(err)
	}
	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", fi.Sys())

}

func TestIsFs(t *testing.T) {
	if !IsTmpfs("/dev/shm") {
		panic("IsTmpfs failed on /dev/shm")
	}
	if IsTmpfs("/dev/hugepages") {
		panic("IsTmpfs didn't fail on /dev/hugepages")
	}
	if !IsHugetlbfs("/dev/hugepages") {
		panic("IsHugetlbfs failed on /dev/hugepages")
	}
	if IsHugetlbfs("/dev/shm") {
		panic("IsHugetlbfs didn't fail on /dev/shm")
	}
}
