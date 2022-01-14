These tests are for tests that cannot run via earthly-in-earthly.

If you add or remote any tests from this directory, the corresponding entry in `.github/workflows/ci.yml`
must also be updated by running:

    earthly +generate-github-tasks
