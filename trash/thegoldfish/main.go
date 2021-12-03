package main

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -L./personlib -lperson -Wl,-rpath=./personlib
// #include "./personlib/person.h"
import "C"

import "fmt"

type Person C.struct_APerson

func GetPerson(name string, long_name string) *Person {
	return (*Person)(C.get_person(C.CString(name), C.CString(long_name)))
}

func main() {
	var p *Person
	p = GetPerson("rdm", "rdm yldz")
	fmt.Printf("Hello Go World! my name is %s, %s\n", C.GoString(p.name), C.GoString(p.long_name))

}
