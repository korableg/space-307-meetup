package db

import (
	"fmt"
	"testing"
	"unsafe"
)

func Test_UnsafeValidOp(t *testing.T) {
	var k = []byte{1, 2, 3, 4, 5, 6}

	p := (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&k[0])) + 2))
	fmt.Println(*p)
}

func Test_UnsafeInvalidOp(t *testing.T) {
	var k = []byte{1, 2, 3, 4, 5, 6}

	uP := uintptr(unsafe.Pointer(&k[0]))
	p := (*byte)(unsafe.Pointer(uP + 2))
	fmt.Println(*p)
}

func Test_UnsafeValidOp2(t *testing.T) {
	var k = []byte{1, 2, 3, 4, 5, 6}

	uP := unsafe.Pointer(&k[0])
	p := (*byte)(unsafe.Add(uP, 2))
	fmt.Println(*p)
}
