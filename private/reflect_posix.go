// +build !windows

package private

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	_PROT_NONE  = 0x0
	_PROT_READ  = 0x1
	_PROT_WRITE = 0x2
	_PROT_EXEC  = 0x4
)

var (
	pageSize uintptr
	pageMask uintptr
)

func init() {
	pageSize = uintptr(syscall.Getpagesize())
	pageMask = ^(pageSize - 1)
}

func sys_mprotect(addr uintptr, n uintptr, prot uintptr) int

func memprotect(v unsafe.Pointer, writeable bool) error {
	var pageStart = uintptr(v) & pageMask
	var newProtect uintptr
	if writeable {
		newProtect = _PROT_READ | _PROT_WRITE
	} else {
		newProtect = _PROT_READ
	}
	no := sys_mprotect(pageStart, pageSize, newProtect)
	if no != 0 {
		return fmt.Errorf("mprotect errno=%d", no)
	}
	return nil
}
