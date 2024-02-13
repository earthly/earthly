# Why do we have a buildkit fork?

Here's a very rough list on some of the features we have in our BuildKit fork, which are too-specific to earthly,
and would not be a good fit to submit upstream.

- ability to pass arbitrary sockets from the host to buildkitd container; used by the interactive debugger
- to support `WITH DOCKER` via host-bind mounts
- faster exports of images and artifacts from remote buildkit instance to local host. This includes the Earthly exporter, the pull-ping call, the embedded registry, and the storage driver to plug the registry into the buildkit cache
- git LFS support
- verbose debugging of which files are sent to buildkit (via fsutils fork)
- `Export` method on gateway client, which is used to enable WAIT / END blocks
- healthcheck overrides, which are required for maintaining stability of our satellite instances
- modifications to the `llbsolver/ops/exec.go` in order to support `LOCALLY` mode
- GCAnalytics for stats on garbage collection and disk space used
- info call that includes the op load and the number of sessions currently executing
- StopIfIdle - the ability to shut down buildkit while ensuring that no session is active and that no session is accepted while closing
- session management capabilities, including the ability to sessions that run for too long
