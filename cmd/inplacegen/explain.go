package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"text/template"

	"github.com/QisFj/godry/cexp"
	"github.com/QisFj/godry/gen/graph"
	"github.com/QisFj/godry/name"
	"github.com/QisFj/godry/slice"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

var (
	inPlaceGenFromRegexp = regexp.MustCompile(`INPLACEGEN_\(.+\)_FROM`)
	inPlaceGenToRegexp   = regexp.MustCompile(`INPLACEGEN_\(.+\)_TO`)
)

func isInPlaceGenFromLine(line string) (prefix string, name string, is bool) {
	loc := inPlaceGenFromRegexp.FindStringIndex(line)
	if loc == nil {
		return "", "", false
	}
	return line[:loc[0]], line[loc[0]+12 : loc[1]-6], true
}

func isInPlaceGenToLine(line string) (name string, is bool) {
	loc := inPlaceGenToRegexp.FindStringIndex(line)
	if loc == nil {
		return "", false
	}
	return line[loc[0]+12 : loc[1]-4], true
}

func explain(lines []string, expName string) (resultLines []string, err error) {
	var prefix, fromName, toName string
	var ok bool
	var i, last int
	for {
		ok = false
		for ; i < len(lines); i++ {
			prefix, fromName, ok = isInPlaceGenFromLine(lines[i])
			if ok && (expName == "" || expName == fromName) {
				break
			}
		}
		if !ok {
			break
		}
		resultLines = append(resultLines, lines[last:i]...)
		last = i

		i += 1
		ok = false
		for ; i < len(lines); i++ {
			toName, ok = isInPlaceGenToLine(lines[i])
			if ok {
				if toName != fromName {
					return nil, fmt.Errorf("inplacegen (from=%s,to=%s) not pair", fromName, toName)
				}
				i += 1
				break
			}
		}
		if !ok {
			return nil, fmt.Errorf("inplacegen (from=%s,to) not pair", fromName)
		}
		log.Printf("try to explain %s line:%d to line:%d", fromName, last, i)
		var gened []string
		gened, err = explainToGen(lines[last:i], prefix)
		if err != nil {
			return nil, err
		}
		resultLines = append(resultLines, gened...)
		last = i
	}
	resultLines = append(resultLines, lines[last:i]...)
	return resultLines, nil
}

func explainToGen(lines []string, prefix string) (result []string, err error) {
	// 1st line is FROM
	result = append(result, lines[0])
	i := 1
	for ; i < len(lines)-1; i++ {
		if !strings.HasPrefix(lines[i], prefix) {
			break
		}
		result = append(result, lines[i])
	}
	var arg Arg
	if arg, err = explainToArg(slice.Map(lines[1:i], func(_ int, v string) string {
		return strings.TrimPrefix(v, prefix)
	})); err != nil {
		return nil, err
	}
	log.Printf("arg.Data: %q", arg.Data)
	log.Printf("arg.ExData: %q", arg.ExData)
	log.Printf("arg.Template: %s", arg.Template)
	result = append(result, arg.Gen()...)
	log.Printf("gen success")

	// last line is TO
	result = append(result, lines[len(lines)-1])
	return result, nil
}

func explainToArg(lines []string) (arg Arg, err error) {
	trimmedFirstLine := strings.Trim(lines[0], " \t")
	if strings.HasPrefix(trimmedFirstLine, "import ") {
		filename := strings.Trim(trimmedFirstLine[7:], " \t")
		lines, err = ReadFileAsLines(filename)
		if err != nil {
			return Arg{}, fmt.Errorf("import %q error: %w", filename, err)
		}
		trimmedFirstLine = strings.Trim(lines[0], " \t")
		log.Printf("imported %q", filename)
	}
	n := mustInt(trimmedFirstLine)
	for i := 1; i <= n; i++ {
		var ex bool
		var g Group
		dataContent := []byte(lines[i])
		if strings.HasPrefix(lines[i], "*") {
			ex = true
			dataContent = dataContent[1:]
		}
		if bytes.HasPrefix(dataContent, []byte("import ")) {
			filename := strings.Trim(string(dataContent[7:]), " \t")
			dataContent, err = ioutil.ReadFile(filename)
			if err != nil {
				return Arg{}, fmt.Errorf("import %q error: %w", filename, err)
			}
			log.Printf("imported %q as %sData", filename, cexp.String(ex, "Ex", ""))
		}
		if err = json5.Unmarshal(dataContent, &g); err != nil {
			return Arg{}, fmt.Errorf("json unmarshal error: %w", err)
		}
		if !ex {
			arg.Data = append(arg.Data, g)
		} else {
			arg.ExData = append(arg.ExData, g)
		}
	}
	arg.Template = strings.Join(lines[n+1:], "\n")
	return arg, nil
}

func (arg Arg) Gen() (result []string) {
	it := graph.NewIter(arg.Data, false)
	for it.Next() {
		entries := slice.Map(it.Get(), func(_ int, v graph.NodeI) Entry {
			return v.(Entry)
		})
		replaced := replaceT(arg.Template, leftD, rightD, entries)
		sb := &strings.Builder{}
		err := template.Must(template.New("").Delims(leftD, rightD).Funcs(template.FuncMap{
			"ToSnake": name.ToSnakeCase,
			"ToCamle": name.ToCamelCase,
		}).Parse(replaced)).Execute(sb, TemplateData(entries, arg.ExData))
		if err != nil {
			log.Fatalf("execute template error: %s", err)
		}
		result = append(result, strings.Split(sb.String(), "\n")...)
	}
	return result
}

func TemplateData(entries []Entry, exData Data) interface{} {
	mm := map[string]interface{}{}
	for i, e := range entries {
		mm[fmt.Sprintf("e%d", i+1)] = entryToMap(e)
	}
	for i, g := range exData {
		mm[fmt.Sprintf("ex%d", i+1)] = slice.Map(g, func(_ int, value Entry) map[string]string {
			return entryToMap(value)
		})
	}
	return mm
}

func entryToMap(e Entry) map[string]string {
	m := map[string]string{}
	for j, v := range e {
		m[fmt.Sprintf("v%d", j+1)] = v
	}
	return m
}
