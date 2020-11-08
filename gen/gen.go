package gen

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

func goSrcFileName(filename string) string {
	if !strings.HasSuffix(filename, ".go") {
		filename += ".go"
	}
	return filename
}

func Gen(filename, tmpl string, data interface{}, options ...Option) {
	cfg := newConfig(options...)
	filename = goSrcFileName(filename)
	var t *template.Template
	var err error
	t, err = template.New("").Delims(cfg.tmplLDelimiter, cfg.tmplRDelimiter).Funcs(map[string]interface{}{
		"Title": strings.Title,
	}).Parse(tmpl)
	if err != nil {
		log.Fatalf("[Error] Template Parse Error: %s", err)
	}

	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("[Error] Create %s Error: %s", filename, err)
	}
	defer f.Close()
	err = t.Execute(f, data)
	if err != nil {
		log.Fatalf("[Error] Execute %s Error: %s", filename, err)
	}
	log.Printf("[Success] Generate %s Success", filename)
	if cfg.gofmt {
		GoFmt(filename)
	}
}

func GoFmt(filename string) {
	filename = goSrcFileName(filename)
	err := exec.Command("gofmt", "-w", filename).Start()
	if err != nil {
		log.Printf("[Warn] gofmt %s Error: %s\n", filename, err)
		return
	}
	log.Printf("[Success] gofmt %s Success\n", filename)
}
