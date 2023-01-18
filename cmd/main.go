package main

import (
	"embed"
	"os"
	"strings"

	"github.com/frain-dev/convoy"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//go:embed VERSION
var f embed.FS

func main() {
	err := os.Setenv("TZ", "") // Use UTC by default :)
	if err != nil {
		log.Fatal("failed to set TZ env - ", err)
	}

	cmd := &cobra.Command{
		Use:     "Convoy CLI",
		Version: convoy.GetVersion(),
		Short:   "Convoy CLI for debugging your events locally",
	}

	cmd.AddCommand(addListenCommand())
	cmd.AddCommand(addLoginCommand())

	err = cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
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
