
# Examples

In this section, you will find some examples of Earthfiles for:
* [Go](./go.md)
* [Java](./java.md)
* [JS](./js.md)
* [Mono-repos](./monorepo.md)
* [Multi-repos](./multirepo.md)

The code for all the examples is available in the [examples GitHub directory](https://github.com/vladaionescu/earthly/tree/master/examples).

Beyond the examples in this section, you can additionally take a look at Earthly's own build (Earthly builds itself), available on GitHub:
* [/build.earth](https://github.com/vladaionescu/earthly/blob/master/build.earth)
* [/buildkitd/build.earth](https://github.com/vladaionescu/earthly/blob/master/buildkitd/build.earth)
* [/earthfile2llb/parser/build.earth](https://github.com/vladaionescu/earthly/blob/master/earthfile2llb/parser/build.earth)
* [/examples/tests/build.earth](https://github.com/vladaionescu/earthly/blob/master/examples/tests/build.earth)
* [/contrib/earthfile-syntax-highlighting/build.earth](https://github.com/vladaionescu/earthly/blob/master/contrib/earthfile-syntax-highlighting/build.earth)

To invoke Earthly's build, check out the code and then run

```bash
earth +all
```
