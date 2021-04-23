package main

import (
	"fmt"
	"os"
)

func writeInit(outputFile *os.File) {
	//fmt.Fprintln(outputFile, "@SP")
	//fmt.Fprintln(outputFile, "M=256")
	fmt.Fprintln(outputFile, "goto Sys.init")
}

func CompileClass(tokenList []string, outputFile *os.File) {
	fmt.Fprintln(outputFile, "<class>")

	var tokeni int

	for i, token := range tokenList {
		if tokeni != 0 && i <= tokeni {
			continue
		}

		_type := TokenType(token)

		if _type == TokenTypeKeyword {
			if token == "static" || token == "field" {
				tokeni = i + compileClassVarDec(tokenList[i:], outputFile)
				continue
			}
			if token == "constructor" || token == "function" || token == "method" {
				tokeni = i + compileSubroutine(tokenList[i:], outputFile)
				continue
			}

			compileKeyword(token, outputFile)
		}

		if _type == TokenTypeIdentifer {
			compileIdentifer(token, outputFile)
		}

		if _type == TokenTypeSymbol {
			compileSymbol(token, outputFile)
		}

	}

	fmt.Fprint(outputFile, "</class>")
}

func compileClassVarDec(tokenList []string, outputFile *os.File) int {
	fmt.Fprintln(outputFile, "<classVarDec>")

	var tokeni int

	for i, token := range tokenList {
		if tokeni != 0 && i <= tokeni {
			continue
		}

		_type := TokenType(token)

		if _type == TokenTypeKeyword {
			compileKeyword(token, outputFile)
		}

		if _type == TokenTypeIdentifer {
			compileIdentifer(token, outputFile)
		}

		if _type == TokenTypeSymbol {
			compileSymbol(token, outputFile)

			if token == ";" {
				fmt.Fprintln(outputFile, "</classVarDec>")
				return i
			}
		}
	}

	return 0
}

func compileSubroutine(tokenList []string, outputFile *os.File) int {
	fmt.Fprintln(outputFile, "<subroutineDec>")

	var tokeni int

	for i, token := range tokenList {
		if tokeni != 0 && i <= tokeni {
			continue
		}

		_type := TokenType(token)

		if _type == TokenTypeKeyword {
			compileKeyword(token, outputFile)
		}

		if _type == TokenTypeIdentifer {
			compileIdentifer(token, outputFile)
		}

		if _type == TokenTypeSymbol {
			if token == "{" {
				fmt.Fprintln(outputFile, "<subroutineBody>")

				var tokenii int
				subroutineTokenList := tokenList[i:]
				for ii, token := range subroutineTokenList {
					if tokenii != 0 && ii <= tokenii {
						continue
					}

					_type := TokenType(token)

					if _type == TokenTypeSymbol && token == "}" {
						fmt.Fprint(outputFile, "<symbol> ")
						fmt.Fprint(outputFile, token)
						fmt.Fprintln(outputFile, " </symbol>")

						fmt.Fprintln(outputFile, "</subroutineBody>")
						fmt.Fprintln(outputFile, "</subroutineDec>")
						return i + ii
					}

					if _type == TokenTypeSymbol {
						fmt.Fprint(outputFile, "<symbol> ")
						fmt.Fprint(outputFile, token)
						fmt.Fprintln(outputFile, " </symbol>")
					}

					if _type == TokenTypeKeyword {
						if token == "var" {
							tokenii = ii + compileVarDec(subroutineTokenList[ii:], outputFile)
							continue
						}

						if token == "let" || token == "if" || token == "while" || token == "do" || token == "return" {
							tokenii = ii + compileStatements(subroutineTokenList[ii:], outputFile)
							continue
						}
					}

				}
			}

			compileSymbol(token, outputFile)

			if token == "(" {
				tokeni = i + compileParameterList(tokenList[i+1:], outputFile)
				continue
			}
		}
	}

	return 0
}

func compileParameterList(tokenList []string, outputFile *os.File) int {
	fmt.Fprintln(outputFile, "<parameterList>")

	var tokeni int

	for i, token := range tokenList {
		if tokeni != 0 && i <= tokeni {
			continue
		}

		_type := TokenType(token)

		if _type == TokenTypeKeyword {
			compileKeyword(token, outputFile)
		}

		if _type == TokenTypeIdentifer {
			compileIdentifer(token, outputFile)
		}

		if _type == TokenTypeSymbol {
			if token == ")" {
				fmt.Fprintln(outputFile, "</parameterList>")
				return i
			}

			compileSymbol(token, outputFile)
		}
	}

	return 0
}

func compileVarDec(tokenList []string, outputFile *os.File) int {
	fmt.Fprintln(outputFile, "<varDec>")

	var tokeni int

	for i, token := range tokenList {
		if tokeni != 0 && i <= tokeni {
			continue
		}

		_type := TokenType(token)

		if _type == TokenTypeKeyword {
			compileKeyword(token, outputFile)
		}

		if _type == TokenTypeIdentifer {
			compileIdentifer(token, outputFile)
		}

		if _type == TokenTypeSymbol {
			compileSymbol(token, outputFile)

			if token == ";" {
				fmt.Fprintln(outputFile, "</varDec>")
				return i
			}
		}
	}

	return 0
}

func compileStatements(tokenList []string, outputFile *os.File) int {
	fmt.Fprintln(outputFile, "<statements>")

	var tokeni int

	for i, token := range tokenList {
		if tokeni != 0 && i <= tokeni {
			continue
		}

		_type := TokenType(token)

		if _type == TokenTypeKeyword {
			if token == "let" {
				tokeni = i + compileLet(tokenList[i:], outputFile)
				continue
			}

			if token == "if" {
				tokeni = i + compileIf(tokenList[i:], outputFile)
				continue
			}

			if token == "while" {
				tokeni = i + compileWhile(tokenList[i:], outputFile)
				continue
			}

			if token == "do" {
				tokeni = i + compileDo(tokenList[i:], outputFile)
				continue
			}

			if token == "return" {
				tokeni = i + compileReturn(tokenList[i:], outputFile)
				continue
			}
		}

		if _type == TokenTypeSymbol {
			if token == "}" {
				fmt.Fprintln(outputFile, "</statements>")
				return i - 1
			}
		}
	}

	return 0
}

func compileDo(tokenList []string, outputFile *os.File) int {
	fmt.Fprintln(outputFile, " <doStatement>")

	var tokeni int

	for i, token := range tokenList {
		if tokeni != 0 && i <= tokeni {
			continue
		}

		_type := TokenType(token)

		if _type == TokenTypeKeyword {
			compileKeyword(token, outputFile)
		}

		if _type == TokenTypeIdentifer {
			compileIdentifer(token, outputFile)
		}

		if _type == TokenTypeSymbol {
			compileSymbol(token, outputFile)

			if token == "(" {
				tokeni = i + 1 + compileExpressionList(tokenList[i+1:], outputFile)
				continue
			}

			if token == ";" {
				fmt.Fprintln(outputFile, " </doStatement>")
				return i
			}
		}
	}

	return 0
}

func compileLet(tokenList []string, outputFile *os.File) int {
	fmt.Fprintln(outputFile, "<letStatement>")

	var tokeni int

	for i, token := range tokenList {
		if tokeni != 0 && i <= tokeni {
			continue
		}

		_type := TokenType(token)

		if _type == TokenTypeKeyword {
			compileKeyword(token, outputFile)
		}

		if _type == TokenTypeIdentifer {
			compileIdentifer(token, outputFile)
		}

		if _type == TokenTypeSymbol {
			compileSymbol(token, outputFile)

			if token == "[" {
				tokeni = i + 1 + compileExpression(tokenList[i+1:], outputFile)
				continue
			}

			if token == "=" {
				tokeni = i + 1 + compileExpression(tokenList[i+1:], outputFile)
				continue
			}

			if token == ";" {
				fmt.Fprintln(outputFile, " </letStatement>")
				return i
			}
		}
	}

	return 0
}

func compileWhile(tokenList []string, outputFile *os.File) int {
	fmt.Fprintln(outputFile, "<whileStatement>")

	var tokeni int

	for i, token := range tokenList {
		if tokeni != 0 && i <= tokeni {
			continue
		}

		_type := TokenType(token)

		if _type == TokenTypeKeyword {
			compileKeyword(token, outputFile)
		}

		if _type == TokenTypeSymbol {
			compileSymbol(token, outputFile)

			if token == "(" {
				tokeni = i + 1 + compileExpression(tokenList[i+1:], outputFile)
				continue
			}

			if token == "{" {
				tokeni = i + 1 + compileStatements(tokenList[i+1:], outputFile)
				continue
			}

			if token == "}" {
				fmt.Fprintln(outputFile, "</whileStatement>")
				return i
			}
		}
	}

	return 0
}

func compileReturn(tokenList []string, outputFile *os.File) int {
	fmt.Fprintln(outputFile, "<returnStatement>")

	var tokeni int

	for i, token := range tokenList {
		if tokeni != 0 && i <= tokeni {
			continue
		}

		_type := TokenType(token)

		if _type == TokenTypeKeyword && token == "return" {
			compileKeyword(token, outputFile)
		} else if _type == TokenTypeSymbol && token == ";" {
			fmt.Fprint(outputFile, "<symbol> ")
			fmt.Fprint(outputFile, token)
			fmt.Fprintln(outputFile, " </symbol>")

			fmt.Fprintln(outputFile, "</returnStatement>")
			return i
		} else {
			tokeni = i + compileExpression(tokenList[i:], outputFile)
			continue
		}
	}

	return 0
}

func compileIf(tokenList []string, outputFile *os.File) int {
	fmt.Fprintln(outputFile, "<ifStatement>")

	var tokeni int

	for i, token := range tokenList {
		if tokeni != 0 && i <= tokeni {
			continue
		}

		_type := TokenType(token)

		if _type == TokenTypeKeyword {
			compileKeyword(token, outputFile)
		}

		if _type == TokenTypeSymbol {
			compileSymbol(token, outputFile)

			if token == "(" {
				tokeni = i + 1 + compileExpression(tokenList[i+1:], outputFile)
				continue
			}

			if token == "{" {
				tokeni = i + 1 + compileStatements(tokenList[i+1:], outputFile)
				continue
			}

			if token == "}" {
				nextToken := tokenList[i+1]
				if nextToken == "else" {
					continue
				}
				fmt.Fprintln(outputFile, "</ifStatement>")
				return i
			}
		}
	}

	return 0
}

func compileExpression(tokenList []string, outputFile *os.File) int {
	fmt.Fprintln(outputFile, "<expression>")

	var tokeni int

	for i, token := range tokenList {
		if tokeni != 0 && i <= tokeni {
			continue
		}

		_type := TokenType(token)

		if _type == TokenTypeIntConst {
			tokeni = i + compileTerm(tokenList[i:], outputFile)
			continue
		}

		if _type == TokenTypeStringConst {
			tokeni = i + compileTerm(tokenList[i:], outputFile)
			continue
		}

		if _type == TokenTypeKeyword {
			tokeni = i + compileTerm(tokenList[i:], outputFile)
			continue
		}

		if _type == TokenTypeIdentifer {
			tokeni = i + compileTerm(tokenList[i:], outputFile)
			continue
		}

		if _type == TokenTypeSymbol {
			if token == ")" || token == "]" || token == ";" || token == "," {
				fmt.Fprintln(outputFile, "</expression>")
				return i - 1
			}

			if token == "~" {
				tokeni = i + compileTerm(tokenList[i:], outputFile)
				continue
			}

			if token == "-" {
				if i == 0 {
					tokeni = i + compileTerm(tokenList[i:], outputFile)
					continue
				} else if TokenType(tokenList[i-1]) == TokenTypeSymbol && Op(tokenList[i-1]) {
					tokeni = i + compileTerm(tokenList[i:], outputFile)
					continue
				} else {
					compileSymbol(token, outputFile)
					continue
				}
			}

			if token == "(" {
				tokeni = i + compileTerm(tokenList[i:], outputFile)
				continue
			}

			compileSymbol(token, outputFile)
		}
	}

	return 0
}

func Op(token string) bool {
	if token == "+" || token == "-" || token == "*" || token == "/" || token == "&" || token == "|" || token == "<" || token == ">" || token == "=" {
		return true
	}

	return false
}

func compileTerm(tokenList []string, outputFile *os.File) int {
	fmt.Fprintln(outputFile, "<term>")

	var tokeni int

	for i, token := range tokenList {
		if tokeni != 0 && i <= tokeni {
			continue
		}

		_type := TokenType(token)

		if _type == TokenTypeIntConst {
			fmt.Fprint(outputFile, "<integerConstant> ")
			fmt.Fprint(outputFile, token)
			fmt.Fprintln(outputFile, " </integerConstant>")
			fmt.Fprintln(outputFile, "</term>")
			return i
		}

		if _type == TokenTypeStringConst {
			fmt.Fprint(outputFile, "<stringConstant> ")
			fmt.Fprint(outputFile, token)
			fmt.Fprintln(outputFile, " </stringConstant>")
			fmt.Fprintln(outputFile, "</term>")
			return i
		}

		if _type == TokenTypeKeyword {
			compileKeyword(token, outputFile)
			fmt.Fprintln(outputFile, "</term>")
			return i
		}

		if _type == TokenTypeIdentifer {
			nextToken := tokenList[i+1]
			compileIdentifer(token, outputFile)
			if nextToken == "[" || nextToken == "(" || nextToken == "." {
				continue
			}
			fmt.Fprintln(outputFile, "</term>")
			return i
		}

		if _type == TokenTypeSymbol {
			compileSymbol(token, outputFile)

			if token == "[" {
				tokeni = i + 1 + compileExpression(tokenList[i+1:], outputFile)
				continue
			}

			if token == "(" && i == 0 {
				tokeni = i + 1 + compileExpression(tokenList[i+1:], outputFile)
				continue
			}
			if token == "(" && i > 0 {
				tokeni = i + 1 + compileExpressionList(tokenList[i+1:], outputFile)
				continue
			}

			if token == "]" || token == ")" {
				fmt.Fprintln(outputFile, " </term>")
				return i
			}

			if token == "-" || token == "~" {
				result := compileTerm(tokenList[i+1:], outputFile)
				fmt.Fprintln(outputFile, "</term>")
				return i + 1 + result
			}
		}
	}

	return 0
}

func compileExpressionList(tokenList []string, outputFile *os.File) int {
	fmt.Fprintln(outputFile, " <expressionList>")

	if tokenList[0] == ")" {
		fmt.Fprintln(outputFile, " </expressionList>")
		return -1
	}

	var tokeni int

	for i, token := range tokenList {
		if tokeni != 0 && i <= tokeni {
			continue
		}

		_type := TokenType(token)

		if _type == TokenTypeSymbol {
			if token == "," {
				fmt.Fprint(outputFile, "<symbol> ")
				fmt.Fprint(outputFile, token)
				fmt.Fprintln(outputFile, " </symbol>")
				continue
			}

			if token == ")" {
				fmt.Fprintln(outputFile, " </expressionList>")
				return i - 1
			}
		}

		tokeni = i + compileExpression(tokenList[i:], outputFile)
		continue
	}

	return 0
}

func compileKeyword(token string, outputFile *os.File) {
	fmt.Fprint(outputFile, "<keyword> ")
	fmt.Fprint(outputFile, token)
	fmt.Fprintln(outputFile, " </keyword>")
}

func compileIdentifer(token string, outputFile *os.File) {
	fmt.Fprint(outputFile, "<identifier> ")
	fmt.Fprint(outputFile, token)
	fmt.Fprintln(outputFile, " </identifier>")
}

func compileSymbol(token string, outputFile *os.File) {
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
}
