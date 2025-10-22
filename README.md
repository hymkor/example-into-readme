example-into-readme
===================

<!-- outline -->

- [example-into-readme](#example-into-readme)
    - [Specify the language for codeblock](#specify-the-language-for-codeblock)
    - [Quoting the output of the command](#quoting-the-output-of-the-command)
    - [Outline generator](#outline-generator)
    - [Include other markdown](#include-other-markdown)
    - [Install](#install)
        - [Use go install](#use-go-install)
        - [Use the scoop-installer](#use-the-scoop-installer)

<!-- -->

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

Specify the language for codeblock
----------------------------------

On GitHub, the codeblock for some languages can not be judged with their extensions only. Then, you can write the name of languge before filename.

    ```LANGNAME FILENAME
    ```

For example:

    ```rust foo.rs
    ```

Quoting the output of the command
-----

    ```COMMANDNAME ARGS ... |
    ```

Outline generator
-----------------

    <!-- outline -->
    <!-- -->

Include other markdown
----------------

    <!-- example.md -->
    <!-- -->

Include command output
----------------------

    <!-- COMMANDNAME ARGS ... | -->
    <!-- -->

Install
-------

Download the binary package from [Releases](https://github.com/hymkor/example-into-readme/releases) and extract the executable.


### Use go install

```
go install github.com/hymkor/example-into-readme@latest
```

### Use the scoop-installer

```
scoop install https://raw.githubusercontent.com/hymkor/example-into-readme/master/example-into-readme.json
```

or

```
scoop bucket add hymkor https://github.com/hymkor/scoop-bucket
scoop install example-into-readme
```
