# GitHub Actions

Documentation for this repo GitHub actions configuration

## Skipping PR Workflows (DISABLED! SEE NOTE BELOW)

The following is disabled due to issue https://github.com/orgs/community/discussions/13261.  
Once this issue is resolved, the configuration can be restored by reverting the changes in [this PR](https://github.com/earthly/earthly/pull/3345)

### Motivation
Some PR workflows (workflows triggered by PRs) in this repo might take substantial time to complete,  
and it is not always necessary to run them all, especially in cases where the affected files are documents or any other file that don't affect our workflows or tests.

### GitHub Builtin Support
GitHub supports defining filters on workflows so that they won't get triggered, and there is a builtin filter that would do so according to affected files - [paths-ignore](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#onpushpull_requestpull_request_targetpathspaths-ignore).  
The shortcoming of this filter is that when the workflow is marked as a `required check` and does not get triggered because of this filter,  
the workflow will be stuck in `pending` and won't allow the PR to be merged.  
Thus, using this filter is a good approach when handling workflows that are not required checks, but it's problematic otherwise.

### Solution for required workflows
An alternative to `paths-ignore` is examining the affected files at the `job` level instead.  
Jobs allow you to set conditions on whether they should execute or not.  
The advantage of using conditions on the job level is that if a job does not run due to a condition evaluation to false,
GitHub will still mark the workflow as passing/successful instead of pending,  
and will not block the PR even if the workflow is marked as a required check ([reference](https://docs.github.com/en/actions/using-jobs/using-conditions-to-control-job-execution#overview)).  

To examine the files, [dorny/paths-filter](https://github.com/dorny/paths-filter) is used.
Similarly to GitHub's `paths-ignore`, `dorny/paths-filter` supports setting rules that match against the affected files, which evaluates to a boolean value.
We can then use that boolean value to determine whether subsequent jobs in the workflow should execute or not. 

### Complexity due to reusable workflows

#### Background
A problem arises when the job we wish to conditionally run uses a reusable workflow.
The problem is that GitHub treats such configuration as a parent-child relationship,
and any job that is skipped, by default will also cause its children jobs to skip as well.
In addition, when we set `required checks`, we cannot properly select the parent job, only the child job.  
In this scenario, if a child job is skipped as a result of skipping its parent, the child job will get stuck in `pending` and the PR will be blocked from merging (just like in the original problem).

#### Final Solution
So, in order to work around the above-mentioned issue, we need to avoid
skipping the parent job. Instead, we evaluate the boolean as before, and pass it as  
an input argument to the child job. As a matter of standardization, we call this input arg `SKIP_JOB`.
This argument is then used to conditionally run the _child_ job instead of the _parent_ job,  
and so if the child job is skipped, even if it is marked as a required check, it won't get stuck in `pending`.

Lastly, because most jobs are dependent on (`needs` directive) `build-earthly` job, if the latter is skipped, by default - so will the dependent jobs.  
As described before, that is something we need to avoid. To change that default configuration, we change
the default condition of each dependent job to `$ {{ !failure() }}` which will allow the dependent jobs to run whether the dependency job was successful or if it was skipped.
