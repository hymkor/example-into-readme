Release notes
=============

- Fix: outline generation broken when headers contained links (`[text](url)`); now only the link text is used (#6)
- Add new directives: `stdout:` and `output:` (#7)

v0.8.2
------
Feb 3, 2026

- Fixed an issue where, when the target file was omitted, a file whose name differed from `README.md` only in letter case (e.g. `readme.md`) could be updated with its name changed to `README.md`. (#3)

v0.8.1
------
Dec 7, 2025

- Fix: Incorrect handling of arguments enclosed in double quotes when embedding command output. (#1)
- Adjust line endings in command output to match the line endings of the target Markdown. (#2)

v0.8.0
------
Oct 22, 2025

- When a line like `<!-- command-name args…| -->` is found, the tool now replaces the lines from the next line up to the nearest `<!-- -->` with the command’s output.

v0.7.1
------
Oct 2, 2025

- Made the tool exit with an error if the closing line for an outline or code block replacement is missing.

v0.7.0
------
Sep 28, 2025

- Support the outline generator
    - It generates between `<!-- outline -->` and `<!-- -->`

v0.6.0
------
Dec 13, 2024

- Tabs are not expanded to spaces now when Makefile is quoted

v0.5.0
------
May 17, 2024

- Change the default value of `-temporary` which was `README.tmp` to `{}.tmp`,  
  and the default value of `-backup` which was  `README.md~` to `{}~`.  
  (`{}` means the value of the target filename)
- The word `-target` can be omitted.

v0.4.0
------
Oct 2, 2023

Support ` ```{LANGNAME} FILENAME` as the format of the header of codeblocks

- `{LANGCODE}` can be omited  
    for example: ` ```rust dir/foo.rs`
- This modify is done because GitHub can not consider ` ```foo.rs` as codeblock for Rust.

v0.3.1
------
Aug 24, 2023

- Use os.Pipe for both STDOUT and STDIN instead of cmd.StdoutPipe and cmd.StderrPipe()

v0.3.0
------
Aug 24, 2023

- Merge the contents of STDOUT and STDERR of ` ```command`

v0.2.0
------
May 1, 2023

- Support quoting the output of the command like `` ```COMMAND | ``

v0.1.0
------
Mar 23, 2023

- The first release
