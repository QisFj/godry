package config

import "sync"

// zero value is ready to use
type MultiAccessor struct {
	rw   sync.RWMutex
	aMap map[string]*Accessor
}

func (ma *MultiAccessor) AddAccessor(key string, a *Accessor) {
	ma.rw.Lock()
	defer ma.rw.Unlock()
	if ma.aMap == nil {
		ma.aMap = map[string]*Accessor{}
	}
	ma.aMap[key] = a
}

func (ma *MultiAccessor) GetAccessor(key string) (*Accessor, bool) {
	ma.rw.RLock()
	defer ma.rw.RUnlock()
	a, exist := ma.aMap[key]
	return a, exist
}

func (ma *MultiAccessor) Accessor(key string) *Accessor {
	a, _ := ma.GetAccessor(key)
	return a
}

func (ma *MultiAccessor) GetConfig(key string) (Config, bool) {
	if a, exist := ma.GetAccessor(key); exist {
		return a.Config(), true
	}
	return Config{}, false
}

func (ma *MultiAccessor) Config(key string) Config {
	c, _ := ma.GetConfig(key)
	return c
}

func (ma *MultiAccessor) Configs() map[string]Config {
	ma.rw.RLock()
	defer ma.rw.RUnlock()
	configs := map[string]Config{}
	for key, a := range ma.aMap {
		configs[key] = a.Config()
	}
	return configs
}
