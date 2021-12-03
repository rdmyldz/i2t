package rand

// #include <stdio.h>
// #include <stdlib.h>
import "C"
import "unsafe"

func Random() int {
	return int(C.random())
}

func Seed(i int) {
	C.srandom(C.uint(i))

}

func Print(s string) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.fputs(cs, (*C.FILE)(C.stdout))
}
