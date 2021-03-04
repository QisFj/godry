package main

import (
	"regexp"
	"strconv"
	"strings"
)

func replaceT(t, ld, rd string, data []Entry) string {
	replaced := ""
	ldN, rdN := len(ld), len(rd)
	for {
		ldIdx := strings.Index(t, ld)
		if ldIdx < 0 {
			break
		}
		rdIdx := strings.Index(t, rd)
		if rdIdx < 0 {
			break
		}
		replaced += t[:ldIdx] + ld + replace(t[ldIdx+ldN:rdIdx], data) + rd
		t = t[rdIdx+rdN:]
	}
	return replaced + t
}

var (
	replaceRegexp1 = regexp.MustCompile(`\$\d+`)
	replaceRegexp2 = regexp.MustCompile(`\$\[\d+]`)
	replaceRegexp3 = regexp.MustCompile(`\$\[\d+\.\d+]`)
)

func replace(t string, data []Entry) string {
	t = replaceRegexp1.ReplaceAllStringFunc(t, func(s string) string {
		return "\"" + data[mustInt(s[1:])-1][0] + "\""
	})
	t = replaceRegexp2.ReplaceAllStringFunc(t, func(s string) string {
		return "\"" + data[mustInt(s[2:len(s)-1])-1][0] + "\""
	})
	t = replaceRegexp3.ReplaceAllStringFunc(t, func(s string) string {
		idx := strings.Index(s, ".")
		return "\"" + data[mustInt(s[2:idx])-1][mustInt(s[idx+1:len(s)-1])-1] + "\""
	})
	return t
}

func mustInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
