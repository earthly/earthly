---
name: Bug report
about: Report an unexpected problem
title: ''
labels: type:bug
assignees: ''

---

**What went wrong?**
<!-- Can you describe what happened? -->

<!-- Can you reproduce it? If so, please list the steps to reproduce the issue. -->

<!-- Do you have an Earthfile that you can share which showcases the problem? If so, please include it here. -->

**What should have happened?**
<!-- Can you describe the expected outcome? -->

<!-- Have you found any workarounds? Can you share them for any other users who might be experiencing the issue? -->

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
