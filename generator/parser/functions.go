package parser

import (
	"strings"
)

type Function struct {
	Name   string
	Params string
	Return string
}

func ParseFunctions(s string) []*Function {
	fs := make([]*Function, 0)
	for _, l := range strings.Split(s, "\n\n") {
		if f := parseFunction(l); f.Name != "" {
			fs = append(fs, f)
		}
	}
	return fs
}

func parseFunction(s string) *Function {
	funcInfo := new(Function)
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		funcString := "func "
		// func Hemlo(...) (...) {...}
		if !strings.HasPrefix(line, funcString) {
			continue
		}
		// Hemlo(...) (...) {...}
		line = strings.TrimPrefix(line, funcString)
		// [Hemlo, (...) (...) {...}]
		funcName := strings.Split(line, "(")[0]
		funcInfo.Name = strings.TrimSpace(funcName)
		// ...) (...) {...}
		line = strings.TrimPrefix(line, funcName)[1:]
		// ...)
		funcParams := strings.SplitN(line, ")", 2)[0]
		// ...
		funcInfo.Params = strings.TrimSpace(funcParams)
		//  (...) {...}
		line = strings.TrimPrefix(line, funcParams)[1:]
		//  (...)
		dirtyfuncReturn := strings.SplitN(line, "{", 2)[0]
		// (...)
		dirtyfuncReturn = strings.TrimSpace(dirtyfuncReturn)
		funcInfo.Return = strings.TrimPrefix(strings.TrimSuffix(dirtyfuncReturn, ")"), "(")
		break
	}
	return funcInfo
}
