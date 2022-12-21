package main

import (
	"github.com/lillianhealth/jsplit/pkg/jsplit"

	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		filename   string
		outputPath string
		overwrite  bool
		err        error
	)

	flag.StringVar(&filename, "file", "", "Source JSON file")
	flag.StringVar(&outputPath, "output", "", "Output path for parsed JSON files (can be an s3:// or gs:// URI")
	flag.BoolVar(&overwrite, "overwrite", false, "Overwrite local filesystem output path if it exists")
	flag.Parse()

	if filename == "" || outputPath == "" {
		fmt.Println("Usage: jsplit -file <json_file> -output <output_path>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	err = jsplit.Split(filename, outputPath, overwrite)
	if err != nil {
		fmt.Printf("Split failed: %s", err)
		os.Exit(1)
	}
}
