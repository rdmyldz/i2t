package tesseract

// #cgo CFLAGS: -g -Wall
// #cgo pkg-config: tesseract
// #cgo LDFLAGS: -llept
// #include "/usr/include/tesseract/capi.h"
// #include "/usr/include/leptonica/allheaders.h"
// #include <stdlib.h>
// #include <stdio.h>
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/rdmyldz/i2t/leptonica"
)

type TessBaseAPI struct {
	TBA *C.TessBaseAPI
}

// type TessBaseAPI interface {
// }

func TessVersion() string {
	vers := C.TessVersion()
	return C.GoString(vers)
}

func TessBaseAPICreate(language string) (*TessBaseAPI, error) {
	handle := C.TessBaseAPICreate()
	// handle := C.TessBaseAPICreate()

	// cDatapath := C.CString(datapath)
	// defer C.free(unsafe.Pointer(cDatapath))

	cLanguage := C.CString(language)
	defer C.free(unsafe.Pointer(cLanguage))

	res := C.TessBaseAPIInit3(handle, nil, cLanguage)
	if res != 0 {
		return nil, fmt.Errorf("error initializing tesseract")
	}

	tess := new(TessBaseAPI)
	tess.TBA = handle

	return tess, nil
}

/*
https://github.com/golang/go/issues/13467
I am slightly more sympathetic to *C.char,
but even there I don't understand why the package API doesn't just use
appropriate Go types instead (like []byte).
*/
func (t *TessBaseAPI) SetImage2(imgPath string) error {
	// pix := leptonica.NewPIX()
	img, err := leptonica.PixRead(imgPath)
	if err != nil {
		return fmt.Errorf("couldn't set image: %w", err)
	}
	defer leptonica.PixDestroy(img)
	// cImagePath := C.CString(imgPath)
	// defer C.free(unsafe.Pointer(cImagePath))
	// img := C.pixRead(cImagePath)
	// defer C.pixDestroy(&img)
	// if img == nil {
	// 	return fmt.Errorf("error reading image")
	// }
	// C.TessBaseAPISetImage2(t.TBA, (*C.PIX)(pix.Pix))
	C.TessBaseAPISetImage2(t.TBA, (*C.PIX)(unsafe.Pointer(img)))

	return nil
}

func (t *TessBaseAPI) SetImage2FromMem(data []byte) error {
	// pix := leptonica.NewPIX()
	img, err := leptonica.PixReadMem(data)
	if err != nil {
		return fmt.Errorf("couldn't set image: %w", err)
	}
	defer leptonica.PixDestroy(img)
	// cImagePath := C.CString(imgPath)
	// defer C.free(unsafe.Pointer(cImagePath))
	// img := C.pixRead(cImagePath)
	// defer C.pixDestroy(&img)
	// if img == nil {
	// 	return fmt.Errorf("error reading image")
	// }
	// C.TessBaseAPISetImage2(t.TBA, (*C.PIX)(pix.Pix))
	C.TessBaseAPISetImage2(t.TBA, (*C.PIX)(unsafe.Pointer(img)))

	return nil
}

func (t *TessBaseAPI) Recognize() error {
	if C.TessBaseAPIRecognize(t.TBA, nil) != 0 {
		return fmt.Errorf("error in tesseract recognition")
	}

	return nil
}

func (t *TessBaseAPI) GetUTF8Text() (string, error) {
	text := C.TessBaseAPIGetUTF8Text(t.TBA)
	defer C.free(unsafe.Pointer(text))
	if text == nil {
		return "", fmt.Errorf("error getting text")
	}

	return C.GoString(text), nil
}

func (t *TessBaseAPI) End() {
	C.TessBaseAPIEnd(t.TBA)
}

func (t *TessBaseAPI) Delete() {
	C.TessBaseAPIDelete(t.TBA)
}

// GetLoadedLanguagesAsVector returns loaded languages as string
// https://github.com/golang/go/wiki/cgo#turning-c-arrays-into-go-slices
func (t *TessBaseAPI) GetLoadedLanguagesAsVector() string {
	cArray := C.TessBaseAPIGetLoadedLanguagesAsVector(t.TBA)
	fmt.Printf("cArray: %v\n", cArray)
	length := unsafe.Sizeof(cArray)
	/*
		above 1.17 we can use unsafe.Slice
		goSlice := unsafe.Slice(cArray, length)
		loadedLanguages := C.GoString(goSlice[0])
	*/
	loadedLanguages := C.GoBytes(unsafe.Pointer(*cArray), C.int(length))
	return string(loadedLanguages)
}

func (t *TessBaseAPI) GetAvailableLanguagesAsVector() string {
	cArray := C.TessBaseAPIGetAvailableLanguagesAsVector(t.TBA)
	fmt.Printf("cArray: %v\n", cArray)
	length := unsafe.Sizeof(cArray)
	/*
		above 1.17 we can use unsafe.Slice
		goSlice := unsafe.Slice(cArray, length)
		loadedLanguages := C.GoString(goSlice[0])
	*/
	avalableLangs := C.GoBytes(unsafe.Pointer(*cArray), C.int(length))
	return string(avalableLangs)
}

func (t *TessBaseAPI) GetInitLanguagesAsString() string {
	cString := C.TessBaseAPIGetInitLanguagesAsString(t.TBA)
	return C.GoString(cString)
}

func (t *TessBaseAPI) SetVariable(name, value string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	res := C.TessBaseAPISetVariable(t.TBA, cName, cValue)
	if res == 0 {
		return fmt.Errorf("error setting variables")
	}
	return nil
}

type Box struct {
	x        int32 /*!< left coordinate                   */
	y        int32 /*!< top coordinate                    */
	w        int32 /*!< box width                         */
	h        int32 /*!< box height                        */
	refcount int32 /*!< reference count (1 if no clones)  */
}

func cropImage() {
	box := C.boxCreate(C.int(20), C.int(20), C.int(20), C.int(20))
	defer C.boxDestroy(&box)
	// croppedImage := C.pixClipRectangle()

}
