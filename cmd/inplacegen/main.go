package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	leftD, rightD = "{{", "}}" // todo[maybe]: over write by option
)

func main() {
	var name, file string
	var printOnly bool
	flag.StringVar(&name, "name", "", "name")
	flag.StringVar(&file, "file", os.Getenv("GOFILE"), "go file")
	flag.BoolVar(&printOnly, "print-only", false, "print only instead of overwrite file")
	flag.Parse()

	log.Printf("flags: name=%s, file=%s", name, file)

	if file == "" {
		log.Fatalf("no go file")
	}

	lines, err := ReadFileAsLines(file)
	if err != nil {
		log.Fatalf("read go file %q error: %s", file, err)
	}
	var result []string
	result, err = explain(lines, name)
	if err != nil {
		log.Fatalf("explain error: %s", err)
	}
	resultText := strings.Join(result, "\n")
	if printOnly {
		fmt.Println(resultText)
		return
	}
	err = ioutil.WriteFile(file, []byte(resultText), 0666)
	if err != nil {
		log.Fatalf("write file error: %s", err)
	}
	log.Print("success!")
}
