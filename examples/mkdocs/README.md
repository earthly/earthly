# Material for MkDocs

This example shows demonstrates [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/) using Earthly, and includes a live reloading development environment with `mkdocs serve`.

## Usage

Run `earthly +dev` to start a live-reloading development environment. A web browser should open automatically. Test this by making a change to `index.md` in `docs` -- the browser should reload on save. Create a production build with `earthly +build`.
