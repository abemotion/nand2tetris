package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const CommandTypeArithmetic = "C_ARITHMETIC"
const CommandTypePush = "C_PUSH"
const CommandTypePop = "C_POP"
const CommandTypeLabel = "C_LABEL"
const CommandTypeGoto = "C_GOTO"
const CommandTypeIf = "C_IF"
const CommandTypeFunction = "C_FUNCTION"
const CommandTypeReturn = "C_RETURN"
const CommandTypeCall = "C_CALL"

type Command struct {
	current string
	_type   string
	arg1    string
	arg2    string
	rom     int
	symbol  string
	dest    string
	comp    string
	jump    string
}

func Parse(inputFile *os.File, outputFile *os.File) {
	scanner := bufio.NewScanner(inputFile)
	fileName := inputFile.Name()[strings.LastIndex(inputFile.Name(), "/")+1 : strings.Index(inputFile.Name(), ".")]
	command := &Command{}

	for hasMoreCommands(scanner) {
		// Skip comment
		scannerText := scanner.Text()
		if strings.Contains(scannerText, "//") && strings.Index(scannerText, "//") == 0 {
			continue
		}

		// Skip blank line
		if !regexp.MustCompile(`.`).MatchString(scannerText) {
			continue
		}
		fmt.Println(scannerText)
		command.current = advance(scannerText)
		command._type = commandType(command)

		if command._type != CommandTypeReturn {
			command.arg1 = arg1(command)
		}

		if command._type == CommandTypePush || command._type == CommandTypePop ||
			command._type == CommandTypeFunction || command._type == CommandTypeCall {
			command.arg2 = arg2(command)
		}

		if command._type == CommandTypeArithmetic {
			writeArithmetic(command, outputFile)
		}

		if command._type == CommandTypePush || command._type == CommandTypePop {
			writePushPop(command._type, command.arg1, command.arg2, fileName, outputFile)
		}
		command.rom++
	}
}

func hasMoreCommands(scanner *bufio.Scanner) bool {
	return scanner.Scan()
}

func advance(scannerText string) string {
	s := strings.TrimSpace(scannerText)
	return s
}

func commandType(command *Command) string {
	if strings.Contains(command.current, "push") {
		return CommandTypePush
	}

	if strings.Contains(command.current, "pop") {
		return CommandTypePop
	}

	return CommandTypeArithmetic
}

func arg1(command *Command) string {
	if command._type == CommandTypeArithmetic {
		return command.current
	}
	fmt.Println(command)
	return command.current[strings.Index(command.current, " ")+1 : strings.LastIndex(command.current, " ")]
}

func arg2(command *Command) string {
	return command.current[strings.LastIndex(command.current, " ")+1:]
}
