package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	dirPath := flag.String("d", "", "Directory path for parsing")
	// outputFilePath := flag.String("o", "", "output file path for parsing")

	flag.Parse()

	files, err := ioutil.ReadDir(*dirPath)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fmt.Println(file.Name())

		if !strings.Contains(file.Name(), "jack") {
			continue
		}

		inputFile, err := os.Open(*dirPath + file.Name())
		if err != nil {
			log.Fatalf("input file error : %v", err)
		}
		defer inputFile.Close()

		outputFile, err := os.OpenFile(file.Name()[:strings.Index(file.Name(), ".")]+".xml", os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatalf("output file error : %v", err)
		}
		defer outputFile.Close()

		JackTokenize(inputFile, outputFile)
	}
}
