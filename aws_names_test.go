package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExamples(t *testing.T) {
	exampleInputs, err := os.ReadDir("examples")
	if err != nil {
		t.Fatal(err)
	}

	for _, ex := range exampleInputs {
		testName := ex.Name()
		t.Run(testName, func(t *testing.T) {
			exampleInput, err := os.Open(filepath.Join("examples", ex.Name()))
			if err != nil {
				t.Fatal(err)
			}

			defer exampleInput.Close()

			exampleOutput := &strings.Builder{}

			args := []string{}
			err = awsNames(exampleInput, exampleOutput, args)
			if err != nil {
				t.Fatal(err)
			}
			// XXX test got output matches expected
		})
	}
}
