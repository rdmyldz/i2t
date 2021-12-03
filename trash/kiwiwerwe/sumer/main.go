package main

// #cgo CFLAGS: -Wall -I${SRCDIR}/../
// #cgo LDFLAGS: -L${SRCDIR}/../ -lsumer_c -Wl,-rpath,${SRCDIR}/../
// #include "sum_array.h"
import "C"
import "fmt"

func main() {
	// gcc -o libperson.so -Wall -g -shared -fPIC person.c
	// gcc -o libsumer_c.so -fPIC -shared -Wall sum_array.c

	size := 5
	numbers := make([]int32, size)
	for i := 0; i < size; i++ {
		numbers[i] = int32(i)
	}

	res := C.sum_array(C.uint(size), (*C.int)(&numbers[0]))
	fmt.Printf("res: %v\n", res)
}
