package main

import (
	"strings"
	"testing"
)

func TestFilterLF(t *testing.T) {
	var output strings.Builder
	source := strings.NewReader("foo\n```go.mod\n```\n")
	expect := "foo\n```go.mod\nmodule github.com/hymkor/example-into-readme\n\ngo 1.20\n```\n"

	err := filter(source, &output, func(...any) {})
	if err != nil {
		t.Fatal(err.Error())
	}
	result := output.String()
	if expect != result {
		t.Fatalf("expect `%s` but `%s`", expect, result)
	}
}

func TestFilterCRLF(t *testing.T) {
	var output strings.Builder
	source := strings.NewReader("foo\r\n```go.mod\r\n```\r\n")
	expect := "foo\r\n```go.mod\r\nmodule github.com/hymkor/example-into-readme\r\n\r\ngo 1.20\r\n```\r\n"

	err := filter(source, &output, func(...any) {})
	if err != nil {
		t.Fatal(err.Error())
	}
	result := output.String()
	if expect != result {
		t.Fatalf("expect `%s` but `%s`", expect, result)
	}
}
