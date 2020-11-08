package gen

type Option func(cfg *config)

type config struct {
	gofmt                          bool
	tmplLDelimiter, tmplRDelimiter string
}

var defaultCfg = config{
	gofmt:          true,
	tmplLDelimiter: "<<",
	tmplRDelimiter: ">>",
}

func (cfg *config) use(opts ...Option) {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(cfg)
	}
}

func newConfig(opts ...Option) config {
	cfg := defaultCfg
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
