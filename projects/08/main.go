package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	firstInputFilePath := flag.String("firsti", "", "input file path for parsing")
	secondInputFilePath := flag.String("secondi", "", "input file path for parsing")
	thirdInputFilePath := flag.String("thirdi", "", "input file path for parsing")
	outputFilePath := flag.String("o", "", "output file path for parsing")

	flag.Parse()

	// TODO: read files from input directory
	// TODO: lt

	inputFile, err := os.Open(*firstInputFilePath)
	if err != nil {
		log.Fatalf("input file error : %v", err)
	}
	defer inputFile.Close()

	secondInputFile, err := os.Open(*secondInputFilePath)
	if err != nil {
		log.Fatalf("input file error : %v", err)
	}
	defer secondInputFile.Close()

	thirdInputFile, err := os.Open(*thirdInputFilePath)
	if err != nil {
		log.Fatalf("input file error : %v", err)
	}
	defer thirdInputFile.Close()

	outputFile, err := os.OpenFile(*outputFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("output file error : %v", err)
	}
	defer outputFile.Close()

	fmt.Fprintln(outputFile, "@261")
	fmt.Fprintln(outputFile, " D=A")
	fmt.Fprintln(outputFile, "@SP")
	fmt.Fprintln(outputFile, " M=D")
	fmt.Fprintln(outputFile, "@Sys.init")
	fmt.Fprintln(outputFile, "0;JMP")
	Parse(inputFile, outputFile)
	Parse(secondInputFile, outputFile)
	Parse(thirdInputFile, outputFile)
}
