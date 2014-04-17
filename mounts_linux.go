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

// +build linux

// Package mount provides file system utiilty functions
package mount

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	"syscall"
)

type Mount struct {
	Device     string
	Path       string
	Filesystem string
	Flags      string
}

func Mounts() ([]Mount, error) {
	file, err := os.Open("/proc/self/mounts")
	if err != nil {
		return nil, err
	}
	defer checkClose(file)
	mounts := []Mount(nil)
	reader := bufio.NewReaderSize(file, 64*1024)
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return mounts, nil
			}
			return nil, err
		}
		if isPrefix {
			return nil, syscall.EIO
		}
		parts := strings.SplitN(string(line), " ", 5)
		if len(parts) != 5 {
			return nil, syscall.EIO
		}
		mounts = append(mounts, Mount{parts[0], parts[1], parts[2], parts[3]})
	}
	panic("unreachable")
}

func isFs(path string, magic int64) bool {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		panic(os.NewSyscallError("statfs", err))
	}
	return stat.Type == magic
}

// Returns whether the path is a tmpfs file system
func IsTmpfs(path string) bool {
	return isFs(path, TMPFS_MAGIC)
}

// Returns whether the path is a hugetlbfs file system
func IsHugetlbfs(path string) bool {
	return isFs(path, HUGETLBFS_MAGIC)
}

// Returns whether the path is a btrfs file system
func IsBtrfs(path string) bool {
	return isFs(path, BTRFS_SUPER_MAGIC)
}

// Returns whether the path is an ext4 file system
func IsExt4(path string) bool {
	return isFs(path, EXT4_SUPER_MAGIC)
}

// Returns whether the path is an in memory file system file system (either hugetlbfs or tmpfs)
func IsMemoryFs(path string) bool {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		panic(os.NewSyscallError("statfs", err))
	}
	t := stat.Type
	return t == HUGETLBFS_MAGIC || t == TMPFS_MAGIC
}

// MountTmpfs mounts a tmpfs file system of the specified size at the provided path
func MountTmpfs(path string, size int64) error {
	if size < 0 {
		panic("MountTmpfs: size < 0")
	}
	var flags uintptr
	flags = syscall.MS_NOATIME | syscall.MS_SILENT
	flags |= syscall.MS_NODEV | syscall.MS_NOEXEC | syscall.MS_NOSUID
	options := ""
	if size >= 0 {
		options = "size=" + strconv.FormatInt(size, 10)
	}
	err := syscall.Mount("tmpfs", path, "tmpfs", flags, options)
	return os.NewSyscallError("mount", err)
}

// MountTmpfs mounts a hugelbfs file system of the specified size with
// a specified page size at the provided path
func MountHugetlbfs(path string, pagesize int, size int64) error {
	if pagesize != 2<<20 && pagesize != 1<<30 {
		panic("MountHugetlbfs: invalid page size")
	}
	if size < 0 {
		panic("MountHugetlbfs: size < 0")
	}
	var flags uintptr
	flags = syscall.MS_NOATIME
	flags |= syscall.MS_NODEV | syscall.MS_NOEXEC | syscall.MS_NOSUID
	options := "pagesize=" + strconv.Itoa(pagesize)
	if size >= 0 {
		options += ",size=" + strconv.FormatInt(size, 10)
	}
	err := syscall.Mount("hugetlbfs", path, "hugetlbfs", flags, options)
	return os.NewSyscallError("mount", err)
}

// RemountReadOnly remounts the specified source device at the target path as read only
func RemountReadOnly(source, target string) error {
	var flags uintptr
	flags = syscall.MS_RDONLY | syscall.MS_REMOUNT
	err := syscall.Mount(source, target, "", flags, "")
	return os.NewSyscallError("mount", err)
}

func readLine(name string) string {
	f, err := os.Open(name)
	if err != nil {
		return ""
	}
	defer checkClose(f)
	b := bufio.NewReaderSize(f, 4096)
	line, isPrefix, err := b.ReadLine()
	if err != nil || isPrefix {
		return ""
	}
	return string(line)
}

func readLines(name string) (lines []string) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer checkClose(f)
	b := bufio.NewReaderSize(f, 4096)
	for {
		line, isPrefix, err := b.ReadLine()
		if err != nil || isPrefix {
			return
		}
		lines = append(lines, string(line))
	}
	return
}

func InMountNamespace() bool {
	initMount := readLine("/proc/1/mountinfo")
	myMount := readLine("/proc/self/mountinfo")
	if initMount == "" || myMount == "" {
		panic("MountNamespace: don't know")
	}
	return initMount != myMount
}

func MountNamespace() {
	if !InMountNamespace() {
		panic("not in a separate mount namespace")
	}
	// disable shared mount propagation
	const flags = syscall.MS_REC | syscall.MS_PRIVATE
	err := syscall.Mount("none,", "/", "none", flags, "")
	if err != nil {
		panic("MountNamespace: mount: " + err.Error())
	}
}

var ext4MountOptions = strings.Join([]string{
	"journal_checksum",
	"journal_ioprio=0",
	"barrier=1",
	"data=ordered",
	"errors=remount-ro",
}, ",")

// MountExt4 mounts the specifed ext4 volume device at the specified path
// The readonly flag allows the device to be mounted as read only
// The sync flag can be used to enforce synchronous device access
func MountExt4(device, path string, readonly, sync bool) error {
	var flags uintptr
	flags = syscall.MS_NOATIME | syscall.MS_SILENT
	flags |= syscall.MS_NODEV | syscall.MS_NOEXEC | syscall.MS_NOSUID
	if readonly {
		flags |= syscall.MS_RDONLY
	}
	if sync {
		flags |= syscall.MS_SYNCHRONOUS | syscall.MS_DIRSYNC
	}
	err := syscall.Mount(device, path, "ext4", flags, ext4MountOptions)
	return os.NewSyscallError("mount", err)
}
