---
name: Feature proposal
about: Submit an idea for a new feature
title: ''
labels: type:proposal
assignees: ''

---

**Use case**
<!-- Please describe the reason that you need or want this change -->

**Expected Behavior**
<!-- Please describe how you expect the change to work. Ideally, follow the style of reproduction steps - the steps that you expect to take, followed by the expected outcome.

Example: if the proposal is to "allow functions to return values", the steps may be:

1. Given the following Earthfile:
```
VERSION 0.8

FOO:
    FUNCTION
    RETURN "bar"

foo:
    ARG foo = DO +FOO
    RUN echo $foo
```
2. Run `earthly +foo`.
3. Confirm that the output echoes "bar".
-->
