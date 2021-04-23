package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	inputFilePath := flag.String("i", "", "input file path for parsing")
	outputFilePath := flag.String("o", "", "output file path for parsing")

	flag.Parse()

	inputFile, err := os.Open(*inputFilePath)
	if err != nil {
		log.Fatalf("input file error : %v", err)
	}
	defer inputFile.Close()

	outputFile, err := os.OpenFile(*outputFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("output file error : %v", err)
	}
	defer outputFile.Close()

	Parse(inputFile, outputFile)
}
