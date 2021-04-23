package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const SegmentConstant = "constant"
const SegmentArg = "argument"
const SegmentLocal = "local"
const SegmentStatic = "static"
const SegmentThis = "this"
const SegmentThat = "that"
const SegmentPointer = "pointer"
const SegmentTemp = "temp"

var symbols = map[string]string{
	"local":    "LCL",
	"argument": "ARG",
	"this":     "THIS",
	"that":     "THAT",
	"pointer":  "3",
	"temp":     "5",
	"static":   "",
}

func writePush(segment string, index string, outputFile *os.File) {
	if segment == SegmentConstant {
		if strings.Contains(index, "-") {
			fmt.Fprintln(outputFile, "push constant "+index[1:])
			fmt.Fprintln(outputFile, "neg")
		} else {
			fmt.Fprintln(outputFile, "push constant "+index)
		}

		return
	}

	fmt.Fprintln(outputFile, "push "+segment+" "+index)

}

func writePop(segment string, index string, outputFile *os.File) {
	if segment == SymbolKindStatic {
		fmt.Fprintln(outputFile, "pop static "+index)
		return
	}
	if segment == SymbolKindField {
		fmt.Fprintln(outputFile, "pop this "+index)
		return
	}
	if segment == SymbolKindArg {
		fmt.Fprintln(outputFile, "pop argument "+index)
		return
	}
	if segment == SymbolKindVar {
		fmt.Fprintln(outputFile, "pop local "+index)
		return
	}

	fmt.Fprintln(outputFile, "pop "+segment+" "+index)
}

func writeArithmetic(command string, outputFile *os.File) {
	// TODO: /
	// fmt.Fprintln(outputFile, command)
	if command == "+" {
		fmt.Fprintln(outputFile, "add")
	} else if command == "-" {
		fmt.Fprintln(outputFile, "sub")
	} else if command == "*" {
		fmt.Fprintln(outputFile, "call Math.multiply 2")
	} else if command == "/" {
		fmt.Fprintln(outputFile, "call Math.divide 2")
	} else if command == "&" {
		fmt.Fprintln(outputFile, "and")
	} else if command == "|" {
		fmt.Fprintln(outputFile, "or")
	} else if command == "<" {
		fmt.Fprintln(outputFile, "lt")
	} else if command == ">" {
		fmt.Fprintln(outputFile, "gt")
	} else if command == "=" {
		fmt.Fprintln(outputFile, "eq")
	} else if command == "~" {
		fmt.Fprintln(outputFile, "not")
	}
	// fmt.Fprintln(outputFile, "push "+name)
}

func writeLabel(label string, outputFile *os.File) {
	fmt.Fprintln(outputFile, "label "+label)
}

func writeGoto(label string, outputFile *os.File) {
	fmt.Fprintln(outputFile, "goto "+label)
}

func writeIf(label string, outputFile *os.File) {
	fmt.Fprintln(outputFile, "if-goto "+label)
}

func writeCall(name string, nArgs int, outputFile *os.File) {
	fmt.Fprintln(outputFile, "call "+name+" "+strconv.Itoa(nArgs))
}

func writeFunction(name string, nLocals int, outputFile *os.File) {
	fmt.Fprintln(outputFile, "function "+name+" "+strconv.Itoa(nLocals))
}

func writeReturn(outputFile *os.File) {
	fmt.Fprintln(outputFile, "return")
}
