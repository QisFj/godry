package gen

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

func GoSrcFileName(filename string) string {
	if !strings.HasSuffix(filename, ".go") {
		filename += ".go"
	}
	return filename
}

func Gen(filename, tmpl string, data interface{}, options ...Option) {
	cfg := newConfig(options...)
	filename = GoSrcFileName(filename)
	var t *template.Template
	var err error
	t, err = template.New("").Delims(cfg.tmplLDelimiter, cfg.tmplRDelimiter).Funcs(map[string]interface{}{
		"Title": strings.Title,
	}).Parse(tmpl)
	if err != nil {
		log.Fatalf("[Error] template parse error: %s", err)
	}

	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("[Error] create %s Error: %s", filename, err)
	}
	defer f.Close()
	err = t.Execute(f, data)
	if err != nil {
		log.Fatalf("[Error] execute %s Error: %s", filename, err)
	}
	log.Printf("[Success] generate %s success", filename)
	if cfg.gofmt {
		GoFmt(filename)
	}
}

func GoFmt(filename string) {
	filename = GoSrcFileName(filename)
	err := exec.Command("gofmt", "-w", filename).Start()
	if err != nil {
		log.Printf("[Warn] gofmt %s error: %s\n", filename, err)
		return
	}
	log.Printf("[Success] gofmt %s success\n", filename)
}
