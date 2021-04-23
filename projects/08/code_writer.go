package main

import (
	"fmt"
	"os"
	"strconv"
)

var symbols = map[string]string{
	"local":    "LCL",
	"argument": "ARG",
	"this":     "THIS",
	"that":     "THAT",
	"pointer":  "3",
	"temp":     "5",
	"static":   "",
}

func writeInit(outputFile *os.File) {
	//fmt.Fprintln(outputFile, "@SP")
	//fmt.Fprintln(outputFile, "M=256")
	fmt.Fprintln(outputFile, "goto Sys.init")
}

func writeArithmetic(command *Command, outputFile *os.File) {
	// true=-1
	// false=0
	// & ビット積(AND)。オペランドは整数に限られ、boolでは使えない。
	// | ビット和(OR)。オペランドは整数に限られ、boolでは使えない。
	// ^ ビットごとの否定(NOT)。符号は反転される。オペランドは整数に限られ、boolでは使えない。

	if command.arg1 == "neg" || command.arg1 == "not" {
		fmt.Fprintln(outputFile, "@SP")
		fmt.Fprintln(outputFile, "A=M-1")

		if command.arg1 == "neg" {
			fmt.Fprintln(outputFile, "M=-M")
		}
		if command.arg1 == "not" {
			fmt.Fprintln(outputFile, "M=!M")
		}
		return
	}

	fmt.Fprintln(outputFile, "@SP")
	fmt.Fprintln(outputFile, "AM=M-1")
	fmt.Fprintln(outputFile, "D=M")
	fmt.Fprintln(outputFile, "M=0")
	fmt.Fprintln(outputFile, "@SP")
	fmt.Fprintln(outputFile, "AM=M-1")

	if command.arg1 == "eq" {
		fmt.Fprintln(outputFile, "D=M-D")
		fmt.Fprintln(outputFile, "M=0")
		fmt.Fprintln(outputFile, "@NE"+strconv.Itoa(command.rom))
		fmt.Fprintln(outputFile, "D;JNE")
		fmt.Fprintln(outputFile, "@SP")
		fmt.Fprintln(outputFile, "A=M")
		fmt.Fprintln(outputFile, "M=-1")
		fmt.Fprintln(outputFile, "(NE"+strconv.Itoa(command.rom)+")")
	}

	if command.arg1 == "gt" {
		fmt.Fprintln(outputFile, "D=M-D")
		fmt.Fprintln(outputFile, "M=0")
		fmt.Fprintln(outputFile, "@LE"+strconv.Itoa(command.rom))
		fmt.Fprintln(outputFile, "D;JLE")
		fmt.Fprintln(outputFile, "@SP")
		fmt.Fprintln(outputFile, "A=M")
		fmt.Fprintln(outputFile, "M=-1")
		fmt.Fprintln(outputFile, "(LE"+strconv.Itoa(command.rom)+")")
	}

	if command.arg1 == "lt" {
		fmt.Fprintln(outputFile, "D=M-D")
		fmt.Fprintln(outputFile, "M=0")
		fmt.Fprintln(outputFile, "@GE"+strconv.Itoa(command.rom))
		fmt.Fprintln(outputFile, "D;JGE")
		fmt.Fprintln(outputFile, "@SP")
		fmt.Fprintln(outputFile, "A=M")
		fmt.Fprintln(outputFile, "M=-1")
		fmt.Fprintln(outputFile, "(GE"+strconv.Itoa(command.rom)+")")
	}

	if command.arg1 == "and" {
		fmt.Fprintln(outputFile, "D=M&D")
		fmt.Fprintln(outputFile, "M=0")
		fmt.Fprintln(outputFile, "@SP")
		fmt.Fprintln(outputFile, "A=M")
		fmt.Fprintln(outputFile, "M=D")
	}

	if command.arg1 == "or" {
		fmt.Fprintln(outputFile, "D=M|D")
		fmt.Fprintln(outputFile, "M=0")
		fmt.Fprintln(outputFile, "@SP")
		fmt.Fprintln(outputFile, "A=M")
		fmt.Fprintln(outputFile, "M=D")
	}

	if command.arg1 == "add" {
		fmt.Fprintln(outputFile, "D=M+D")
		fmt.Fprintln(outputFile, "M=0")
		fmt.Fprintln(outputFile, "@SP")
		fmt.Fprintln(outputFile, "A=M")
		fmt.Fprintln(outputFile, "M=D")
	}

	if command.arg1 == "sub" {
		fmt.Fprintln(outputFile, "D=M-D")
		fmt.Fprintln(outputFile, "M=0")
		fmt.Fprintln(outputFile, "@SP")
		fmt.Fprintln(outputFile, "A=M")
		fmt.Fprintln(outputFile, "M=D")
	}

	fmt.Fprintln(outputFile, "@SP")
	fmt.Fprintln(outputFile, "M=M+1")
}

func writePushPop(commandType string, segment string, index string, inputFileName string, outputFile *os.File) {
	v, ok := symbols[segment]
	segmentAssembly := "@"
	intIndex, _ := strconv.Atoi(index)

	if ok {
		if segment == "pointer" {
			segmentAssembly += strconv.Itoa(3 + intIndex)
		} else if segment == "temp" {
			segmentAssembly += strconv.Itoa(5 + intIndex)
		} else if segment == "static" {
			segmentAssembly += inputFileName + "." + index
		} else {
			segmentAssembly += v
		}
	} else {
		segmentAssembly += index
	}

	if commandType == CommandTypePush {
		fmt.Println(segmentAssembly)
		fmt.Fprintln(outputFile, segmentAssembly)

		if ok {
			if index != "0" && segment != "pointer" && segment != "temp" && segment != "static" {
				fmt.Fprintln(outputFile, "A=M+1")
				if (intIndex - 1) > 0 {
					for i := 0; i < intIndex-1; i++ {
						fmt.Fprintln(outputFile, "A=A+1")
					}
				}
				fmt.Fprintln(outputFile, "D=M")
			} else if segment == "pointer" || segment == "temp" || segment == "static" {
				fmt.Fprintln(outputFile, "D=M")
			} else {
				fmt.Fprintln(outputFile, "A=M")
				fmt.Fprintln(outputFile, "D=M")
			}
		} else {
			fmt.Fprintln(outputFile, "D=A")
		}

		fmt.Fprintln(outputFile, "@SP")
		fmt.Fprintln(outputFile, "A=M")
		fmt.Fprintln(outputFile, "M=D")
		fmt.Fprintln(outputFile, "@SP")
		fmt.Fprintln(outputFile, "M=M+1")
	}

	if commandType == CommandTypePop {
		fmt.Fprintln(outputFile, "@SP")
		fmt.Fprintln(outputFile, "AM=M-1")
		fmt.Fprintln(outputFile, "D=M")
		fmt.Fprintln(outputFile, "M=0")
		fmt.Fprintln(outputFile, segmentAssembly)
		if index != "0" && segment != "pointer" && segment != "temp" && segment != "static" {
			fmt.Fprintln(outputFile, "A=M+1")
			if (intIndex - 1) > 0 {
				for i := 0; i < intIndex-1; i++ {
					fmt.Fprintln(outputFile, "A=A+1")
				}
			}
			fmt.Fprintln(outputFile, "M=D")
		} else if segment == "pointer" || segment == "temp" || segment == "static" {
			fmt.Fprintln(outputFile, "M=D")
		} else {
			fmt.Fprintln(outputFile, "A=M")
			fmt.Fprintln(outputFile, "M=D")
		}
	}
}

func writeLabel(command *Command, outputFile *os.File) {
	fmt.Fprintln(outputFile, "("+command.arg1+")")
}

func writeGoto(command *Command, outputFile *os.File) {
	fmt.Fprintln(outputFile, "@"+command.arg1)
	fmt.Fprintln(outputFile, "0;JMP")
}

func writeIf(command *Command, outputFile *os.File) {
	fmt.Fprintln(outputFile, "@SP")
	fmt.Fprintln(outputFile, "AM=M-1")
	fmt.Fprintln(outputFile, "D=M")
	fmt.Fprintln(outputFile, "M=0")
	fmt.Fprintln(outputFile, "@"+command.arg1)
	fmt.Fprintln(outputFile, "D;JLT")
}

func writeCall(command *Command, outputFile *os.File) {
	// fmt.Fprintln(outputFile, "call "+functionName+" "+strconv.Itoa(numArgs))
	//fmt.Fprintln(outputFile, "@LE"+)
	// fmt.Fprintln(outputFile, "@LCL")
	// fmt.Fprintln(outputFile, "A=M")
	// fmt.Fprintln(outputFile, "D=M")
	// fmt.Fprintln(outputFile, "@SP")
	// fmt.Fprintln(outputFile, "A=M")
	// fmt.Fprintln(outputFile, "M=D")
	// fmt.Fprintln(outputFile, "@SP")
	// fmt.Fprintln(outputFile, "M=M+1")

	fmt.Fprintln(outputFile, "@return-address"+strconv.Itoa(command.rom))
	fmt.Fprintln(outputFile, "D=A")
	fmt.Fprintln(outputFile, "@SP")
	fmt.Fprintln(outputFile, "A=M")
	fmt.Fprintln(outputFile, "M=D")
	fmt.Fprintln(outputFile, "@SP")
	fmt.Fprintln(outputFile, "M=M+1")

	writePush("LCL", outputFile)
	writePush("ARG", outputFile)
	writePush("THIS", outputFile)
	writePush("THAT", outputFile)

	argIndex, _ := strconv.Atoi(command.arg2)

	fmt.Fprintln(outputFile, "@SP")
	fmt.Fprintln(outputFile, "D=M")
	if argIndex > 0 {
		fmt.Fprintln(outputFile, "D=M-1")
	}

	if (argIndex - 1) > 0 {
		for i := 0; i < argIndex-1; i++ {
			fmt.Fprintln(outputFile, "D=D-1")
		}
	}

	fmt.Fprintln(outputFile, "D=D-1")
	fmt.Fprintln(outputFile, "D=D-1")
	fmt.Fprintln(outputFile, "D=D-1")
	fmt.Fprintln(outputFile, "D=D-1")
	fmt.Fprintln(outputFile, "D=D-1")
	fmt.Fprintln(outputFile, "@ARG")
	fmt.Fprintln(outputFile, "M=D")

	fmt.Fprintln(outputFile, "@SP")
	fmt.Fprintln(outputFile, "D=M")
	fmt.Fprintln(outputFile, "@LCL")
	fmt.Fprintln(outputFile, "M=D")

	fmt.Fprintln(outputFile, "@"+command.arg1)
	fmt.Fprintln(outputFile, "0;JMP")

	fmt.Fprintln(outputFile, "(return-address"+strconv.Itoa(command.rom)+")")
}

func writePush(input string, outputFile *os.File) {
	fmt.Fprintln(outputFile, "@"+input)
	fmt.Fprintln(outputFile, "D=M")
	fmt.Fprintln(outputFile, "@SP")
	fmt.Fprintln(outputFile, "A=M")
	fmt.Fprintln(outputFile, "M=D")
	fmt.Fprintln(outputFile, "@SP")
	fmt.Fprintln(outputFile, "M=M+1")
}

func writeReturn(command *Command, outputFile *os.File) {

	// fmt.Fprintln(outputFile, "goto D-5")

	fmt.Fprintln(outputFile, "@LCL")
	fmt.Fprintln(outputFile, "D=M")
	fmt.Fprintln(outputFile, "@R13")
	fmt.Fprintln(outputFile, "M=D")
	fmt.Fprintln(outputFile, "D=D-1")
	fmt.Fprintln(outputFile, "D=D-1")
	fmt.Fprintln(outputFile, "D=D-1")
	fmt.Fprintln(outputFile, "D=D-1")
	fmt.Fprintln(outputFile, "A=D-1")
	fmt.Fprintln(outputFile, "D=M")
	fmt.Fprintln(outputFile, "@R14")
	fmt.Fprintln(outputFile, "M=D")

	fmt.Fprintln(outputFile, "@SP")
	fmt.Fprintln(outputFile, "AM=M-1")
	fmt.Fprintln(outputFile, "D=M")
	fmt.Fprintln(outputFile, "@ARG")
	fmt.Fprintln(outputFile, "A=M")
	fmt.Fprintln(outputFile, "M=D")

	fmt.Fprintln(outputFile, "@ARG")
	fmt.Fprintln(outputFile, "D=M+1")
	fmt.Fprintln(outputFile, "@SP")
	fmt.Fprintln(outputFile, "M=D")

	fmt.Fprintln(outputFile, "@R13")
	fmt.Fprintln(outputFile, "D=M")
	fmt.Fprintln(outputFile, "A=D-1")
	fmt.Fprintln(outputFile, "D=M")
	fmt.Fprintln(outputFile, "@THAT")
	fmt.Fprintln(outputFile, "M=D")

	fmt.Fprintln(outputFile, "@R13")
	fmt.Fprintln(outputFile, "D=M")
	fmt.Fprintln(outputFile, "D=D-1")
	fmt.Fprintln(outputFile, "A=D-1")
	fmt.Fprintln(outputFile, "D=M")
	fmt.Fprintln(outputFile, "@THIS")
	fmt.Fprintln(outputFile, "M=D")

	fmt.Fprintln(outputFile, "@R13")
	fmt.Fprintln(outputFile, "D=M")
	fmt.Fprintln(outputFile, "D=D-1")
	fmt.Fprintln(outputFile, "D=D-1")
	fmt.Fprintln(outputFile, "A=D-1")
	fmt.Fprintln(outputFile, "D=M")
	fmt.Fprintln(outputFile, "@ARG")
	fmt.Fprintln(outputFile, "M=D")

	fmt.Fprintln(outputFile, "@R13")
	fmt.Fprintln(outputFile, "D=M")
	fmt.Fprintln(outputFile, "D=D-1")
	fmt.Fprintln(outputFile, "D=D-1")
	fmt.Fprintln(outputFile, "D=D-1")
	fmt.Fprintln(outputFile, "A=D-1")
	fmt.Fprintln(outputFile, "D=M")
	fmt.Fprintln(outputFile, "@LCL")
	fmt.Fprintln(outputFile, "M=D")

	fmt.Fprintln(outputFile, "@R14")
	fmt.Fprintln(outputFile, "A=M")
	fmt.Fprintln(outputFile, "0;JMP")
}

func writeFunction(command *Command, outputFile *os.File) {
	fmt.Fprintln(outputFile, "("+command.arg1+")")

	intLocal, _ := strconv.Atoi(command.arg2)
	for i := 0; i < intLocal; i++ {
		fmt.Fprintln(outputFile, "@0")
		fmt.Fprintln(outputFile, "D=A")
		fmt.Fprintln(outputFile, "@SP")
		fmt.Fprintln(outputFile, "A=M")
		fmt.Fprintln(outputFile, "M=D")
		fmt.Fprintln(outputFile, "@SP")
		fmt.Fprintln(outputFile, "M=M+1")
	}
}
