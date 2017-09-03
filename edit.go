package libedit

import (
	"errors"
	"io"
	"unsafe"
)

// OSX bundles libedit, but places headers in /usr/include/readline directly.
// FreeBSD bundles libedit, but places headers in /usr/include/edit/readline.

// #cgo darwin LDFLAGS: -ledit
// #cgo darwin CPPFLAGS: -I/usr/include/readline
// #cgo freebsd LDFLAGS: -ledit
// #cgo freebsd CPPFLAGS: -I/usr/include/edit/readline
//
// #include <readline.h>
// #include <stdio.h>
// #include <stdlib.h>
import "C"

func Initialize() error {
	r := C.rl_initialize()
	if r < 0 {
		return errors.New("unable to initialize libedit")
	}
	return nil
}

var cPrompt *C.char

func Cleanup() {
	if cPrompt != nil {
		C.free(unsafe.Pointer(cPrompt))
	}
}

func SetPrompt(prompt string) {
	if cPrompt != nil {
		C.free(unsafe.Pointer(cPrompt))
	}
	cPrompt = C.CString(prompt)
}

func Readline() (string, error) {
	l := C.readline(cPrompt)
	if l == nil {
		return "", io.EOF
	}
	C.add_history(l)
	defer C.free(unsafe.Pointer(l))
	return C.GoString(l), nil
}
