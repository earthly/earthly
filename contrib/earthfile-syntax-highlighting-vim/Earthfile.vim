" Vim syntax file
" Language: Earthfile
" Maintainer: Thomas Hobson <git@hexf.me>
" Latest Revision: 04 October 2020
" Source: https://docs.earthly.dev/earthfile

if exists("b:current_syntax")
  finish
endif

" Comments
" # <comment> (EOL)
syn region earthfileComment start="#" end="\n"

" Escapes
" \n
syn match earthfileEscape '\\.'
" \(EOL)
syn match earthfileEscape '\\$'

" Strings
" "<string>"
syn region earthfileString start="\"" end = "\""
" '<string>'
syn region earthfileString start="'" end = "'"

" Variables
" $<varname>
syn match earthfileVariable '\\$\(\w\-\)\+'
" ${<varname>}
syn region earthfileVariable start="${" end = "}"

" Operators
" && >> << | ; > <
syn match earthfileOperatorShell '&&\|>>\|<<\|;\|>\||'
" =
syn match earthfileOperatorAssign '='
" --...
syn match earthfileOperatorFlag '\s\-\+\(\w\|\-\)\+'

" Target
" debian:
syn match earthfileTargetLabel '^\zs\s*\w*\ze\:'
syn match earthfileTargetReference '\(\w\|_\|\-\|/\|:\|+\|\.\)*\s' contained nextgroup=earthfileKeyword

" Keywords
syn match earthfileKeyword '^\s*FROM DOCKERFILE\s*\|^\s*COPY\s*\|^\s*SAVE ARTIFACT\s*\|^\s*SAVE IMAGE\s*\|^\s*RUN\s*\|^\s*LABEL\s*\|^\s*EXPOSE\s*\|^\s*VOLUME\s*\|^\s*USER\s*\|^\s*ENV\s*\|^\s*ARG\s*\|^\s*BUILD\s*\|^\s*WORKDIR\s*\|^\s*ENTRYPOINT\s*\|^\s*CMD\s*\|^\s*GIT CLONE\s*\|^\s*DOCKER LOAD\s*\|^\s*DOCKER PULL\s*\|^\s*HEALTHCHECK\s*NONE\|^\s*HEALTHCHECK\s*CMD\|^\s*WITH DOCKER\|^\s*END'
syn match earthfileKeyword '^\s*FROM\s*' nextgroup=earthfileBaseImage
syn match earthfileBaseImage '\S\+' contained


syn match earthfileKeyword '^\s*SAVE ARTIFACT\s*' nextgroup=earthfileTargetReference

syn match earthfileKeyword '\s*AS LOCAL' contained

" Highlights
hi def link earthfileKeyword Keyword

hi def link earthfileOperatorShell Operator
hi def link earthfileOperatorAssign Operator
hi def link earthfileOperatorFlag Special

hi def link earthfileBaseImage Constant

hi def link earthfileTargetLabel Statement
hi def link earthfileTargetReference Constant

hi def link earthfileComment Comment
hi def link earthfileEscape SpecialChar
hi def link earthfileString String
hi def link earthfileVariable Identifier

let b:current_syntax = "earthfile"