package utils

import "regexp"

const (
	timeRangeRegex = `^-?([0]{1}\.{1}[0-9]+|[1-9]{1}[0-9]*\.{1}[0-9]+|[0-9]+|0)[smhdw]$`
)

func IsValidTimeRange(s string) bool {
	re := regexp.MustCompile(timeRangeRegex)

	return re.MatchString(s)

}
