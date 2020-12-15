# Earthly example

As a distinct example of a complete build, you can take a look at Earthly's own build. Earthly builds itself, and the build files are available on GitHub:

* [Earthfile](https://github.com/earthly/earthly/blob/main/Earthfile) - the root build file
* [buildkitd/Earthfile](https://github.com/earthly/earthly/blob/main/buildkitd/Earthfile) - the build of the buildkit daemon
* [earthfile2llb/parser/Earthfile](https://github.com/earthly/earthly/blob/main/earthfile2llb/parser/Earthfile) - the build of the parser, which generates .go files
* [examples/tests/Earthfile](https://github.com/earthly/earthly/blob/main/examples/tests/Earthfile) - system and smoke tests
* [contrib/earthfile-syntax-highlighting/Earthfile](https://github.com/earthly/earthly/blob/main/contrib/earthfile-syntax-highlighting/Earthfile) - the build of the VS Code extension

To invoke Earthly's build, check out the code and then run the following in the root of the repository

```bash
earthly +all
```

[![asciicast](https://asciinema.org/a/313845.svg)](https://asciinema.org/a/313845)
