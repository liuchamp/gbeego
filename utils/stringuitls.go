package utils

import "fmt"

// 转换interface为string

func checkString(value interface{}) (string, error) {
	ps, ok := value.(string)
	if ok {
		return ps, nil
	}
	return "", fmt.Errorf("value not is string")
}
