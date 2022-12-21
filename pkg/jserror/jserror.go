package jserror

import (
	"fmt"
	"os"
	"runtime"
)

func ErrExit(err error) {
	if err != nil {
		_, file, no, ok := runtime.Caller(1)
		if ok {
			fmt.Fprintf(os.Stderr, "Fatal error in %s#%d: %s", file, no, err.Error())
		} else {
			fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		}

		os.Exit(1)
	}
}
