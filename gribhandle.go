package eccodes

// #cgo pkg-config: eccodes
// #include <eccodes.h>
// #include <stdlib.h>
import "C"

import (
	"fmt"
	"unsafe"
)

type GRIBHandle struct {
	file     *GRIBFile
	c_handle *C.struct_grib_handle
	c_err    C.int
}

func (grib_handle *GRIBHandle) MakeKeyIterator(namespace string) (grib_keys_iterator GRIBKeysIterator, err error) {
	c_namespace := C.CString(namespace)
	defer C.free(unsafe.Pointer(c_namespace))
	grib_keys_iterator.c_keys_iterator_filter_flags = C.CODES_KEYS_ITERATOR_ALL_KEYS | C.CODES_KEYS_ITERATOR_SKIP_DUPLICATES
	grib_keys_iterator.c_keys_iterator = C.codes_keys_iterator_new(grib_handle.c_handle, grib_keys_iterator.c_keys_iterator_filter_flags, c_namespace)
	if grib_keys_iterator.c_keys_iterator == nil {
		err = fmt.Errorf("Cannot create GRIB key iterator on namespace '%s'", namespace)
	}
	return
}

func (grib_handle *GRIBHandle) Free() {
	if grib_handle.c_handle != nil {
		C.codes_handle_delete(grib_handle.c_handle)
	}
}
