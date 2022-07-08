package main

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type ListAddFunc func(item []byte) error

func errExit(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) != 2 {
		errExit(fmt.Errorf("usage: jsplit <file>"))
	}

	filename := os.Args[1]
	dir := strings.Replace(filename, ".", "_", -1)

	if _, err := os.Stat(dir); err == nil {
		errExit(fmt.Errorf("error: %s already exists", filename))
	} else if !os.IsNotExist(err) {
		errExit(err)
	}

	err := os.Mkdir(dir, os.ModePerm)
	errExit(err)

	rd, err := AsyncReaderFromFile(filename, 1024*1024)
	errExit(err)

	fmt.Printf("Reading %s\n", filename)
	ctx := context.Background()
	ctx = rd.Start(ctx)

	err = SplitStream(ctx, rd, dir)
	errExit(err)
}
