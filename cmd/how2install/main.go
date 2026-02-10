package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//go:embed en.txt
var enText string

//go:embed ja.txt
var jaText string

var rxVersion = regexp.MustCompile(`/v[1-9]$`)

type Module struct {
	module string
	url    string
	user   string
	repo   string
}

func module() (*Module, error) {
	fd, err := os.Open("go.mod")
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	sc := bufio.NewScanner(fd)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "module") {
			module := strings.TrimSpace(line[6:])
			gitUrl := rxVersion.ReplaceAllString(module, "")
			parts := strings.Split(gitUrl, "/")
			if len(parts) < 3 {
				return nil, fmt.Errorf("%s: invalid module name", module)
			}
			return &Module{
				module: module,
				url:    gitUrl,
				user:   parts[1],
				repo:   parts[2],
			}, nil
		}
	}
	return nil, sc.Err()
}

func mains(args []string) error {
	mod, err := module()
	if err != nil {
		return err
	}

	cmddir := ""
	if _, err := os.Stat(filepath.Join("cmd", mod.repo)); err == nil {
		cmddir = "/cmd/" + mod.repo
	}
	if len(args) >= 1 && strings.EqualFold(args[0],"ja") {
		fmt.Printf(jaText, mod.user, mod.repo, mod.module, mod.url, cmddir)
	} else {
		fmt.Printf(enText, mod.user, mod.repo, mod.module, mod.url, cmddir)
	}
	return nil
}

func main() {
	if err := mains(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
