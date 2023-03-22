package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func readUntilQQQ(sc *bufio.Scanner, w io.Writer) error {
	for sc.Scan() {
		fmt.Fprintln(w, sc.Text())
		if strings.HasPrefix(sc.Text(), "```") {
			return nil
		}
	}
	return sc.Err()
}

func copyWithDetab(r io.Reader, w io.Writer) error {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		text := sc.Text()
		for len(text) > 0 && text[0] == '\t' {
			io.WriteString(w, "    ")
			text = text[1:]
		}
		fmt.Fprintln(w, text)
	}
	return sc.Err()
}

func conv(srcFile, dstFile string, log func(...any)) error {
	r, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer w.Close()

	bw := bufio.NewWriter(w)
	defer bw.Flush()

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		text := sc.Text()
		fmt.Fprintln(bw, text)
		if strings.HasPrefix(text, "```") {
			filename := strings.TrimSpace(text[3:])
			qr, err := os.Open(filename)
			if err != nil {
				if !os.IsNotExist(err) {
					return err
				}
				if err = readUntilQQQ(sc, bw); err != nil {
					return err
				}
				continue
			}
			copyWithDetab(qr, bw)
			qr.Close()
			if err := readUntilQQQ(sc, io.Discard); err != nil {
				return err
			}
			bw.WriteString("```\n")
			log("Include", filename)
		}
	}
	return sc.Err()
}

func mains() error {
	const md = "README.md"
	const tmp = "README.tmp"
	const bak = "README.md~"

	fmt.Fprintln(os.Stderr, "Convert from", md, "to", tmp)
	if err := conv(md, tmp, func(s ...any) { fmt.Fprintln(os.Stderr, s...) }); err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "Rename", md, "to", bak)
	if err := os.Rename(md, bak); err != nil {
		return fmt.Errorf("rename `%s` to `%s`: %w", md, bak, err)
	}
	fmt.Fprintln(os.Stderr, "Rename", tmp, "to", md)
	if err := os.Rename(tmp, md); err != nil {
		return fmt.Errorf("rename `%s` to `%s`: %w", tmp, md, err)
	}
	return nil
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
