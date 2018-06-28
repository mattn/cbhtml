# cbhtml

Copy HTML into Clipboard on Windows

## Usage

```
cat README.md | markdown | cbhtml
```

## Installation

```
go get github.com/mattn/cbhtml
```

## Misc

Copy text with syntax on Vim

```vim
function! s:clipboard_html(...) abort
  '<,'>TOhtml
  %!cbhtml
  bw!
endfunction

xnoremap <silent> <plug>(clipboard-html) :<c-u>call <sid>clipboard_html()<cr>
xmap C <plug>(clipboard-html)
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
