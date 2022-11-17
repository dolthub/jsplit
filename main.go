package main

import (
	"context"
	"flag"
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
	var filename string
	var outputPath string

	flag.StringVar(&filename, "file", "", "Source JSON file")
	flag.StringVar(&outputPath, "output", "", "Output path for parsed JSON files (optional)")
	flag.Parse()

	if len(filename) == 0 {
		fmt.Println("Usage: jsplit -file <json_file> -output <output_path>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(outputPath) == 0 {
		outputPath = strings.Replace(filename, ".", "_", -1)
	}

	if _, err := os.Stat(outputPath); err == nil {
		errExit(fmt.Errorf("error: %s already exists", filename))
	} else if !os.IsNotExist(err) {
		errExit(err)
	}

	err := os.Mkdir(outputPath, os.ModePerm)
	errExit(err)

	rd, err := AsyncReaderFromFile(filename, 1024*1024)
	errExit(err)

	fmt.Printf("Reading %s\n", filename)
	ctx := context.Background()
	ctx = rd.Start(ctx)

	err = SplitStream(ctx, rd, outputPath)
	errExit(err)
}
