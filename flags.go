package main

import (
	"strings"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	res := strings.Join([]string(*i), " ")
	return "[" + res + "]"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *arrayFlags) Map() map[string]string {
	result := make(map[string]string)
	for _, e := range []string(*i) {
		keyValuePair := strings.SplitN(e, "=", 2)
		result[keyValuePair[0]] = keyValuePair[1]
	}
	return result
}
