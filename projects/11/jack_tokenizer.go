package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const TokenTypeKeyword = "KEYWORD"
const TokenTypeSymbol = "SYMBOL"
const TokenTypeIntConst = "INT_CONST"
const TokenTypeStringConst = "STRING_CONST"
const TokenTypeIdentifer = "IDENTIFIER"

const KeywordClass = "CLASS"
const KeywordMethod = "METHOD"
const KeywordFunction = "FUNCTION"
const KeywordConstructor = "CONSTRUCTOR"
const KeywordInt = "INT"
const KeywordBoolean = "BOOLEAN"
const KeywordChar = "CHAR"
const KeywordVoid = "VOID"
const KeywordVar = "VAR"
const KeywordStatic = "STATIC"
const KeywordField = "FIELD"
const KeywordLet = "LET"
const KeywordDo = "DO"
const KeywordIf = "IF"
const KeywordElse = "ELSE"
const KeywordWhile = "WHILE"
const KeywordReturn = "RETURN"
const KeywordTrue = "TRUE"
const KeywordFalse = "FALSE"
const KeywordNull = "NULL"
const KeywordThis = "THIS"

func JackTokenize(inputFile *os.File, outputFile *os.File) {
	scanner := bufio.NewScanner(inputFile)

	var tokenList []string

	for hasMoreTokens(scanner) {
		scannerText := scanner.Text()

		// Skip comment
		if strings.Contains(scannerText, "//") && strings.Index(scannerText, "//") == 0 {
			continue
		}
		if strings.Contains(scannerText, "/*") {
			continue
		}
		if strings.Contains(scannerText, "*") && strings.Index(scannerText, "*") == 1 {
			continue
		}

		// Skip blank line
		if !regexp.MustCompile(`.`).MatchString(scannerText) {
			continue
		}

		tokenList = advance(scannerText, tokenList)
	}

	// fmt.Fprintln(outputFile, "<tokens>")
	// for _, token := range tokenList {
	// 	tokenToXml(token, outputFile)
	// }
	// fmt.Fprintln(outputFile, "</tokens>")

	fmt.Println(tokenList)
	CompileClass(tokenList, outputFile)
}

func hasMoreTokens(scanner *bufio.Scanner) bool {
	return scanner.Scan()
}

func advance(scannerText string, tokenList []string) []string {
	s := strings.TrimSpace(scannerText)
	// s := scannerText

	if strings.Contains(scannerText, "//") {
		s = scannerText[0:strings.Index(scannerText, "//")]
		s = strings.TrimSpace(s)
	}

	var token string
	var tokeni int

	for i, c := range s {
		if tokeni != 0 && i <= tokeni {
			continue
		}
		stringChar := string(c)
		// fmt.Println(string(c))

		if regexp.MustCompile(`[\{\}\(\)\[\]\.,;\+\-\*\/&\|<>=~]`).MatchString(stringChar) {
			if token != "" {
				tokenList = append(tokenList, token)
				token = ""
			}

			if stringChar == "-" {
				if s[i+1:i+2] == " " {
					tokenList = append(tokenList, stringChar)
				} else {
					token += stringChar
				}
			} else {
				tokenList = append(tokenList, stringChar)
			}

			continue
		}

		if stringChar == "\"" {
			tokenList = append(tokenList, s[i:strings.LastIndex(s, "\"")+1])
			token = ""
			tokeni = strings.LastIndex(s, "\"")
			continue
		}

		if stringChar == " " {
			if token != "" {
				trimToken := strings.TrimSpace(token)
				tokenList = append(tokenList, strings.TrimSpace(trimToken))
				token = ""
			}
			continue
		}
		token += stringChar
	}

	return tokenList
}

func TokenType(token string) string {
	keywordList := []string{"class", "constructor", "function", "method", "field", "static", "var", "int", "char", "boolean", "void", "true", "false", "null", "this", "let", "do", "if", "else", "while", "return"}

	for _, v := range keywordList {
		if token == v {
			return TokenTypeKeyword
		}
	}

	if regexp.MustCompile(`^[\-]*[0-9]+$`).MatchString(token) {
		return TokenTypeIntConst
	}

	if regexp.MustCompile(`^"`).MatchString(token) {
		return TokenTypeStringConst
	}

	if regexp.MustCompile(`[\{\}\(\)\[\]\.,;\+\-\*\/&\|<>=~]`).MatchString(token) {
		return TokenTypeSymbol
	}

	return TokenTypeIdentifer
}

func tokenToXml(token string, outputFile *os.File) {
	keywordList := []string{"class", "constructor", "function", "method", "field", "static", "var", "int", "char", "boolean", "void", "true", "false", "null", "this", "let", "do", "if", "else", "while", "return"}

	for _, v := range keywordList {
		if token == v {
			fmt.Fprint(outputFile, "<keyword> ")
			fmt.Fprint(outputFile, token)
			fmt.Fprintln(outputFile, " </keyword>")
			return

		}
	}

	if regexp.MustCompile(`[\{\}\(\)\[\]\.,;\+\-\*\/&\|<>=~]`).MatchString(token) {
		fmt.Fprint(outputFile, "<symbol> ")
		if token == "<" {
			fmt.Fprint(outputFile, "&lt;")
		} else if token == ">" {
			fmt.Fprint(outputFile, "&gt;")
		} else if token == "&" {
			fmt.Fprint(outputFile, "&amp;")
		} else {
			fmt.Fprint(outputFile, token)
		}

		fmt.Fprintln(outputFile, " </symbol>")
		return
	}

	if regexp.MustCompile(`^[0-9]+$`).MatchString(token) {
		fmt.Fprint(outputFile, "<integerConstant> ")
		fmt.Fprint(outputFile, token)
		fmt.Fprintln(outputFile, " </integerConstant>")
		return
	}

	if regexp.MustCompile(`"`).MatchString(token) {
		fmt.Fprint(outputFile, "<stringConstant> ")
		fmt.Fprint(outputFile, token[1:len(token)-1])
		fmt.Fprintln(outputFile, " </stringConstant>")
		return
	}

	fmt.Fprint(outputFile, "<identifier> ")
	fmt.Fprint(outputFile, token)
	fmt.Fprintln(outputFile, " </identifier>")
	return
}
