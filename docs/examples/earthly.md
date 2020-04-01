# Earthly example

As a distinct example of a complete build, you can take a look at Earthly's own build. Earthly builds itself, and the build files are available on GitHub:

* [build.earth](https://github.com/vladaionescu/earthly/blob/master/build.earth) - the root build file
* [buildkitd/build.earth](https://github.com/vladaionescu/earthly/blob/master/buildkitd/build.earth) - the build of the buildkit daemon
* [earthfile2llb/parser/build.earth](https://github.com/vladaionescu/earthly/blob/master/earthfile2llb/parser/build.earth) - the build of the parser, which generates .go files
* [examples/tests/build.earth](https://github.com/vladaionescu/earthly/blob/master/examples/tests/build.earth) - system and smoke tests
* [contrib/earthfile-syntax-highlighting/build.earth](https://github.com/vladaionescu/earthly/blob/master/contrib/earthfile-syntax-highlighting/build.earth) - the build of the VS Code extension

To invoke Earthly's build, check out the code and then run the following in the root of the repository

```bash
earth +all
```

[![asciicast](https://asciinema.org/a/313845.svg)](https://asciinema.org/a/313845)
