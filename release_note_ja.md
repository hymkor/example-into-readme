- `README.tmp` だった `-temporary` のデフォルト値を `{}.tmp` に、  
  `README.md~` だった `-backup` のデフォルト値を `{}~` に  
  変更した(`{}` はターゲットファイル名) 
- ターゲットファイル名を指定する `-target` の語句を省略できるようにした。

v0.4.0
=======
Oct 2, 2023

` ```{LANGNAME} FILENAME` をコードブロックのヘッダフォーマットとしてサポートしました。

- `{LANGCODE}` は省略可能です。  
    例： ` ```rust dir/foo.rs`
- GitHub が ` ```foo.rs` を Rust のコードブロックとして判断してくれなかったため、この修正を行いました。

v0.3.1
=======
Aug 24, 2023

- cmd.StdoutPipe() と cmd.StderrPipe() のかわりに STDOUT,STDERR の両方に os.Pipe を使うようにした。

v0.3.0
=======
Aug 24, 2023

- ` ```command` の標準出力と標準エラー出力をマージするようにした

v0.2.0
=======
May 1, 2023

- コマンドの出力を `` ```COMMAND | `` みたいに引用できるようにした。

v0.1.0
=======
Mar 23, 2023

- 初版