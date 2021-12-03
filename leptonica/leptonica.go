package leptonica

// #cgo LDFLAGS: -llept
// #include "leptonica/allheaders.h"
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"unsafe"
)

type PIX struct {
	Pix *C.PIX
}

/*
https://github.com/golang/go/issues/13467
I am slightly more sympathetic to *C.char,
but even there I don't understand why the package API doesn't just use
appropriate Go types instead (like []byte).
*/
// PixRead returns unsafe.Pointer otherwise we can't use it in another package easily
func PixRead(path string) (*C.PIX, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	img := C.pixRead(cPath)
	if img == nil {
		return nil, fmt.Errorf("error reading image: %s", path)
	}
	// defer C.pixDestroy(&img)
	// p := PIX{Pix: img}
	return img, nil
}

func PixDestroy(p *C.PIX) {
	C.pixDestroy(&p)
}

func PixReadMem(data []byte) (*C.PIX, error) {
	cSize := len(data)
	fmt.Printf("lenSize: %v\n", cSize)
	img := C.pixReadMem((*C.uchar)(unsafe.Pointer(&data[0])), C.ulong(cSize))
	if img == nil {
		return nil, fmt.Errorf("error reading image: ")
	}
	// defer C.pixDestroy(&img)
	// p := PIX{Pix: img}
	return img, nil
}
