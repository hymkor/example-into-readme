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
	"regexp"
	"runtime"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/transform"

	"github.com/hymkor/example-into-readme/internal/realpath"
	"github.com/hymkor/example-into-readme/outline"
)

var (
	rxCodeBlock = regexp.MustCompile("^```")
	rxComment   = regexp.MustCompile(`^<!--\s+-->`)
	rxMarker    = regexp.MustCompile(`^<!--\s*(.*?)\s*-->\s*$`)
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

func skipUntil(br *bufio.Reader, rx *regexp.Regexp, w io.Writer) error {
	for {
		line, err := br.ReadString('\n')
		io.WriteString(w, line)
		if rx.MatchString(line) {
			return nil
		}
		if err != nil {
			if err == io.EOF {
				return io.ErrUnexpectedEOF
			}
			return err
		}
	}
}

func copyWithNoDetab(r io.Reader, newline string, w io.Writer) error {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		io.WriteString(w, sc.Text())
		io.WriteString(w, newline)
	}
	return sc.Err()
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

func splitField(s string) (result []string) {
	for len(s) > 0 {
		for len(s) > 0 && strings.IndexByte(" \t\v\r\n", s[0]) >= 0 {
			s = s[1:]
		}
		quote := false
		var buffer strings.Builder
		for len(s) > 0 {
			c, siz := utf8.DecodeRuneInString(s)
			if !quote && strings.ContainsRune(" \t\v\r\n", c) {
				break
			}
			if c == '"' {
				quote = !quote
			} else {
				buffer.WriteString(s[:siz])
			}
			s = s[siz:]
		}
		result = append(result, buffer.String())
	}
	return
}

func open(s string) (io.ReadCloser, error) {
	if len(s) > 0 && s[len(s)-1] == '|' {
		args := splitField(s[:len(s)-1])
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
	}
	fd, err := os.Open(s)
	if os.IsNotExist(err) {
		// remove language text and retry
		_, s, ok := strings.Cut(s, " ")
		if !ok {
			return nil, err
		}
		fd, err = os.Open(s)
	}
	return fd, err
}

func filter(r io.Reader, w io.Writer, headers []*outline.Header, log func(...any)) error {
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
				if err = skipUntil(br, rxCodeBlock, bw); err != nil {
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
			} else if strings.EqualFold(filename, "Makefile") {
				copyWithNoDetab(qr, newline, bw)
			} else {
				copyWithDetab(qr, newline, bw)
			}

			qr.Close()
			if err := skipUntil(br, rxCodeBlock, io.Discard); err != nil {
				return err
			}
			bw.WriteString("```")
			bw.WriteString(newline)
			log("Include", filename)
		} else if m := rxMarker.FindStringSubmatch(text); m != nil {
			newline := "\n"
			if strings.HasSuffix(text, "\r\n") {
				newline = "\r\n"
			}
			if m[1] == "outline" {
				bw.WriteString(newline)
				outline.List(headers, "", newline, bw)
				bw.WriteString(newline)
				bw.WriteString("<!-- -->")
				bw.WriteString(newline)
				if err := skipUntil(br, rxComment, io.Discard); err != nil {
					return err
				}
				log("Make Outline")
			} else if fd, err := open(m[1]); err == nil {
				sc := bufio.NewScanner(fd)
				for sc.Scan() {
					io.WriteString(bw, sc.Text())
					io.WriteString(bw, newline)
				}
				if err := fd.Close(); err != nil {
					return err
				}
				if err := sc.Err(); err != nil {
					return err
				}
				bw.WriteString("<!-- -->")
				if err := skipUntil(br, rxComment, io.Discard); err != nil {
					return err
				}
				bw.WriteString(newline)
				log("Include", m[1])
			}
		}
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			}
			return nil
		}
	}
}

func conv(srcFile, dstFile string, log func(...any)) (string, error) {
	r, err := os.Open(srcFile)
	if err != nil {
		return "", err
	}
	defer r.Close()

	realName, err := realpath.FromFile(r)
	if err != nil {
		return "", err
	}

	headers, err := outline.Make(srcFile)
	if err != nil {
		return "", err
	}

	w, err := os.Create(dstFile)
	if err != nil {
		return "", err
	}
	defer w.Close()

	return realName, filter(r, w, headers, log)
}

var (
	flagTarget = flag.String("target", "README.md", "Rewrite filename (Deprecated: remove `-target`)")
	flagTemp   = flag.String("temporary", "{}.tmp", "Temporary filename ({} means original filepath)")
	flagBackup = flag.String("backup", "{}~", "Backup filename ({} means original filepath)")
)

func logToStderr(s ...any) {
	fmt.Fprintln(os.Stderr, s...)
}

func mains(args []string) error {
	const lockKey = "EXAMPLEINTOREADME"
	_, ok := os.LookupEnv(lockKey)
	if ok {
		return errors.New("Locked")
	}
	os.Setenv(lockKey, "RUNNING")

	md := *flagTarget
	if len(args) >= 1 {
		md = args[0]
	}
	tmp := strings.Replace(*flagTemp, "{}", md, 1)

	fmt.Fprintln(os.Stderr, "Convert from", md, "to", tmp)
	var err error
	md, err = conv(md, tmp, logToStderr)
	if err != nil {
		return err
	}

	bak := strings.Replace(*flagBackup, "{}", md, 1)

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

var version string

func main() {
	fmt.Fprintf(os.Stderr, "%s %s-%s-%s\n", os.Args[0], version, runtime.GOOS, runtime.GOARCH)
	flag.Parse()
	if err := mains(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
