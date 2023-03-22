package main

import (
	"strings"
	"testing"
)

func TestFilterLF(t *testing.T) {
	var output strings.Builder
	source := strings.NewReader("foo\n```testdata.txt\n```\n")
	expect := "foo\n```testdata.txt\nhogehoge\n```\n"

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
	source := strings.NewReader("foo\r\n```testdata.txt\r\n```\r\n")
	expect := "foo\r\n```testdata.txt\r\nhogehoge\r\n```\r\n"

	err := filter(source, &output, func(...any) {})
	if err != nil {
		t.Fatal(err.Error())
	}
	result := output.String()
	if expect != result {
		t.Fatalf("expect `%s` but `%s`", expect, result)
	}
}
