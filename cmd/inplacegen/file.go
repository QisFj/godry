package main

import (
	"io/ioutil"
	"strings"
)

func ReadFileAsLines(filename string) ([]string, error) {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(fileContent), "\n"), nil
}
