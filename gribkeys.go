package eccodes

// #cgo pkg-config: eccodes
// #include <eccodes.h>
// #include <stdlib.h>
import "C"

import (
	"fmt"
	"unsafe"
)

const MAX_VAL_LEN int = 1024

type GRIBKeysIterator struct {
	handle                       *GRIBHandle
	namespace                    string
	c_keys_iterator              *C.struct_grib_keys_iterator
	c_keys_iterator_filter_flags C.ulong
}

type GRIBKeyValue struct {
	Key   string
	Value string
}

func (grib_keys_iterator *GRIBKeysIterator) Next() (value GRIBKeyValue, err error) {
	if grib_keys_iterator.c_keys_iterator == nil {
		err = fmt.Errorf("GRIB key iterator handle is NULL, cannot get values")
		return
	}
	if C.codes_keys_iterator_next(grib_keys_iterator.c_keys_iterator) == 0 {
		err = fmt.Errorf("Reached end of keys iterator")
		return
	}
	var c_vlen C.size_t
	c_vlen = (C.size_t)(MAX_VAL_LEN)
	buf := make([]byte, MAX_VAL_LEN)
	value.Key = C.GoString(C.codes_keys_iterator_get_name(grib_keys_iterator.c_keys_iterator))
	C.codes_keys_iterator_get_string(grib_keys_iterator.c_keys_iterator, (*C.char)(unsafe.Pointer(&buf[0])), &c_vlen)
	value.Value = string(buf)
	return
}

func (grib_keys_iterator *GRIBKeysIterator) Free() {
	if grib_keys_iterator.c_keys_iterator != nil {
		C.codes_keys_iterator_delete(grib_keys_iterator.c_keys_iterator)
	}
}
