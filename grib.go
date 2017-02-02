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

type GRIBHandle struct {
	file     *GRIBFile
	c_handle *C.struct_grib_handle
	c_err    C.int
}

type GRIBKeysIterator struct {
	handle                       *GRIBHandle
	namespace                    string
	c_keys_iterator              *C.struct_grib_keys_iterator
	c_keys_iterator_filter_flags C.ulong
}

type GRIBFile struct {
	c_file *C.struct__IO_FILE
}

type Value struct {
	Key   string
	Value string
}

func (grib_file *GRIBFile) Open(path string) {
	var err error
	c_path := C.CString(path)
	c_fopen_flags := C.CString("r")
	defer C.free(unsafe.Pointer(c_path))
	defer C.free(unsafe.Pointer(c_fopen_flags))
	grib_file.c_file, err = C.fopen(c_path, c_fopen_flags)
	if err != nil {
		panic(fmt.Sprintf("Cannot open GRIB file %s: %s", path, err))
	}
}

func (grib_file *GRIBFile) Close() {
	C.fclose(grib_file.c_file)
}

func (grib_file *GRIBFile) Next() (grib_handle GRIBHandle, err error) {
	grib_handle.file = grib_file
	grib_handle.c_handle = C.codes_handle_new_from_file(nil, grib_file.c_file, C.PRODUCT_GRIB, &grib_handle.c_err)
	if grib_handle.c_handle == nil {
		err = fmt.Errorf("Cannot create GRIB handle from file, got NULL pointer from eccodes library with error code %d", grib_handle.c_err)
	}
	return
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

func (grib_keys_iterator *GRIBKeysIterator) Next() (value Value, err error) {
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
