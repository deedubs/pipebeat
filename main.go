package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/deedubs/pipebeat/beater"
)

func main() {
	err := beat.Run("pipebeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
