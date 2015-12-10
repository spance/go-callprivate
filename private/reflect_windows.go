package private

import (
	"syscall"
	"unsafe"
)

var (
	virtualProtect *syscall.Proc
	virtualQuery   *syscall.Proc
)

const (
	_PAGE_NOACCESS          = 0x01
	_PAGE_READONLY          = 0x02
	_PAGE_READWRITE         = 0x04
	_PAGE_EXECUTE_READ      = 0x20
	_PAGE_EXECUTE_READWRITE = 0x40
)

func init() {
	k := syscall.MustLoadDLL("kernel32.dll")
	virtualProtect = k.MustFindProc("VirtualProtect")
	virtualQuery = k.MustFindProc("VirtualQuery")
}

type MEMORY_BASIC_INFORMATION struct {
	BaseAddress       uintptr
	AllocationBase    uintptr
	AllocationProtect uint32
	RegionSize        uintptr
	State             uint32
	Protect           uint32
	Type              uint32
}

func memprotect(v unsafe.Pointer, writeable bool) error {
	var mbi MEMORY_BASIC_INFORMATION
	virtualQuery.Call(uintptr(v), uintptr(unsafe.Pointer(&mbi)), unsafe.Sizeof(mbi))
	var newProtect uintptr
	var oldProtect uint32
	var err error
	if writeable {
		newProtect = _PAGE_EXECUTE_READWRITE
	} else {
		newProtect = _PAGE_EXECUTE_READ
	}
	_, _, err = virtualProtect.Call(mbi.BaseAddress, mbi.RegionSize, newProtect, uintptr(unsafe.Pointer(&oldProtect)))
	if en, y := err.(syscall.Errno); y && en == 0 {
		return nil
	}
	return err
}
