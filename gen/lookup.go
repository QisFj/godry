package gen

import (
	"fmt"
	"os"
	"path/filepath"
)

func GoFileDir() (string, error) {
	goFile, ok := os.LookupEnv("GOFILE")
	if !ok {
		return "", fmt.Errorf("environment variable GOFILE is not set")
	}
	return filepath.Dir(goFile), nil
}
