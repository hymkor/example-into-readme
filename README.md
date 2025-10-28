example-into-readme
===================

<!-- badges.cmd | -->
[![Go Test](https://github.com/hymkor/example-into-readme/actions/workflows/go.yml/badge.svg)](https://github.com/hymkor/example-into-readme/actions/workflows/go.yml)
[![License](https://img.shields.io/badge/License-MIT-red)](https://github.com/hymkor/example-into-readme/blob/master/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/hymkor/example-into-readme.svg)](https://pkg.go.dev/github.com/hymkor/example-into-readme)
<!-- -->

<!-- outline -->

- [example-into-readme](#example-into-readme)
    - [Overview](#overview)
        - [Before running](#before-running)
        - [After running](#after-running)
    - [Usage examples](#usage-examples)
        - [Include a file](#include-a-file)
        - [Include command output](#include-command-output)
        - [Include another markdown file](#include-another-markdown-file)
        - [Specify the language for code blocks](#specify-the-language-for-code-blocks)
        - [Quoting the output of commands](#quoting-the-output-of-commands)
        - [Outline generator](#outline-generator)
        - [Include other markdown files](#include-other-markdown-files)
        - [Include command output](#include-command-output-1)
    - [Install](#install)
        - [Use go install](#use-go-install)
        - [Use the scoop-installer](#use-the-scoop-installer)

<!-- -->

Overview
--------

**example-into-readme** automatically inserts example files or command outputs into your `README.md`.
It helps you keep your documentation up-to-date — without manually copying and pasting code or results.

This tool is meant to be used **locally** as a CLI helper, not as a GitHub Action.

### Before running

    ```go.mod
    ```

### After running

    ```go.mod
    module github.com/hymkor/example-into-readme
    
    go 1.20
    
    require golang.org/x/text v0.8.0
    ```

The program finds code blocks whose info-string contains a filename, then replaces the block’s content with the actual file content.

```
$ example-into-readme
Convert from README.md to README.tmp
Include go.mod
Rename README.md to README.md~
Rename README.tmp to README.md
```

When the filename ends with `*.go`, it skips lines before `package` to ignore `//go:build` directives.

Usage examples
--------------

### Include a file

    ```go.mod
    ```

→ replaced with the actual contents of `go.mod`.

### Include command output

    ```go run example.go |
    ```

→ replaced with the result of running the command.

### Include another markdown file

    <!-- example.md -->
    <!-- -->

→ replaced with the contents of `example.md`.

### Specify the language for code blocks

Some file extensions are not automatically recognized by GitHub’s syntax highlighting.
You can specify the language name explicitly before the filename:

    ```rust foo.rs
    ```

### Quoting the output of commands

To embed command results, write the command followed by a `|`:

    ```COMMANDNAME ARGS ... |
    ```

### Outline generator

You can insert a markdown outline between these markers:

    <!-- outline -->
    <!-- -->

### Include other markdown files

    <!-- filename.md -->
    <!-- -->

### Include command output

    <!-- COMMANDNAME ARGS ... | -->
    <!-- -->

Install
-------

Download the binary from [Releases](https://github.com/hymkor/example-into-readme/releases) and extract the executable.

### Use go install

```cmd
go install github.com/hymkor/example-into-readme@latest
```

### Use the scoop-installer

```cmd
scoop install https://raw.githubusercontent.com/hymkor/example-into-readme/master/example-into-readme.json
```

or

```cmd
scoop bucket add hymkor https://github.com/hymkor/scoop-bucket
scoop install example-into-readme
```
