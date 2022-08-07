package parser

import (
	"strings"
)

type Method struct {
	Owner   string
	Pointer bool
	Name    string
	Params  string
	Return  string
}

func ParseMethods(s string) []*Method {
	fs := make([]*Method, 0)
	for _, l := range strings.Split(s, "\n\n") {
		if f := parseMethod(l); f.Name != "" {
			fs = append(fs, f)
		}
	}
	return fs
}

func parseMethod(s string) *Method {
	funcInfo := new(Method)
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		funcString := "func "
		// func (c *Something) Hemlo(...) (...) {...}
		if !strings.HasPrefix(line, funcString) {
			continue
		}
		// (c *Something) Hemlo(...) (...) {...}
		line = strings.TrimPrefix(line, funcString)
		if !strings.HasPrefix(strings.TrimSpace(line), "(") {
			continue
		}
		dirtyMethodOwner := strings.SplitN(line, ")", 2)[0]
		if arr := strings.Fields(strings.TrimSpace(dirtyMethodOwner)); arr != nil {
			var dirtyMethod string
			if len(arr) == 1 {
				dirtyMethod = arr[0]
			} else {
				dirtyMethod = arr[1]
			}
			dirtyMethod = strings.TrimPrefix(dirtyMethod, "(")
			if strings.HasPrefix(dirtyMethod, "*") {
				funcInfo.Pointer = true
				dirtyMethod = dirtyMethod[1:]
			}
			funcInfo.Owner = dirtyMethod
		}
		line = strings.SplitN(line, ")", 2)[1]
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
