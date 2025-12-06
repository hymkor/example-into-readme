package main

import (
	"strings"
	"testing"
)

func TestFilterLF(t *testing.T) {
	var output strings.Builder
	source := strings.NewReader("foo\n```testdata.txt\n```\n")
	expect := "foo\n```testdata.txt\nhogehoge\n```\n"

	err := filter(source, &output, nil, func(...any) {})
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

	err := filter(source, &output, nil, func(...any) {})
	if err != nil {
		t.Fatal(err.Error())
	}
	result := output.String()
	if expect != result {
		t.Fatalf("expect `%s` but `%s`", expect, result)
	}
}

func testSplitField(t *testing.T, source string, expect ...string) {
	t.Helper()
	result := splitField(source)
	if len(expect) != len(result) {
		t.Fatalf("%#v: len: expect %d, but %d", source, len(expect), len(result))
	}
	for i := range expect {
		if expect[i] != result[i] {
			t.Fatalf("%#v: [%d]: expect %#v, but %#v", source, i, expect[i], result[i])
		}
	}
}

func TestSplitField(t *testing.T) {
	testSplitField(t, `foo bar  baz`, "foo", "bar", "baz")
	testSplitField(t, `foo " bar " baz`, "foo", " bar ", "baz")
}
