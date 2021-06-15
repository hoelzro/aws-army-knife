package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type example struct {
	name           string
	input          []byte
	expectedOutput []byte
}

func TestExamples(t *testing.T) {
	exampleFiles, err := os.ReadDir("examples")
	if err != nil {
		t.Fatal(err)
	}

	var examples []example

	for _, file := range exampleFiles {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		exampleInputBytes, err := os.ReadFile(filepath.Join("examples", file.Name()))
		if err != nil {
			t.Fatal(err)
		}

		testName := file.Name()
		testName = strings.TrimSuffix(testName, filepath.Ext(testName))

		examples = append(examples, example{
			name:  testName,
			input: exampleInputBytes,
		})
	}

	for _, ex := range examples {
		t.Run(ex.name, func(t *testing.T) {
			exampleOutput := &strings.Builder{}

			args := []string{}
			err = awsNames(bytes.NewReader(ex.input), exampleOutput, args)
			if err != nil {
				t.Fatal(err)
			}
			// XXX test got output matches expected
		})
	}
}
