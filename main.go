package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func readUntilQQQ(br *bufio.Reader, w io.Writer) error {
	for {
		line, err := br.ReadString('\n')
		io.WriteString(w, line)
		if strings.HasPrefix(line, "```") {
			return nil
		}
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
	}
}

func copyWithDetab(r io.Reader, newline string, w io.Writer) error {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		text := sc.Text()
		for len(text) > 0 && text[0] == '\t' {
			io.WriteString(w, "    ")
			text = text[1:]
		}
		io.WriteString(w, text)
		io.WriteString(w, newline)
	}
	return sc.Err()
}

func filter(r io.Reader, w io.Writer, log func(...any)) error {
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	br := bufio.NewReader(r)
	for {
		text, errRead := br.ReadString('\n')
		io.WriteString(bw, text)
		if strings.HasPrefix(text, "```") {
			filename := strings.TrimSpace(text[3:])
			qr, err := os.Open(filename)
			if err != nil {
				if !os.IsNotExist(err) {
					return err
				}
				if err = readUntilQQQ(br, bw); err != nil {
					return err
				}
				continue
			}
			newline := "\n"
			if strings.HasSuffix(text, "\r\n") {
				newline = "\r\n"
			}
			copyWithDetab(qr, newline, bw)
			qr.Close()
			if err := readUntilQQQ(br, io.Discard); err != nil {
				return err
			}
			bw.WriteString("```")
			bw.WriteString(newline)
			log("Include", filename)
		}
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			}
			return nil
		}
	}
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

	return filter(r, w, log)
}

var (
	flagTarget = flag.String("target", "README.md", "Rewrite filename")
	flagTemp   = flag.String("temporary", "README.tmp", "Temporary filename")
	flagBackup = flag.String("backup", "README.md~", "Backup filename")
)

func mains() error {
	md := *flagTarget
	tmp := *flagTemp
	bak := *flagBackup

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
	flag.Parse()
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
