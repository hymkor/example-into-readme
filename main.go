package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/text/transform"
)

type goFilter struct {
	pass bool
}

func (g *goFilter) Reset() {}

func (g *goFilter) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	if g.pass {
		n := copy(dst, src)
		return n, n, nil
	}
	for {
		newlinePos := bytes.IndexByte(src, '\n')
		if newlinePos < 0 && !atEOF {
			return nDst, nSrc, transform.ErrShortSrc
		}
		if bytes.HasPrefix(src, []byte("package")) {
			if len(dst) < newlinePos+1 {
				return nDst, nSrc, transform.ErrShortDst
			}
			n := copy(dst, src)
			nSrc += n
			nDst += n
			g.pass = true
			return
		}
		nSrc += newlinePos + 1
		src = src[newlinePos+1:]
	}
}

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

func open(s string) (io.ReadCloser, error) {
	if len(s) > 0 && s[len(s)-1] == '|' {
		args := strings.Fields(s[:len(s)-1])
		cmd := exec.Command(args[0], args[1:]...)
		r, w, err := os.Pipe()
		if err != nil {
			return nil, err
		}
		cmd.Stdin = os.Stdin
		cmd.Stdout = w
		cmd.Stderr = w
		err = cmd.Start()
		if err != nil {
			return r, err
		}
		go func() {
			cmd.Wait()
			w.Close()
		}()
		return r, nil
	} else {
		return os.Open(s)
	}
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
			qr, err := open(filename)
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
			if strings.HasSuffix(filename, ".go") {
				copyWithDetab(transform.NewReader(qr, &goFilter{}), newline, bw)
			} else {
				copyWithDetab(qr, newline, bw)
			}

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
	const lockKey = "EXAMPLEINTOREADME"
	_, ok := os.LookupEnv(lockKey)
	if ok {
		return errors.New("Locked")
	}
	os.Setenv(lockKey, "RUNNING")

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
