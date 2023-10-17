# Community Contributions

Unfortunately outside contributor PRs can't run our tests due to our GHA never sharing our secrets policy   
(which means they can't access our mirror, or make use of the way we share binaries between stages of the tests).

As a workaround, we can create a new temporary PR off of the PR.  
First, review the PR change is satisfactory and ensure it is safe to run the tests for (no malicious attempt to steal secrets).  
Then apply the label [approved-for-tests](https://github.com/earthly/earthly/labels/approved-for-tests). This would trigger a GHA that would open a new draft PR, linked to the original.  
If and when all the tests pass, the original PR can then be approved and the temporary draft PR and branch can be closed and deleted respectively.
