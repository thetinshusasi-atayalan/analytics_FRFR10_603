package utils

import "strconv"

func Str2Int(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}

func Str2Int64(s string) int64 {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return val
}

func Int64ToStr(i int64) string {
	val := strconv.FormatInt(i, 10)

	return val
}

func IndexOfForStringArray(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}
