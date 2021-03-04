package gen

import (
	"strings"
	"text/template"
)

type Option func(cfg *config)

type config struct {
	gofmt                          bool
	tmplLDelimiter, tmplRDelimiter string
	funcMap                        template.FuncMap
}

var defaultCfg = config{
	gofmt:          true,
	tmplLDelimiter: "<<",
	tmplRDelimiter: ">>",
	funcMap: map[string]interface{}{
		"Title": strings.Title,
	},
}

func (cfg *config) use(opts ...Option) {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(cfg)
	}
}

func (cfg *config) deepCopy() config {
	newCfg := *cfg
	newCfg.funcMap = map[string]interface{}{}
	for k, v := range cfg.funcMap {
		newCfg.funcMap[k] = v
	}
	return newCfg
}

func newConfig(opts ...Option) config {
	cfg := defaultCfg.deepCopy()
	cfg.use(opts...)
	return cfg
}

// to support using code like `Options.XXX()` to create a Option
type optionFactory struct{}

var Options optionFactory

func (optionFactory) GoFmt(enable bool) Option {
	return func(cfg *config) {
		cfg.gofmt = enable
	}
}

func (optionFactory) TemplateDelimiter(left, right string) Option {
	return func(cfg *config) {
		cfg.tmplLDelimiter = left
		cfg.tmplRDelimiter = right
	}
}

func (optionFactory) Func(key string, f interface{}) Option {
	return func(cfg *config) {
		// assert: cfg.funcMap != nil
		cfg.funcMap[key] = f
	}
}
