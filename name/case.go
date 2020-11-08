package name

import (
	"regexp"
	"strings"
)

// reference: https://gist.github.com/stoewer/fbe273b711e6a06315d19552dd4d33e6

// in above gist, matchFirstCap use regexp `(.)([A-Z][a-z]+)`, where `.` would match any character,
// witch make separators such as `.` `,` also be matched
// so use `[A-Za-z0-9]` instead, to avoid match special character
var matchFirstCap = regexp.MustCompile(`([A-Za-z0-9])([A-Z][a-z]+)`)
var matchAllCap = regexp.MustCompile(`([a-z0-9])([A-Z])`)

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

var link = regexp.MustCompile(`(^[A-Za-z])|_([A-Za-z])`)

func ToCamelCase(str string) string {
	return link.ReplaceAllStringFunc(str, func(s string) string {
		return strings.ToUpper(strings.Replace(s, "_", "", -1))
	})
}
