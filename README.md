example-into-readme
===================

This program inserts example-files into `README.md` at the code block in the current directory.

```go.mod
module github.com/hymkor/example-into-readme

go 1.20
```

The info-string at the header of codeblocks has to have a filename to include the file.
When a filename is not written, the block will not be changed.

```
$ ./example-into-readme.exe
Convert from README.md to README.tmp
Include go.mod
Rename README.md to README.md~
Rename README.tmp to README.md
```
