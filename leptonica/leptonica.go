package leptonica

// #cgo LDFLAGS: -llept
// #include "leptonica/allheaders.h"
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"log"
	"unsafe"
)

type PIX struct {
	Pix *C.PIX
}

type Pix *C.PIX

const (
	IFF_UNKNOWN = iota
	IFF_BMP
	IFF_JFIF_JPEG
	IFF_PNG
	IFF_TIFF
	IFF_TIFF_PACKBITS
	IFF_TIFF_RLE
	IFF_TIFF_G3
	IFF_TIFF_G4
	IFF_TIFF_LZW
	IFF_TIFF_ZIP
	IFF_PNM
	IFF_PS
	IFF_GIF
	IFF_JP2
	IFF_WEBP
	IFF_LPDF
	IFF_DEFAULT
	IFF_SPIX
)

/*
https://github.com/golang/go/issues/13467
I am slightly more sympathetic to *C.char,
but even there I don't understand why the package API doesn't just use
appropriate Go types instead (like []byte).
*/
// PixRead returns unsafe.Pointer otherwise we can't use it in another package easily
func PixRead(path string) ([]*C.PIX, error) {
	var img *C.PIX
	var pages []*C.PIX
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	f := C.fopenReadStream(cPath)
	var cFormat C.int
	if ok := C.findFileFormatStream(f, &cFormat); ok != 0 {
		return nil, fmt.Errorf("could not find the file format")
	}
	switch cFormat {
	case IFF_PNG:
		log.Printf("the pic is PNG File: %v\n", cFormat)
		img = C.pixReadStreamPng(f)
		if img == nil {
			return nil, fmt.Errorf("error reading image: %s", path)
		}
		pages = append(pages, img)
	case IFF_TIFF, IFF_TIFF_PACKBITS, IFF_TIFF_RLE, IFF_TIFF_G3, IFF_TIFF_G4, IFF_TIFF_LZW, IFF_TIFF_ZIP:
		log.Printf("the pic is TIFF File: %v\n", cFormat)

		var cPageNum C.int
		// defer C.free(unsafe.Pointer(&cPageNum))
		if ok := C.tiffGetCount(f, &cPageNum); ok != 0 {
			return nil, fmt.Errorf("in PixSlice, error getting the number of pages")
		}
		log.Printf("page num: %v \n", cPageNum)
		for i := 0; i < int(cPageNum); i++ {
			img = C.pixReadStreamTiff(f, C.int(i))
			if img == nil {
				return nil, fmt.Errorf("in PixSlice, error reading")
			}
			pages = append(pages, img)
		}
	case IFF_UNKNOWN:
		return nil, fmt.Errorf("unknown format: no pix returned")
	}
	return pages, nil
}

func PixaReadTiff(path string) (*C.PIXA, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	f := C.fopenReadStream(cPath)
	var cFormat C.int
	if ok := C.findFileFormatStream(f, &cFormat); ok != 0 {
		return nil, fmt.Errorf("could not find the file format")
	}
	pixa := C.pixaReadStream(f)
	if pixa == nil {
		return nil, fmt.Errorf("error reading pixa")
	}

	return pixa, nil
}

func PixDestroy(p *C.PIX) {
	C.pixDestroy(&p)
}

// LEPT_DLL extern PIX * pixReadMem ( const l_uint8 *data, size_t size );
// LEPT_DLL extern l_ok findFileFormatBuffer ( const l_uint8 *buf, l_int32 *pformat );
// LEPT_DLL extern PIX * pixReadMemTiff ( const l_uint8 *cdata, size_t size, l_int32 n );
func PixReadMem(data []byte) ([]*C.PIX, error) {
	var img *C.PIX
	var pages []*C.PIX
	size := len(data)

	var cFormat C.int
	f := (*C.uchar)(unsafe.Pointer(&data[0]))
	if ok := C.findFileFormatBuffer(f, &cFormat); ok != 0 {
		return nil, fmt.Errorf("could not find the file format")
	}
	log.Printf("file format is: %v\n", cFormat)
	switch cFormat {
	case IFF_PNG:
		log.Printf("the pic is PNG File: %v\n", cFormat)
		img = C.pixReadMemPng(f, C.ulong(size))
		if img == nil {
			return nil, fmt.Errorf("error reading png")
		}
		pages = append(pages, img)
	case IFF_TIFF, IFF_TIFF_PACKBITS, IFF_TIFF_RLE, IFF_TIFF_G3, IFF_TIFF_G4, IFF_TIFF_LZW, IFF_TIFF_ZIP:
		log.Printf("the pic is TIFF File: %v\n", cFormat)

		for pageNum := 0; ; pageNum++ {
			img = C.pixReadMemTiff(f, C.ulong(size), C.int(pageNum))
			if img == nil {
				// return nil, fmt.Errorf("in PixForNetwork, error reading")
				break
			}
			pages = append(pages, img)
		}
	case IFF_UNKNOWN:
		return nil, fmt.Errorf("unknown format: no pix returned")
	}
	return pages, nil
}
