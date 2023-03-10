package convoy_cli

import (
	"embed"
	"strings"
)

//go:embed VERSION
var f embed.FS

func GetVersion() string {
	v := "0.1.0"

	f, err := ReadVersion()
	if err != nil {
		return v
	}

	v = strings.TrimSuffix(string(f), "\n")
	return v
}

func ReadVersion() ([]byte, error) {
	data, err := f.ReadFile("VERSION")
	if err != nil {
		return nil, err
	}

	return data, nil
}
