// +build !windows

package private

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	_PROT_NONE       = 0x0
	_PROT_READ       = 0x1
	_PROT_WRITE      = 0x2
	_PROT_EXEC       = 0x4
	_PROT_READ_WRITE = _PROT_READ | _PROT_WRITE
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

func memprotect(v unsafe.Pointer, proc func()) error {
	var pageStart = uintptr(v) & pageMask
	var no int
	no = sys_mprotect(pageStart, pageSize, _PROT_READ_WRITE)
	if no != 0 {
		return fmt.Errorf("mprotect errno=%d", no)
	}
	proc()
	no = sys_mprotect(pageStart, pageSize, _PROT_READ)
	if no != 0 {
		return fmt.Errorf("mprotect errno=%d", no)
	}
	return nil
}
