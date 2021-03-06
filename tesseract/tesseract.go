package tesseract

// #cgo CFLAGS: -g
// #cgo pkg-config: tesseract
// #cgo LDFLAGS: -llept
// #include <stdbool.h>
// #include <tesseract/capi.h>
// #include <leptonica/allheaders.h>
/*
bool gocb(void *ctx, int a);
*/
import "C"
import (
	"fmt"
	"log"
	"os"
	"unsafe"

	"github.com/rdmyldz/i2t/leptonica"
)

//export gocb
func gocb(ctx unsafe.Pointer, words C.int) C.bool {
	fmt.Println("callback invoked")
	return true
}

func (t *TessBaseAPI) SetCancelFunc() {
	C.TessMonitorSetCancelFunc(
		t.Mon,
		(C.TessCancelFunc)(unsafe.Pointer(C.gocb)),
	)
}

type TessBaseAPI struct {
	TBA *C.TessBaseAPI
	Mon *C.ETEXT_DESC
}

func TessVersion() string {
	vers := C.TessVersion()
	return C.GoString(vers)
}

func TessBaseAPICreate(language string) (*TessBaseAPI, error) {
	handle := C.TessBaseAPICreate()

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

// TessBaseAPICreateWithMonitor use this if you wanna cancel
func TessBaseAPICreateWithMonitor(language string) (*TessBaseAPI, error) {
	handle := C.TessBaseAPICreate()
	mon := C.TessMonitorCreate()

	cLanguage := C.CString(language)
	defer C.free(unsafe.Pointer(cLanguage))

	res := C.TessBaseAPIInit3(handle, nil, cLanguage)
	if res != 0 {
		return nil, fmt.Errorf("error initializing tesseract")
	}

	tess := new(TessBaseAPI)
	tess.TBA = handle
	tess.Mon = mon

	return tess, nil
}

/*
https://tpgit.github.io/Leptonica/readfile_8c_source.html#l00285
 pixReadStream()

      Input:  fp (file stream)
              hint (bitwise OR of L_HINT_* values for jpeg; use 0 for no hint)
      Return: pix if OK; null on error

  Notes:
      (1) The hint only applies to jpeg.
*/

func (t *TessBaseAPI) ProcessImage(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("in ProcessImage, error reading file: %w", err)
	}

	return t.ProcessImageMem(data)
}

func (t *TessBaseAPI) getText(pixes []leptonica.Pix) ([]string, error) {
	var texts []string
	for _, pix := range pixes {
		t.SetImage2((*C.PIX)(unsafe.Pointer(pix)))
		err := t.Recognize()
		if err != nil {
			return nil, err
		}
		str, err := t.GetUTF8Text()
		if err != nil {
			return nil, err
		}
		texts = append(texts, str)

	}
	return texts, nil
}

func (t *TessBaseAPI) ProcessImageMem(data []byte) ([]string, error) {
	pixes, err := leptonica.PixReadMem(data)
	if err != nil {
		return nil, fmt.Errorf("in ProcessImageMem, %w", err)
	}
	defer leptonica.DestroyPixes(pixes)

	return t.getText(pixes)
}

/*
https://github.com/golang/go/issues/13467
I am slightly more sympathetic to *C.char,
but even there I don't understand why the package API doesn't just use
appropriate Go types instead (like []byte).
*/
func (t *TessBaseAPI) SetImage2(img *C.PIX) {
	C.TessBaseAPISetImage2(t.TBA, img)
}

func (t *TessBaseAPI) Recognize() error {
	if err := C.TessBaseAPIRecognize(t.TBA, t.Mon); err != 0 {
		return fmt.Errorf("error in tesseract recognition: %v", err)
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
	log.Println("End Called")
	C.TessBaseAPIEnd(t.TBA)
}

func (t *TessBaseAPI) Delete() {
	log.Println("Delete Called")
	if t.Mon != nil {
		C.TessMonitorDelete(t.Mon)
	}
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
	availableLangs := C.GoBytes(unsafe.Pointer(*cArray), C.int(length))
	return string(availableLangs)
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
