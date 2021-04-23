package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const CommandTypeA = "A_COMMAND"
const CommandTypeC = "C_COMMAND"
const CommandTypeL = "L_COMMAND"

type Command struct {
	current string
	_type   string
	symbol  string
	dest    string
	comp    string
	jump    string
}

func Parse(inputFile *os.File, outputFile *os.File, table *SymbolTable) {
	scanner := bufio.NewScanner(inputFile)

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

		command := &Command{}
		command.current = advance(scannerText)
		command._type = commandType(command)

		if command._type == CommandTypeL {
			continue
		}

		if command._type == CommandTypeA {
			command.symbol = symbol(command)
			r := "0"
			var i int

			if regexp.MustCompile(`^[0-9]`).Match([]byte(command.symbol)) {
				i, _ = strconv.Atoi(command.current[1:])
			} else {
				if !table.contains(command.symbol) {
					table.addEntry(command.symbol, table.ram)
					table.ram++
				}

				i = table.getAddress(command.symbol)
			}

			r += fmt.Sprintf("%015b", i)
			fmt.Println("A")
			fmt.Println(r)
			fmt.Fprintln(outputFile, r)
			table.rom++

			continue
		}

		if command._type == CommandTypeC {
			command.dest = dest(command)
			command.comp = comp(command)
			command.jump = jump(command)
		}

		r := "111"
		r += Comp(command.comp)
		r += Dest(command.dest)
		r += Jump(command.jump)

		fmt.Println("C")
		fmt.Println(command.current)
		fmt.Println(command.comp)
		fmt.Println(r)
		fmt.Fprintln(outputFile, r)

		table.rom++
	}
}

func hasMoreCommands(scanner *bufio.Scanner) bool {
	return scanner.Scan()
}

func advance(scannerText string) string {
	s := strings.TrimSpace(scannerText)
	if strings.Contains(s, " ") {
		return s[:strings.Index(s, " ")]
	}
	return s
}

func commandType(command *Command) string {
	firstChar := command.current[0:1]

	if firstChar == "@" {
		return CommandTypeA
	}

	if firstChar == "(" {
		return CommandTypeL
	}

	return CommandTypeC
}

func symbol(command *Command) string {
	// be called when A_COMMAND or L_COMMAND

	if command._type == CommandTypeA {
		return command.current[1:]
	}

	return command.current[1 : len(command.current)-1]
}

func dest(command *Command) string {
	// be called when C_COMMAND

	if strings.Contains(command.current, "=") {
		return command.current[:strings.Index(command.current, "=")]
	}

	return ""
}

func comp(command *Command) string {
	// be called when C_COMMAND

	if strings.Contains(command.current, "=") {
		equalIndex := strings.Index(command.current, "=")

		if strings.Contains(command.current, ";") {
			return command.current[equalIndex+1 : strings.Index(command.current, ";")]
		}

		return command.current[equalIndex+1:]
	}

	if strings.Contains(command.current, ";") {
		return command.current[:strings.Index(command.current, ";")]
	}

	return ""
}

func jump(command *Command) string {
	// be called when C_COMMAND

	if strings.Contains(command.current, ";") {
		return command.current[strings.Index(command.current, ";")+1:]
	}

	return ""
}
