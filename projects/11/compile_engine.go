package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"strings"
	// "github.com/google/uuid"
)

func writeInit(outputFile *os.File) {
	//fmt.Fprintln(outputFile, "@SP")
	//fmt.Fprintln(outputFile, "M=256")
	fmt.Fprintln(outputFile, "goto Sys.init")
}

func CompileClass(tokenList []string, outputFile *os.File) {
	// fmt.Fprintln(outputFile, "<class>")

	var tokeni int

	class = tokenList[1]

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

	// fmt.Fprint(outputFile, "</class>")
}

func compileClassVarDec(tokenList []string, outputFile *os.File) int {
	// fmt.Fprintln(outputFile, "<classVarDec>")

	var tokeni int

	kind := ""
	if tokenList[0] == "static" {
		kind = SymbolKindStatic
	} else if tokenList[0] == "field" {
		kind = SymbolKindField
	}
	varType := tokenList[1]

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
			if i != 1 {
				classSymbolTable.define(token, varType, kind)
			}
		}

		if _type == TokenTypeSymbol {
			compileSymbol(token, outputFile)

			if token == ";" {
				// fmt.Fprintln(outputFile, "</classVarDec>")
				return i
			}
		}
	}

	return 0
}

func compileSubroutine(tokenList []string, outputFile *os.File) int {
	// fmt.Fprintln(outputFile, "<subroutineDec>")

	var tokeni int

	subroutineSymbolTable = NewSubroutineSymbolTable()
	var subroutineName string
	varCount = 0

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
				// fmt.Fprintln(outputFile, "<subroutineBody>")

				var tokenii int
				subroutineTokenList := tokenList[i:]
				for ii, token := range subroutineTokenList {
					if tokenii != 0 && ii <= tokenii {
						continue
					}

					_type := TokenType(token)

					if _type == TokenTypeSymbol && token == "}" {
						// fmt.Fprint(outputFile, "<symbol> ")
						// fmt.Fprint(outputFile, token)
						// fmt.Fprintln(outputFile, " </symbol>")

						// fmt.Fprintln(outputFile, "</subroutineBody>")
						// fmt.Fprintln(outputFile, "</subroutineDec>")
						return i + ii
					}

					if _type == TokenTypeSymbol {
						// fmt.Fprint(outputFile, "<symbol> ")
						// fmt.Fprint(outputFile, token)
						// fmt.Fprintln(outputFile, " </symbol>")
					}

					if _type == TokenTypeKeyword {
						if token == "var" {
							tokenii = ii + compileVarDec(subroutineTokenList[ii:], outputFile)
							continue
						}

						if token == "let" || token == "if" || token == "while" || token == "do" || token == "return" {
							writeFunction(subroutineName, varCount, outputFile)

							if tokenList[0] == "constructor" {
								size := classSymbolTable.varCount(SymbolKindField)
								size++
								writePush(SegmentConstant, strconv.Itoa(size), outputFile)
								writeCall("Memory.alloc", 1, outputFile)
								writePop(SegmentPointer, "0", outputFile)
							}
							if tokenList[0] == "method" {
								writePush(SegmentArg, "0", outputFile)
								// this
								writePop(SegmentPointer, "0", outputFile)
							}

							tokenii = ii + compileStatements(subroutineTokenList[ii:], outputFile)
							continue
						}
					}

				}
			}

			compileSymbol(token, outputFile)

			if token == "(" {
				if tokenList[0] == "method" {
					subroutineSymbolTable.define("0_this", "0_this", SymbolKindArg)
				}
				tokeni = i + compileParameterList(tokenList[i+1:], outputFile)
				if tokenList[0] == "method" {
					parameterCount++
				}
				
				subroutineName = class + "." + tokenList[2]
				continue
			}
		}
	}

	return 0
}

func compileParameterList(tokenList []string, outputFile *os.File) int {
	// fmt.Fprintln(outputFile, "<parameterList>")

	var tokeni int

	var varType string
	parameterCount = 0

	for i, token := range tokenList {
		if tokeni != 0 && i <= tokeni {
			continue
		}

		_type := TokenType(token)

		if _type == TokenTypeKeyword || token == "Array" {
			compileKeyword(token, outputFile)
			varType = token
			continue
		}

		if _type == TokenTypeIdentifer {
			compileIdentifer(token, outputFile)
			subroutineSymbolTable.define(token, varType, SymbolKindArg)
			parameterCount++
		}

		if _type == TokenTypeSymbol {
			if token == ")" {
				// fmt.Fprintln(outputFile, "</parameterList>")
				return i
			}

			compileSymbol(token, outputFile)
		}
	}

	return 0
}

func compileVarDec(tokenList []string, outputFile *os.File) int {
	// fmt.Fprintln(outputFile, "<varDec>")

	var tokeni int

	varType := tokenList[1]

	for i, token := range tokenList {
		if tokeni != 0 && i <= tokeni {
			continue
		}

		_type := TokenType(token)

		if _type == TokenTypeKeyword {
			compileKeyword(token, outputFile)
		}

		if _type == TokenTypeIdentifer {
			// compileIdentifer(token, outputFile)
			if i != 1 {
				varCount++
				subroutineSymbolTable.define(token, varType, SymbolKindVar)
			}
		}

		if _type == TokenTypeSymbol {
			compileSymbol(token, outputFile)

			if token == ";" {
				// fmt.Fprintln(outputFile, "</varDec>")
				return i
			}
		}
	}

	return 0
}

func compileStatements(tokenList []string, outputFile *os.File) int {
	// fmt.Fprintln(outputFile, "<statements>")

	var tokeni int

	fmt.Println("statements")
	fmt.Println(tokenList)

	for i, token := range tokenList {
		if tokeni != 0 && i <= tokeni {
			continue
		}

		fmt.Println(token)

		_type := TokenType(token)

		if _type == TokenTypeKeyword {
			if token == "let" {
				tokeni = i + compileLet(tokenList[i:], outputFile, false)
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
				// fmt.Fprintln(outputFile, "</statements>")
				return i - 1
			}
		}
	}

	return 0
}

func compileDo(tokenList []string, outputFile *os.File) int {
	// fmt.Fprintln(outputFile, " <doStatement>")

	var tokeni int
	var subroutineName string

	for i, token := range tokenList {
		if tokeni != 0 && i <= tokeni {
			continue
		}
		fmt.Println(token)

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
				for _, t := range tokenList[1:i] {
					subroutineName += t
				}

				var objectFlag bool
				subroutineName, objectFlag = compileObject(subroutineName, outputFile)

				tokeni = i + 1 + compileExpressionList(tokenList[i+1:], outputFile)

				if objectFlag {
					expressionCount++
				}
				continue
			}

			if token == ";" {
				// fmt.Fprintln(outputFile, " </doStatement>")
				
				if !strings.Contains(subroutineName, ".") {
					subroutineName = class + "." + subroutineName
				}
				writeCall(subroutineName, expressionCount, outputFile)
				return i
			}
		}
	}

	return 0
}

func compileLet(tokenList []string, outputFile *os.File, equalFlag bool) int {
	// fmt.Fprintln(outputFile, "<letStatement>")

	var tokeni int

	var varName string
	var indexFlag bool

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
			nextToken := tokenList[i+1]
			varName += token

			if nextToken == "[" {
				if equalFlag {
					kind, index := symbolKindIndex(token)
					writePush(kind, index, outputFile)
				}
				indexFlag = true
			}
		}

		if _type == TokenTypeSymbol {
			compileSymbol(token, outputFile)

			if token == "[" {
				tokeni = i + 1 + compileExpression(tokenList[i+1:], outputFile)
				if equalFlag {
					writeArithmetic("+", outputFile)
				}
				continue
			}

			if token == "=" {
				if equalFlag {
					if indexFlag {
						writePop(SegmentPointer, "1", outputFile)
						writePop(SegmentThat, "0", outputFile)
						return 0
					} else {
						kind, index := symbolKindIndex(varName)
						writePop(kind, index, outputFile)
						return 0
					}
				} else {
					tokeni = i + 1 + compileExpression(tokenList[i+1:], outputFile)
				}
				continue
			}

			if token == ";" {
				// fmt.Fprintln(outputFile, " </letStatement>")

				_ = compileLet(tokenList, outputFile, true)

				return i
			}
		}
	}

	return 0
}

func compileWhile(tokenList []string, outputFile *os.File) int {
	// fmt.Fprintln(outputFile, "<whileStatement>")

	var tokeni int

	startLabel := generateUuid()
	var endLabel string

	writeLabel(startLabel, outputFile)

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

				endLabel = generateUuid()
				writeArithmetic("~", outputFile)
				writeIf(endLabel, outputFile)

				continue
			}

			if token == "{" {
				tokeni = i + 1 + compileStatements(tokenList[i+1:], outputFile)
				continue
			}

			if token == "}" {
				writeGoto(startLabel, outputFile)
				writeLabel(endLabel, outputFile)
				// fmt.Fprintln(outputFile, "</whileStatement>")
				return i
			}
		}
	}

	return 0
}

func compileReturn(tokenList []string, outputFile *os.File) int {
	// fmt.Fprintln(outputFile, "<returnStatement>")
	// fmt.Println(tokenList)

	var tokeni int

	for i, token := range tokenList {
		if tokeni != 0 && i <= tokeni {
			continue
		}

		_type := TokenType(token)

		if _type == TokenTypeKeyword && token == "return" {
			compileKeyword(token, outputFile)
		} else if _type == TokenTypeSymbol && token == ";" {
			// fmt.Fprint(outputFile, "<symbol> ")
			// fmt.Fprint(outputFile, token)
			// fmt.Fprintln(outputFile, " </symbol>")

			// fmt.Fprintln(outputFile, "</returnStatement>")
			writeReturn(outputFile)
			return i
		} else {
			tokeni = i + compileExpression(tokenList[i:], outputFile)
			continue
		}
	}

	writeReturn(outputFile)

	return 0
}

func compileIf(tokenList []string, outputFile *os.File) int {
	// fmt.Fprintln(outputFile, "<ifStatement>")

	var tokeni int
	var elseFlag bool

	endLabel := generateUuid()
	elseLabel := generateUuid()

	for i, token := range tokenList {
		if tokeni != 0 && i <= tokeni {
			continue
		}

		_type := TokenType(token)
		fmt.Println(token)
		

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
				if !elseFlag {
					writeArithmetic("~", outputFile)
					writeIf(endLabel, outputFile)
				}

				tokeni = i + 1 + compileStatements(tokenList[i+1:], outputFile)
				continue
			}

			if token == "}" {
				if elseFlag {
					writeLabel(elseLabel, outputFile)
					return i
				}

				nextToken := tokenList[i+1]
				if nextToken == "else" {
					writeGoto(elseLabel, outputFile)
					writeLabel(endLabel, outputFile)
					elseFlag = true
					continue
				} else {
					writeLabel(endLabel, outputFile)
				}
				// fmt.Fprintln(outputFile, "</ifStatement>")

				return i
			}
		}
	}

	return 0
}

func compileExpression(tokenList []string, outputFile *os.File) int {
	// fmt.Fprintln(outputFile, "<expression>")

	var tokeni int

	var symbol string

	for i, token := range tokenList {
		if tokeni != 0 && i <= tokeni {
			continue
		}

		_type := TokenType(token)
		fmt.Println(token)
		fmt.Println(_type)

		if _type == TokenTypeIntConst {
			tokeni = i + compileTerm(tokenList[i:], outputFile)
			writePush(SegmentConstant, token, outputFile)
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
			if token == ")" || token == ";" || token == "," {
				// fmt.Fprintln(outputFile, "</expression>")

				writeArithmetic(symbol, outputFile)
				return i - 1
			}

			if token == "]" {
				// fmt.Fprintln(outputFile, "</expression>")

				writeArithmetic(symbol, outputFile)
				return i
			}

			if token == "~" || token == "(" || token == "[" {
				tokeni = i + compileTerm(tokenList[i:], outputFile)
				continue
			}

			if token == "-" {
				if i == 0 {
					tokeni = i + compileTerm(tokenList[i:], outputFile)
					continue
				} else if TokenType(tokenList[i-1]) == TokenTypeSymbol && Op(tokenList[i-1]) {
					tokeni = i + compileTerm(tokenList[i:], outputFile)
					writeArithmetic(tokenList[i-1], outputFile)
					continue
				} else {
					compileSymbol(token, outputFile)
					symbol = token
					continue
				}
			}

			compileSymbol(token, outputFile)
			symbol = token
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
	// fmt.Fprintln(outputFile, "<term>")

	var tokeni int

	var subroutineName string

	for i, token := range tokenList {
		if tokeni != 0 && i <= tokeni {
			continue
		}

		_type := TokenType(token)

		if _type == TokenTypeIntConst {
			// fmt.Fprint(outputFile, "<integerConstant> ")
			// fmt.Fprint(outputFile, token)
			// fmt.Fprintln(outputFile, " </integerConstant>")
			// fmt.Fprintln(outputFile, "</term>")
			return i
		}

		if _type == TokenTypeStringConst {
			// fmt.Fprint(outputFile, "<stringConstant> ")
			// fmt.Fprint(outputFile, token)
			// fmt.Fprintln(outputFile, " </stringConstant>")
			// fmt.Fprintln(outputFile, "</term>")

			writePush(SegmentConstant, strconv.Itoa(len(token)-2), outputFile)
			writeCall("String.new", 1, outputFile)
			for _, c := range token {
				stringChar := string(c)

				if stringChar == "\"" {
					continue
				}

				writePush(SegmentConstant, strconv.Itoa(int(c)), outputFile)
				writeCall("String.appendChar", 2, outputFile)
			}
			return i
		}

		if _type == TokenTypeKeyword {
			// fmt.Fprintln(outputFile, "</term>")
	
			compileKeyword(token, outputFile)
			if token == "false" || token == "null" {
				writePush(SegmentConstant, "0", outputFile)
			}
			if token == "true" {
				writePush(SegmentConstant, "-1", outputFile)
			}
			if token == "this" {
				writePush(SegmentPointer, "0", outputFile)
			}
			return i
		}

		if _type == TokenTypeIdentifer {
			nextToken := tokenList[i+1]
			compileIdentifer(token, outputFile)
			subroutineName += token
			
			if nextToken == "(" || nextToken == "." {
				continue
			}

			kind, index := symbolKindIndex(token)
			writePush(kind, index, outputFile)

			if nextToken == "[" {
				continue
			}

			// fmt.Fprintln(outputFile, "</term>")
			return i
		}

		if _type == TokenTypeSymbol {
			compileSymbol(token, outputFile)
			if token == "." {
				subroutineName += token
			}

			if token == "[" {
				tokeni = i + 1 + compileExpression(tokenList[i+1:], outputFile)
				writeArithmetic("+", outputFile)
				writePop(SegmentPointer, "1", outputFile)
				writePush(SegmentThat, "0", outputFile)
				return tokeni
			}

			if token == "(" && i == 0 {
				tokeni = i + 1 + compileExpression(tokenList[i+1:], outputFile)
				continue
			}
			if token == "(" && i > 0 {
				var objectFlag bool
				subroutineName, objectFlag = compileObject(subroutineName, outputFile)

				tokeni = i + 1 + compileExpressionList(tokenList[i+1:], outputFile)

				if objectFlag {
					expressionCount++
				}
				if !strings.Contains(subroutineName, ".") {
					subroutineName = class + "." + subroutineName
				}
				writeCall(subroutineName, expressionCount, outputFile)
				continue
			}

			if token == "]" {
				return i - 1
			}

			if token == ")" {
				// fmt.Fprintln(outputFile, " </term>")
				return i
			}

			if token == "-" || token == "~" {
				// fmt.Fprintln(outputFile, "</term>")

				result := compileTerm(tokenList[i+1:], outputFile)
				writeArithmetic(token, outputFile)
				return i + 1 + result
			}
		}
	}

	return 0
}

func compileExpressionList(tokenList []string, outputFile *os.File) int {
	// fmt.Fprintln(outputFile, " <expressionList>")
	expressionCount = 0

	if tokenList[0] == ")" {
		// fmt.Fprintln(outputFile, " </expressionList>")
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
				// fmt.Fprint(outputFile, "<symbol> ")
				// fmt.Fprint(outputFile, token)
				// fmt.Fprintln(outputFile, " </symbol>")
				continue
			}

			if token == ")" {
				// fmt.Fprintln(outputFile, " </expressionList>")
				return i - 1
			}
		}

		tokeni = i + compileExpression(tokenList[i:], outputFile)
		expressionCount++
		continue
	}

	return 0
}

func generateUuid() string {
	t := time.Now()

	return strconv.Itoa(t.Hour()) + strconv.Itoa(t.Minute()) + strconv.Itoa(t.Second()) + "_" + strconv.Itoa(t.Nanosecond())
}

func symbolKindIndex(varName string) (kind, index string) {
	kind = subroutineSymbolTable.kindOf(varName)
	index = subroutineSymbolTable.indexOf(varName)
	if kind == SymbolKindNone {
		kind = classSymbolTable.kindOf(varName)
		index = classSymbolTable.indexOf(varName)
	}

	if kind == SymbolKindStatic {
		kind = SegmentStatic
	}
	if kind == SymbolKindField {
		kind = SegmentThis
	}
	if kind == SymbolKindArg {
		kind = SegmentArg
	}
	if kind == SymbolKindVar {
		kind = SegmentLocal
	}

	return kind, index
}

func symbolType(objectName string) string {
	_type := subroutineSymbolTable.typeOf(objectName)
	if _type == "" {
		_type = classSymbolTable.typeOf(objectName)
	}

	return _type
}

func compileObject(subroutineName string, outputFile *os.File) (string, bool) {
	var objectFlag bool

	if strings.Contains(subroutineName, ".") {
		objectName := subroutineName[:strings.Index(subroutineName, ".")]
		kind, index := symbolKindIndex(objectName)
		
		if kind != SymbolKindNone {
			_type := symbolType(objectName)
			if _type != "int" && _type != "char" &&_type != "boolean" {
				subroutineName = _type + subroutineName[strings.Index(subroutineName, "."):]
				writePush(kind, index, outputFile)
				objectFlag = true
			}
		}
	}

	if !strings.Contains(subroutineName, ".") {
		writePush(SegmentPointer, "0", outputFile)
		objectFlag = true
	}

	return subroutineName, objectFlag
}

func compileKeyword(token string, outputFile *os.File) {
	return
	fmt.Fprint(outputFile, "<keyword> ")
	fmt.Fprint(outputFile, token)
	fmt.Fprintln(outputFile, " </keyword>")
}

func compileIdentifer(token string, outputFile *os.File) {
	return
	fmt.Fprint(outputFile, "<identifier> ")
	fmt.Fprint(outputFile, token)
	fmt.Fprintln(outputFile, " </identifier>")
}

func compileSymbol(token string, outputFile *os.File) {
	return
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
