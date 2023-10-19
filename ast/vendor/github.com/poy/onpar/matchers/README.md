# OnPar Matchers

OnPar provides a set of minimalistic matchers to get you started.
However, the intention is to be able to write your own custom matchers so that
your code is more readable.


## Matchers List
- [String Matchers](#string-matchers)
- [Logical Matchers](#logical-matchers)
- [Error Matchers](#error-matchers)
- [Channel Matchers](#channel-matchers)
- [Collection Matchers](#collection-matchers)
- [Other Matchers](#other-matchers)


## String Matchers
### StartWith
StartWithMatcher accepts a string and succeeds if the actual string starts with
the expected string.

```go
Expect(t, "foobar").To(StartWith("foo"))
```

### EndWith
EndWithMatcher accepts a string and succeeds if the actual string ends with
the expected string.

```go
Expect(t, "foobar").To(EndWith("bar"))
```

### ContainSubstring
ContainSubstringMatcher accepts a string and succeeds if the expected string is a
sub-string of the actual.

```go
Expect(t, "foobar").To(ContainSubstring("ooba"))
```
### MatchRegexp

## Logical Matchers
### Not
NotMatcher accepts a matcher and will succeed if the specified matcher fails.

```go
Expect(t, false).To(Not(BeTrue()))
```

### BeAbove
BeAboveMatcher accepts a float64. It succeeds if the actual is greater
than the expected.

```go
Expect(t, 100).To(BeAbove(99))
```

### BeBelow
BeBelowMatcher accepts a float64. It succeeds if the actual is
less than the expected.

```go
Expect(t, 100).To(BeBelow(101))
```

### BeFalse
BeFalseMatcher will succeed if actual is false.

```go
Expect(t, 2 == 3).To(BeFalse())
```

### BeTrue
BeTrueMatcher will succeed if actual is true.

```go
Expect(t, 2 == 2).To(BeTrue())
```

### Equal
EqualMatcher performs a DeepEqual between the actual and expected.

```go
Expect(t, 42).To(Equal(42))
```

## Error Matchers
### HaveOccurred
HaveOccurredMatcher will succeed if the actual value is a non-nil error.

```go
Expect(t, err).To(HaveOccurred())

Expect(t, nil).To(Not(HaveOccurred()))
```

## Channel Matchers
### Receive
ReceiveMatcher only accepts a readable channel. It will error for anything else.
It will attempt to receive from the channel but will not block.
It fails if the channel is closed.

```go
c := make(chan bool, 1)
c <- true
Expect(t, c).To(Receive())
```

### BeClosed
BeClosedMatcher only accepts a readable channel. It will error for anything else.
It will succeed if the channel is closed.

```go
c := make(chan bool)
close(c)
Expect(t, c).To(BeClosed())
```

## Collection Matchers
### HaveCap
This matcher works on Slices, Arrays, Maps and Channels and will succeed if the
type has the specified capacity.

```go
Expect(t, []string{"foo", "bar"}).To(HaveCap(2))
```
### HaveKey
HaveKeyMatcher accepts map types and will succeed if the map contains the
specified key.

```go
Expect(t, fooMap).To(HaveKey("foo"))
```

### HaveLen
HaveLenMatcher accepts Strings, Slices, Arrays, Maps and Channels. It will
succeed if the type has the specified length.

```go
Expect(t, "12345").To(HaveLen(5))
```

## Other Matchers
### Always
AlwaysMatcher matches by polling the child matcher until it returns an error.
It will return an error the first time the child matcher returns an error.
If the child matcher never returns an error, then it will return a nil.

By default, the duration is 100ms with an interval of 10ms.

```go
isTrue := func() bool {
  return true
}
Expect(t, isTrue).To(Always(BeTrue()))
```

### Chain
### ViaPolling
