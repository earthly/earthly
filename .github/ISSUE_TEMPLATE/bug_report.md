---
name: Bug report
about: Report an unexpected problem
title: ''
labels: type:bug
assignees: ''

---

**Description**
<!-- A quick description of the bug -->

**Reproduction Steps**
<!-- An ordered list of steps to take to reproduce the problem.

Example:
1. Set up earthly to use podman
2. Run `earthly boostrap`
3. Run any earthly target
-->

**Expected behavior**
<!-- What you expected to happen.

Continuing the example from earlier: earthly runs my target in podman.
-->

**Actual behavior**
<!-- What actually happened.

Continuing the example from earlier: earthly fails to connect to buildkitd due to TLS certificate errors.
-->

**Other Helpful Information**
<!-- Please include any additional information you might have. This may include:

- The logs from the `earthly-buildkitd` container, usually from running `docker logs earthly-buildkitd`.
- The output of the failing command.
- Output of the failing command with `--verbose ` enabled, e.g. `earthly --verbose +some-failing-target`.
- Stack trace.
    - If you `kill -SIGABRT` your earthly process, you will get a full stack
      trace. This is useful if the process is stuck or slow, so that we can get an
      idea of which functions might be stuck.
- In some rare circumstances: screenshots (although if the output is text, we
  prefer copy/pasted text over screenshots).
-->
