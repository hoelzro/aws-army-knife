package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var exceptions = map[string]string{
	"SecurityGroups":           "GroupName",
	"LoadBalancerDescriptions": "LoadBalancerName",
	"SecretList":               "Name",
}

func awsNames(r io.Reader, w io.Writer, args []string) error {
	input, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	var doc map[string]interface{}
	err = json.Unmarshal(input, &doc)
	if err != nil {
		return err
	}

	var key string
	for key = range doc {
		break
	}

	rawItems := doc[key].([]interface{})
	items := make([]map[string]interface{}, len(rawItems))

	for i, item := range rawItems {
		items[i] = item.(map[string]interface{})
	}

	if key == "Reservations" {
		key = "Instances"
		expandedItems := make([]map[string]interface{}, 0)

		for _, item := range items {
			instances := item["Instances"].([]map[string]interface{})
			expandedItems = append(expandedItems, instances...)
		}

		items = expandedItems
	}

	nameKeyPrefix := key
	if strings.HasSuffix(nameKeyPrefix, "ies") {
		nameKeyPrefix = strings.TrimSuffix(nameKeyPrefix, "ies") + "ys"
	}
	if strings.HasSuffix(nameKeyPrefix, "s") {
		nameKeyPrefix = strings.TrimSuffix(nameKeyPrefix, "s")
	}

	if len(items) == 0 {
		return nil
	}

	nameKey, isExceptional := exceptions[key]
	if !isExceptional {
		nameKeys := make([]string, 0)
		for k := range items[0] {
			if strings.EqualFold(k, nameKeyPrefix+"name") {
				nameKeys = append(nameKeys, k)
			}
			if strings.EqualFold(k, nameKeyPrefix+"id") {
				nameKeys = append(nameKeys, k)
			}
			if strings.EqualFold(k, nameKeyPrefix+"identifier") {
				nameKeys = append(nameKeys, k)
			}
		}

		if len(nameKeys) > 1 {
			sort.Strings(nameKeys)
			return errors.New("Ambiguous name keys: " + strings.Join(nameKeys, " "))
		} else if len(nameKeys) == 0 {
			return errors.New("Unable to determine name key for " + key)
		}

		nameKey = nameKeys[0]
	}

	for _, item := range items {
		fmt.Fprintln(w, item[nameKey])
	}

	return nil
}

func main() {
	invokedAs := filepath.Base(os.Args[0])
	args := os.Args[1:]

	if invokedAs == "aws-army-knife" {
		invokedAs = os.Args[1]
		args = os.Args[2:]
	}

	switch invokedAs {
	case "aws-names":
		err := awsNames(os.Stdin, os.Stdout, args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Printf("I don't know how to run as %s\n", invokedAs)
	}
}
