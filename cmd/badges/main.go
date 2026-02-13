package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var rxVersion = regexp.MustCompile(`/v[1-9]$`)

func module() (string, error) {
	fd, err := os.Open("go.mod")
	if err != nil {
		return "", err
	}
	defer fd.Close()

	sc := bufio.NewScanner(fd)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "module") {
			module := strings.TrimSpace(line[6:])
			gitUrl := rxVersion.ReplaceAllString(module, "")
			return gitUrl, nil
		}
	}
	return "", sc.Err()
}

func mains() error {
	url, err := module()
	if err != nil {
		return err
	}
	if _, err := os.Stat(`.github\workflows\go.yml`); err == nil {
		fmt.Printf("[![Go Test](https://%[1]s/actions/workflows/go.yml/badge.svg)](https://%[1]s/actions/workflows/go.yml)\n", url)
	}
	if _, err := os.Stat("LICENSE"); err == nil {
		fmt.Printf("[![License](https://img.shields.io/badge/License-MIT-red)](https://%[1]s/blob/master/LICENSE)\n", url)
	}
	fmt.Printf("[![Go Reference](https://pkg.go.dev/badge/%[1]s.svg)](https://pkg.go.dev/%[1]s)\n", url)

	fmt.Printf("[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://%s)\n", url)

	return nil
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
