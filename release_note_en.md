v0.5.0
======
May 17, 2024

- Change the default value of `-temporary` which was `README.tmp` to `{}.tmp`,  
  and the default value of `-backup` which was  `README.md~` to `{}~`.  
  (`{}` means the value of the target filename)
- The word `-target` can be omitted.

v0.4.0
=======
Oct 2, 2023

Support ` ```{LANGNAME} FILENAME` as the format of the header of codeblocks

- `{LANGCODE}` can be omited  
    for example: ` ```rust dir/foo.rs`
- This modify is done because GitHub can not consider ` ```foo.rs` as codeblock for Rust.

v0.3.1
=======
Aug 24, 2023

- Use os.Pipe for both STDOUT and STDIN instead of cmd.StdoutPipe and cmd.StderrPipe()

v0.3.0
=======
Aug 24, 2023

- Merge the contents of STDOUT and STDERR of ` ```command`

v0.2.0
=======
May 1, 2023

+ Support quoting the output of the command like `` ```COMMAND | ``

v0.1.0
=======
Mar 23, 2023

+ The first release
