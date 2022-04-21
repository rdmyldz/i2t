package leptonica

// #cgo LDFLAGS: -llept
// #include <leptonica/allheaders.h>
import "C"
import (
	"fmt"
	"log"
	"unsafe"
)

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
/*
IFF_DEFAULT added to the switch statement like it is a TIFF
because C.findFileFormatStream returns 17 for the test file 'multipage-sample.tif
but for the same test file, it just works as it should, in PixReadMem.
Because C.findFileFormatBuffer returns 4
*/
func PixRead(path string) ([]Pix, error) {
	var img Pix
	var pages []Pix
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	f := C.fopenReadStream(cPath)
	var cFormat C.int
	if ok := C.findFileFormatStream(f, &cFormat); ok != 0 {
		return nil, fmt.Errorf("could not find the file format")
	}
	log.Printf("format after findfileFormat %v\n", cFormat)
	switch cFormat {
	case IFF_PNG:
		log.Printf("the pic is PNG File: %v\n", cFormat)
		img = C.pixReadStreamPng(f)
		if img == nil {
			return nil, fmt.Errorf("error reading image: %s", path)
		}
		pages = append(pages, img)
	case IFF_TIFF, IFF_TIFF_PACKBITS, IFF_TIFF_RLE, IFF_TIFF_G3,
		IFF_TIFF_G4, IFF_TIFF_LZW, IFF_TIFF_ZIP, IFF_DEFAULT:
		log.Printf("the pic is TIFF File: %v\n", cFormat)

		for i := 0; ; i++ {
			img = C.pixReadStreamTiff(f, C.int(i))
			if img == nil {
				break
			}
			pages = append(pages, img)
		}
	case IFF_UNKNOWN:
		return nil, fmt.Errorf("unknown format: no pix returned")
	}

	if len(pages) < 1 {
		return nil, fmt.Errorf("no pages returned")
	}
	return pages, nil
}

func PixDestroy(p Pix) {
	C.pixDestroy((**C.PIX)(&p))
}

func DestroyPixes(pixes []Pix) {
	for _, p := range pixes {
		PixDestroy(p)
	}
}

// LEPT_DLL extern PIX * pixReadMem ( const l_uint8 *data, size_t size );
// LEPT_DLL extern l_ok findFileFormatBuffer ( const l_uint8 *buf, l_int32 *pformat );
// LEPT_DLL extern PIX * pixReadMemTiff ( const l_uint8 *cdata, size_t size, l_int32 n );
/*
	for JPEG
	https://tpgit.github.io/Leptonica/readfile_8c_source.html
	row:=00306  if ((pix = pixReadStreamJpeg(fp, READ_24_BIT_COLOR, 1, NULL, hint))
	pixReadMemJpeg( const l_uint8  *cdata, size_t size, l_int32 cmflag, l_int32 reduction, l_int32 *pnwarn, l_int32 hint)
	pixReadStreamJpeg( FILE *fp, l_int32 cmflag, l_int32 reduction, l_int32 *pnwarn, l_int32 hint)
	LEPT_DLL extern PIX * pixReadMemJpeg (const l_uint8 *data, size_t size, l_int32 cmflag, l_int32 reduction, l_int32 *pnwarn, l_int32 hint );
*/
func PixReadMem(data []byte) ([]Pix, error) {
	var img Pix
	var pages []Pix
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
		img = C.pixReadMemPng(f, c_ulong(size))
		if img == nil {
			return nil, fmt.Errorf("error reading png")
		}
		pages = append(pages, img)

	case IFF_TIFF, IFF_TIFF_PACKBITS, IFF_TIFF_RLE, IFF_TIFF_G3,
		IFF_TIFF_G4, IFF_TIFF_LZW, IFF_TIFF_ZIP:
		log.Printf("the pic is TIFF File: %v\n", cFormat)

		for pageNum := 0; ; pageNum++ {
			img = C.pixReadMemTiff(f, c_ulong(size), C.int(pageNum))
			if img == nil {
				break
			}
			pages = append(pages, img)
		}

	case IFF_BMP:
		log.Printf("the pic is BMP File: %v\n", cFormat)
		img = C.pixReadMemBmp(f, c_ulong(size))
		if img == nil {
			return nil, fmt.Errorf("error reading BMP")
		}
		pages = append(pages, img)

	case IFF_JFIF_JPEG:
		log.Printf("the pic is JPEG File: %v\n", cFormat)
		img = C.pixReadMemJpeg(f, c_ulong(size), 0, 1, nil, 0)
		if img == nil {
			return nil, fmt.Errorf("error reading JPEG")
		}
		pages = append(pages, img)

	case IFF_PNM:
		log.Printf("the pic is PNM File: %v\n", cFormat)
		img = C.pixReadMemPnm(f, c_ulong(size))
		if img == nil {
			return nil, fmt.Errorf("error reading PNM")
		}
		pages = append(pages, img)

	case IFF_GIF:
		log.Printf("the pic is GIF File: %v\n", cFormat)
		img = C.pixReadMemGif(f, c_ulong(size))
		if img == nil {
			return nil, fmt.Errorf("error reading GIF")
		}
		pages = append(pages, img)

	case IFF_WEBP:
		log.Printf("the pic is WEBP File: %v\n", cFormat)
		img = C.pixReadMemWebP(f, c_ulong(size))
		if img == nil {
			return nil, fmt.Errorf("error reading WEBP")
		}
		pages = append(pages, img)

	case IFF_SPIX:
		log.Printf("the pic is SPIX File: %v\n", cFormat)
		img = C.pixReadMemSpix(f, c_ulong(size))
		if img == nil {
			return nil, fmt.Errorf("error reading SPIX")
		}
		pages = append(pages, img)

	case IFF_JP2:
		return nil, fmt.Errorf("JP2 not supported")

	case IFF_UNKNOWN:
		return nil, fmt.Errorf("unknown format: no pix returned")
	}

	if len(pages) < 1 {
		return nil, fmt.Errorf("no pages returned")
	}
	return pages, nil
}
