package config

import (
	"encoding/json"
	"regexp"
	"time"
)

// types for json marshal and unmarshal
// - Regexp : *regexp.Regexp
// - Duration : time.Duration

type Regexp struct {
	*regexp.Regexp // note: zero value is nil, not ""
}

func (reg Regexp) MarshalJSON() ([]byte, error) {
	if reg.Regexp == nil {
		// don't return ""
		// "" is not zero value
		return []byte(`null`), nil
	}
	return json.Marshal(reg.String())
}

func (reg *Regexp) UnmarshalJSON(data []byte) error {
	if len(data) == 4 && string(data) == "null" {
		reg.Regexp = nil // set to nil
		return nil
	}
	var rStr string
	err := json.Unmarshal(data, &rStr)
	if err != nil {
		return err
	}
	r, err := regexp.Compile(rStr)
	if err != nil {
		return err
	}
	reg.Regexp = r
	return nil
}

type Duration time.Duration

func (dur Duration) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Duration(dur).String() + "\""), nil
}

func (dur *Duration) UnmarshalJSON(data []byte) error {
	var dStr string
	err := json.Unmarshal(data, &dStr)
	if err != nil {
		return err
	}
	if dStr == "" {
		*dur = 0 // zero value
		return nil
	}
	d, err := time.ParseDuration(dStr)
	if err != nil {
		return err
	}
	*dur = Duration(d)
	return nil
}
