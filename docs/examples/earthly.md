# Earthly example

As a distinct example of a complete build, you can take a look at Earthly's own build. Earthly builds itself, and the build files are available on GitHub:

* [Earthfile](https://github.com/earthly/earthly/blob/master/Earthfile) - the root build file
* [buildkitd/Earthfile](https://github.com/earthly/earthly/blob/master/buildkitd/Earthfile) - the build of the buildkit daemon
* [earthfile2llb/parser/Earthfile](https://github.com/earthly/earthly/blob/master/earthfile2llb/parser/Earthfile) - the build of the parser, which generates .go files
* [examples/tests/Earthfile](https://github.com/earthly/earthly/blob/master/examples/tests/Earthfile) - system and smoke tests
* [contrib/earthfile-syntax-highlighting/Earthfile](https://github.com/earthly/earthly/blob/master/contrib/earthfile-syntax-highlighting/Earthfile) - the build of the VS Code extension

To invoke Earthly's build, check out the code and then run the following in the root of the repository

```bash
earth +all
```

[![asciicast](https://asciinema.org/a/313845.svg)](https://asciinema.org/a/313845)
