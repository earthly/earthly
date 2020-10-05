#!/usr/bin/env bash
set -e

function install {
  vimpath=$1
  mkdir -p "$vimpath/"{syntax,ftdetect}
  echo "installing vim highlighting to $vimpath"
  cp "Earthfile.vim" "$vimpath/syntax/Earthfile.vim"
  echo "au BufRead,BufNewFile Earthfile set filetype=Earthfile" > "$vimpath/ftdetect/Earthfile.vim"
  echo "au BufRead,BufNewFile build.earth set filetype=Earthfile" >> "$vimpath/ftdetect/Earthfile.vim"
}

declare -a vimpaths=("$HOME/.vim" "$HOME/.config/nvim")
for vimpath in "${vimpaths[@]}"; do
  if [[ -d "$vimpath" ]]; then
    install "$vimpath"
  fi
done
