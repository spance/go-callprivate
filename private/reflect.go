package private

import (
	"reflect"
	"unsafe"
)

const (
	kindMask        = (1 << 5) - 1
	flagMethodShift = 9
)

type flag_t uintptr

type value struct {
	typ *rtype
	ptr unsafe.Pointer
	flag_t
}

type rtype struct {
	size          uintptr
	ptrdata       uintptr
	hash          uint32         // hash of type; avoids computation in hash tables
	_             uint8          // unused/padding
	align         uint8          // alignment of variable with this type
	fieldAlign    uint8          // alignment of struct field with this type
	kind          uint8          // enumeration for C
	alg           uintptr        // algorithm table
	gcdata        uintptr        // garbage collection data
	_string       uintptr        // string form; unnecessary but undeniably useful
	*uncommonType                // (relatively) uncommon fields
	ptrToThis     *rtype         // type for pointer to this type, if used in binary or has methods
	zero          unsafe.Pointer // pointer to zero value
}

type uncommonType struct {
	name    *string  // name of type
	pkgPath *string  // import path; nil for built-in types like int, string
	methods []method // methods associated with type
}

// Method on non-interface type
type method struct {
	name    *string        // name of method
	pkgPath *string        // nil for exported Names; otherwise import path
	mtyp    *rtype         // method type (without receiver)
	typ     *rtype         // .(*FuncType) underneath (with receiver)
	ifn     unsafe.Pointer // fn used in interface call (one-word receiver)
	tfn     unsafe.Pointer // fn used for normal method call
}

// imethod represents a method on an interface type
type imethod struct {
	name    *string // name of method
	pkgPath *string // nil for exported Names; otherwise import path
	typ     *rtype  // .(*FuncType) underneath
}

type interfaceType struct {
	rtype   `reflect:"interface"`
	methods []imethod // sorted by hash
}

func (t *rtype) Kind() reflect.Kind { return reflect.Kind(t.kind & kindMask) }

func SetAccessible(val reflect.Value) (err error) {
	v := (*value)(unsafe.Pointer(&val))
	i := int(v.flag_t) >> flagMethodShift
	if v.flag_t&kindMask == 0 { // invalid reflect
		return
	}
	var pkgPath **string
	if v.typ.Kind() == reflect.Interface {
		tt := (*interfaceType)(unsafe.Pointer(v.typ))
		m := &tt.methods[i]
		if m.pkgPath != nil {
			pkgPath = &m.pkgPath
		}
	} else {
		ut := v.typ.uncommonType
		m := &ut.methods[i]
		if m.pkgPath != nil {
			pkgPath = &m.pkgPath
		}
	}
	if pkgPath != nil && *pkgPath != nil {
		err = memprotect(unsafe.Pointer(v.typ), unsafe.Pointer(*pkgPath), func() {
			*pkgPath = nil
		})
	}
	return
}
