package eccodes

// #cgo pkg-config: eccodes
// #include <eccodes.h>
// #include <stdlib.h>
import "C"

import (
	"fmt"
	"unsafe"
)

type GRIBFile struct {
	c_file *C.struct__IO_FILE
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
