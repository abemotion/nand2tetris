package main

import (
	"strings"
)

func Dest(n string) (r string) {
	if n == "" {
		return "000"
	}

	r = "0"
	if strings.Contains(n, "A") {
		r = "1"
	}

	if strings.Contains(n, "D") {
		r += "1"
	} else {
		r += "0"
	}

	if strings.Contains(n, "M") {
		r += "1"
	} else {
		r += "0"
	}

	return
}

func Comp(n string) (r string) {
	r = "0"
	if strings.Contains(n, "M") {
		r = "1"
	}

	if strings.Contains(n, "0") {
		r += "101010"
		return
	}

	if strings.Contains(n, "|") {
		r += "010101"
		return
	}

	if strings.Contains(n, "&") {
		r += "000000"
		return
	}

	if strings.Contains(n, "-D") {
		r += "000111"
		return
	}

	if strings.Contains(n, "D-") {
		r += "010011"
		return
	}

	if strings.Contains(n, "D+") {
		r += "000010"
		return
	}

	if strings.Contains(n, "D-1") {
		r += "001110"
		return
	}

	if strings.Contains(n, "A-1") || strings.Contains(n, "M-1") {
		r += "110010"
		return
	}

	if strings.Contains(n, "D+1") {
		r += "011111"
		return
	}

	if strings.Contains(n, "+1") {
		r += "110111"
		return
	}

	if strings.Contains(n, "-D") {
		r += "001111"
		return
	}

	if strings.Contains(n, "-A") || strings.Contains(n, "-M") {
		r += "110011"
		return
	}

	if strings.Contains(n, "!D") {
		r += "001101"
		return
	}

	if strings.Contains(n, "!") {
		r += "110001"
		return
	}

	if strings.Contains(n, "D") {
		r += "001100"
		return
	}

	if strings.Contains(n, "A") || strings.Contains(n, "M") {
		r += "110000"
		return
	}

	if strings.Contains(n, "-1") {
		r += "111010"
		return
	}

	r += "111111"
	return
}

func Jump(n string) (r string) {
	if n == "" {
		return "000"
	}

	if n == "JMP" {
		return "111"
	}

	if strings.Contains(n, "JG") || strings.Contains(n, "JE") {
		r = "0"
	} else {
		r = "1"
	}

	if strings.Contains(n, "T") || strings.Contains(n, "N") {
		r += "0"
	} else {
		r += "1"
	}

	if strings.Contains(n, "L") || strings.Contains(n, "Q") {
		r += "0"
	} else {
		r += "1"
	}

	return
}
