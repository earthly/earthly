# earthfile-syntax-highlighting-vim README

<div align="center"><img alt="Earthly" width="700px" src="https://github.com/earthly/earthly/raw/main/img/logo-banner-white-bg.png" /></div>

Syntax highlighting for [Earthly](https://earthly.dev) Earthfiles for Vim.

For an introduction of Earthly see the [Earthly GitHub repository](https://github.com/earthly/earthly) or the [Earthly documentation](https://docs.earthly.dev).

## Installation Notes

### Automatic

The easiest way to install syntax highlighting for Earthfile is by running `make` in this directory. This craetes and copies all the files required.

### Manual

To install manually, copy `Earthfile.vim` to `~/.vim/syntax/Earthfile.vim`, you may need to create the directories.

Now write the following into the file at `~/.vim/ftdetect/Earthfile.vim`

```vim
au BufRead,BufNewFile Earthfile set filetype=Earthfile
au BufRead,BufNewFile build.earth set filetype=Earthfile
```

#### Neovim

Neovim users will have to change the `~/.vim/` prefix in the above steps to `~/.config/nvim`.

## Screenshot

![Java example Earthfile in Vim](https://raw.githubusercontent.com/vishnugt/earthly/main/contrib/earthfile-syntax-highlighting-vim/Screenshot.png)
