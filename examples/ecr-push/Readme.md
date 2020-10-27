# Push an image to ECR with Earthly

We are still making this easier, we hope it gets much shorter!

The steps to build your image go in the `+some-thing` target. The `+ecr-push` target will push it to ECS. Assuming you keep your AWS credentials in the standard environment variables, you can build your image and push it to ECR like this:

```
earth -P \
    --secret AWS_ACCESS_KEY_ID \
    --secret AWS_SECRET_ACCESS_KEY \
    --secret AWS_SESSION_TOKEN \
    --build-arg AWS_DEFAULT_REGION \
    --build-arg AWS_ACCOUNT_ID=<your-id> \
     +ecr-push

```

You can omit the session token if you are not using `assume-role`.

If you do not have the environment variables, you can set the AWS credentials like this via the already existing AWS CLI setup:

```
...
    --secret AWS_ACCESS_KEY_ID=$(aws configure get default.aws_access_key_id)
    --secret AWS_SECRET_ACCESS_KEY=$(aws configure get default.aws_secret_access_key)
    --secret AWS_SESSION_TOKEN=$(aws configure get default.aws_session_token)
...
```