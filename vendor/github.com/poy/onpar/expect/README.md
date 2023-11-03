# expect

`Expect` is used to used to invoke matchers and give a helpful error message. `Expect` helps a test read well and minimize overhead and scaffolding.

```go
func TestSomething(t *testing.T) {
    x := "foo"
    Expect(t, x).To(Equal("bar"))
}
```

The previous example will fail and pass the error from the matcher `Equal` to `t.Fatal`.

```go
func TestSomething(t *testing.T) {
    x := "foo"
    Expect(t, x).To(StartsWith("foobar"))
}
```

This example will pass. 
