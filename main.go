package main

// #cgo CFLAGS: -g -Wall
// #cgo pkg-config: tesseract
// #cgo LDFLAGS: -llept
// #include <stdio.h>
// #include <stdlib.h>
// #include "/usr/include/leptonica/allheaders.h"
// #include "/usr/include/tesseract/capi.h"
import "C"
import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rdmyldz/i2t/tesseract"
)

func main() {
	testdir := "testdata"
	ext := ".png"
	// testdir := "trashdata"
	// ext := ".tif"
	dir, err := os.ReadDir(testdir)
	if err != nil {
		log.Fatalln(err)
	}
	var pictures []string
	for _, f := range dir {
		if f.IsDir() {
			continue
		}

		if strings.HasSuffix(f.Name(), ext) {
			pictures = append(pictures, filepath.Join(testdir, f.Name()))
		}
	}
	fmt.Printf("pics: %v\n", pictures)
	vers := tesseract.TessVersion()
	fmt.Printf("verson: %v\n", vers)
	// handle := C.TessBaseAPICreate()
	// defer C.free(unsafe.Pointer(handle))

	handle, err := tesseract.TessBaseAPICreate("tur")
	if err != nil {
		log.Fatalln(err)
	}
	// err = handle.SetVariable("tessedit_char_whitelist", "123456789ABCDEFGHIJKLMNOPQRSTUVWXYZİĞÜŞÇÖĞ")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// loadedLanguages := handle.GetLoadedLanguagesAsVector()
	// fmt.Printf("loadedLanguages: %s\n", loadedLanguages)

	// availableLangs := handle.GetAvailableLanguagesAsVector()
	// fmt.Printf("avalableLangs: %s\n", availableLangs)

	// initLangs := handle.GetInitLanguagesAsString()
	// fmt.Printf("initLangs: %s\n", initLangs)

	// var datapath string
	// cdp := C.CString(datapath)
	// defer C.free(unsafe.Pointer(cdp))
	// r1Ctype := C.TessBaseAPIInit3(handle, nil, C.CString("eng"))

	// fmt.Printf("r1Ctype: %v\n", r1Ctype)
	// if r1Ctype != 0 {
	// 	log.Fatalln("error initializing tesseract")
	// }

	// imgPath := C.CString("qq1.png")
	// defer C.free(unsafe.Pointer(imgPath))

	// img := C.pixRead(imgPath)
	// if img == nil {
	// 	log.Fatalln("error reading image")
	// }

	// C.TessBaseAPISetImage2((*C.TessBaseAPI)(handle.TBA), img)
	// if C.TessBaseAPIRecognize() {

	// }
	// if C.TessBaseAPIRecognize((*C.TessBaseAPI)(handle.TBA), nil) != 0 {
	// 	log.Fatal("error in tesseract recogntion")
	// }
	// text := C.TessBaseAPIGetUTF8Text((*C.TessBaseAPI)(handle.TBA))
	// if text == nil {
	// 	log.Fatal("error getting text")
	// }
	/*
		// cropping code
		cImagePath := C.CString("e.png")
		defer C.free(unsafe.Pointer(cImagePath))
		img := C.pixRead(cImagePath)
		defer C.pixDestroy(&img)
		box := C.boxCreate(C.int(10), C.int(0), C.int(20), C.int(45))
		defer C.boxDestroy(&box)
		croppedImage := C.pixClipRectangle(img, box, nil)
		res := C.pixWrite(C.CString("e5.png"), croppedImage, C.int(3))
		if res != 0 {
			log.Fatalln("error writing to file")
		}
		return
	*/
	for _, pic := range pictures {

		// texts, err := handle.ProcessImage(pic)
		texts, err := handle.ProcessImage(pic)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("--------------%s----------------------\n", pic)
		for i, t := range texts {
			fmt.Printf("page %d\n", i+1)
			fmt.Println(t)
		}

	}

}
