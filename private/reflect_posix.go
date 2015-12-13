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

func memprotect(_p1, _p2 unsafe.Pointer, proc func()) error {
	var pageStart, plen uintptr
	var p1, p2 = uintptr(_p1) & pageMask, uintptr(_p2) & pageMask
	var no int
	if p1 <= p2 {
		pageStart, plen = p1, p2-p1+pageSize
	} else {
		pageStart, plen = p2, p1-p2+pageSize
	}
	no = sys_mprotect(pageStart, plen, _PROT_READ_WRITE)
	if no != 0 {
		return fmt.Errorf("mprotect errno=%d", no)
	}
	proc()
	no = sys_mprotect(pageStart, plen, _PROT_READ)
	if no != 0 {
		return fmt.Errorf("mprotect errno=%d", no)
	}
	return nil
}
