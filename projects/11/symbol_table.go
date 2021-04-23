package main

import "strconv"

const SymbolKindStatic = "STATIC"
const SymbolKindField = "FIELD"
const SymbolKindArg = "ARGUMENT"
const SymbolKindVar = "VAR"
const SymbolKindNone = "NONE"

var classSymbolTable *ClassSymbolTable
var subroutineSymbolTable *SubroutineSymbolTable

type ClassSymbolTable struct {
	_static map[string][]string
	_field  map[string][]string
}

type SubroutineSymbolTable struct {
	_argument map[string][]string
	_var      map[string][]string
}

func NewClassSymbolTable() *ClassSymbolTable {
	return &ClassSymbolTable{make(map[string][]string), make(map[string][]string)}
}

func NewSubroutineSymbolTable() *SubroutineSymbolTable {
	return &SubroutineSymbolTable{map[string][]string{}, map[string][]string{}}
}

func (t *ClassSymbolTable) define(name string, _type string, kind string) {
	if kind == SymbolKindStatic {
		t._static[name] = []string{_type, strconv.Itoa(len(t._static))}
	} else if kind == SymbolKindField {
		t._field[name] = []string{_type, strconv.Itoa(len(t._field))}
	}
}

func (t *ClassSymbolTable) varCount(kind string) int {
	if kind == SymbolKindStatic {
		return len(t._static)
	} else if kind == SymbolKindField {
		return len(t._field)
	}

	return 0
}

func (t *ClassSymbolTable) kindOf(name string) string {
	if len(t._static[name]) != 0 {
		return SymbolKindStatic
	}

	if len(t._field[name]) != 0 {
		return SymbolKindField
	}

	return SymbolKindNone
}

func (t *ClassSymbolTable) typeOf(name string) string {
	if len(t._static[name]) != 0 {
		return t._static[name][0]
	}

	if len(t._field[name]) != 0 {
		return t._field[name][0]
	}

	return ""
}

func (t *ClassSymbolTable) indexOf(name string) string {
	if len(t._static[name]) != 0 {
		return t._static[name][1]
	}

	if len(t._field[name]) != 0 {
		return t._field[name][1]
	}

	return "0"
}

func (st *SubroutineSymbolTable) define(name string, _type string, kind string) {
	if kind == SymbolKindArg {
		st._argument[name] = []string{_type, strconv.Itoa(len(st._argument))}
	} else if kind == SymbolKindVar {
		st._var[name] = []string{_type, strconv.Itoa(len(st._var))}
	}
}

func (st *SubroutineSymbolTable) varCount(kind string) int {
	if kind == SymbolKindArg {
		return len(st._argument)
	} else if kind == SymbolKindVar {
		return len(st._var)
	}

	return 0
}

func (st *SubroutineSymbolTable) kindOf(name string) string {
	if len(st._argument[name]) != 0 {
		return SymbolKindArg
	}

	if len(st._var[name]) != 0 {
		return SymbolKindVar
	}

	return SymbolKindNone
}

func (st *SubroutineSymbolTable) typeOf(name string) string {
	if len(st._argument[name]) != 0 {
		return st._argument[name][0]
	}

	if len(st._var[name]) != 0 {
		return st._var[name][0]
	}

	return ""
}

func (st *SubroutineSymbolTable) indexOf(name string) string {
	if len(st._argument[name]) != 0 {
		return st._argument[name][1]
	}

	if len(st._var[name]) != 0 {
		return st._var[name][1]
	}

	return "0"
}
