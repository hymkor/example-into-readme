example-into-readme
===================

This program inserts example-files into `README.md` at the code block in the current directory.

`README.md` before running

    ```go.mod
    ```

`README.md` after running

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

When `*.go` is given as a filename, skip lines until `package` is found to ignore `//go:build`.

Quoting the output of the command
-----

    ```COMMANDNAME ARGS ... |
    ```

Install
-------

Download the binary package from [Releases](https://github.com/hymkor/example-into-readme/releases) and extract the executable.

### for scoop-installer

```
scoop install https://raw.githubusercontent.com/hymkor/example-into-readme/master/example-into-readme.json
```

or

```
scoop bucket add hymkor https://github.com/hymkor/scoop-bucket
scoop install example-into-readme
```
