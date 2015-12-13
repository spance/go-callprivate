package private

import (
	"syscall"
	"unsafe"
)

var (
	virtualProtect *syscall.Proc
	virtualQuery   *syscall.Proc
	pageMask       uintptr
	mbiSize        uintptr
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
	mbiSize = unsafe.Sizeof(_MEMORY_BASIC_INFORMATION{})
	pageMask = ^(uintptr(syscall.Getpagesize()) - 1)
}

type _MEMORY_BASIC_INFORMATION struct {
	BaseAddress       uintptr
	AllocationBase    uintptr
	AllocationProtect uint32
	RegionSize        uintptr
	State             uint32
	Protect           uint32
	Type              uint32
}

func (m *_MEMORY_BASIC_INFORMATION) equals(o *_MEMORY_BASIC_INFORMATION) bool {
	return m.BaseAddress == o.BaseAddress && m.RegionSize == o.RegionSize
}

func queryVirtual(p1, p2 uintptr) []_MEMORY_BASIC_INFORMATION {
	var mbi [2]_MEMORY_BASIC_INFORMATION
	var m1 = &mbi[0]
	virtualQuery.Call(p1, uintptr(unsafe.Pointer(m1)), mbiSize)
	if p1&pageMask == p2&pageMask {
		return mbi[:1]
	}
	var m2 = &mbi[1]
	virtualQuery.Call(p2, uintptr(unsafe.Pointer(m2)), mbiSize)
	if m1.equals(m2) {
		return mbi[:1]
	}
	return mbi[:]
}

func memprotect(p1, p2 unsafe.Pointer, proc func()) (err error) {
	var mbi = queryVirtual(uintptr(p1), uintptr(p2))
	var oldProtect [2]uintptr
	for i := 0; i < len(mbi); i++ {
		m, op := &mbi[i], uintptr(unsafe.Pointer(&oldProtect[i]))
		_, _, err = virtualProtect.Call(m.BaseAddress, m.RegionSize, _PAGE_EXECUTE_READWRITE, op)
		if en, y := err.(syscall.Errno); !y || en != 0 {
			return
		}
	}
	proc()
	for i := 0; i < len(mbi); i++ {
		m, op := &mbi[i], uintptr(unsafe.Pointer(&oldProtect[i]))
		_, _, err = virtualProtect.Call(m.BaseAddress, m.RegionSize, op, 0)
		if en, y := err.(syscall.Errno); !y || en != 0 {
			return
		}
	}
	return nil
}
