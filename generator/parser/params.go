package parser

import "strings"

func MapParams(paramStr string) map[string]string {
	paramMap := make(map[string]string)
	params := strings.Split(paramStr, ", ")
	pendingParams := make([]string, 0)
	for _, param := range params {
		paramFields := strings.Fields(param)
		if len(param) == 1 {
			pendingParams = append(pendingParams, paramFields[0])
			continue
		}
		if len(pendingParams) > 0 {
			for _, pendingParam := range pendingParams {
				paramMap[pendingParam] = paramFields[1]
			}
			pendingParams = make([]string, 0)
		}
		paramMap[paramFields[0]] = paramFields[1]
	}
	return paramMap
}
