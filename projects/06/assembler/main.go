package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"regexp"
	"strings"
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

	scanner := bufio.NewScanner(inputFile)
	table := NewSymbolTable()

	for scanner.Scan() {
		// Skip comment
		scannerText := scanner.Text()
		if strings.Contains(scannerText, "//") && strings.Index(scannerText, "//") == 0 {
			continue
		}

		// Skip blank line
		if !regexp.MustCompile(`.`).MatchString(scannerText) {
			continue
		}

		command := &Command{}
		command.current = advance(scannerText)
		command._type = commandType(command)
		command.symbol = symbol(command)

		if command._type == CommandTypeL {
			if !table.contains(command.symbol) {
				table.addEntry(command.symbol, table.rom)
			}
			continue
		}

		table.rom++
	}

	inputFileAgain, err := os.Open(*inputFilePath)
	if err != nil {
		log.Fatalf("input file error : %v", err)
	}
	defer inputFileAgain.Close()

	table.rom = 0
	Parse(inputFileAgain, outputFile, table)
}
