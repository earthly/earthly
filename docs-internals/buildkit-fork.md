# Why do we have a buildkit fork?

Here's a very rough list on some of the features we have in our BuildKit fork, which are too-specific to earthly,
and would not be a good fit to submit upstream.

- ability to pass arbitrary sockets from the host to buildkitd container; used by the interactive debugger
- to support `WITH DOCKER` via host-bind mounts
- faster exports of images from remote buildkit instance to local host; used for "pull ping" registry
- git LFS support
- verbose debugging of which files are sent to buildkit (via fsutils fork)
- `Export` method on gateway client, which is used to enable WAIT / END blocks
- healthcheck overrides, which are required for maintaining stability of our satellite instances
- modifications to the `llbsolver/ops/exec.go` in order to support `LOCALLY` mode
