# Pushing images to AWS ECR

## Introduction

Amazons' Elastic Container Registry (ECR) is a hosted docker repository that requires extra configuration for day-to-day use. This configuration is not typical of other repositories, and there are some considerations to account for when using it with Earthly. This guide will walk you through creating an Earthfile, building an image, and pushing it to ECR.

## Create an Earthfile

No special considerations are needed in the Earthfile itself. You can use `SAVE IMAGE` just like any other repository.

```
FROM alpine:3.11

build:
    RUN echo "Hello from Earthly!" > motd
    ENTRYPOINT cat motd
    SAVE IMAGE --push <aws_account_id>.dkr.ecr.<region>.amazonaws.com/hello-earthly:with-love
```

## Install and Configure the ECR Credential Helper

ECR does not issue permanent credentials. Instead, it relies on your AWS credentials to issue docker credentials. You can follow instructions [here](https://docs.aws.amazon.com/cli/latest/reference/ecr/get-login.html) to log in with generated credentials, but the process will need to be repeated every 12 hours. In practice, this often means lots of glue code in your CI pipeline to keep credentials up to date.

AWS has released a credential helper to ease logging into ECR. It may be that you already have the credential helper installed, since it has been included with Docker Desktop as of version [2.4.0.0](https://docs.docker.com/docker-for-windows/release-notes/#docker-desktop-community-2400). If not, you can follow installation instructions on their GitHub repository here. Here is a sample `.docker/config.json` to enable the usage of this helper:

```
{
        "credHelpers": {
                "<aws_account_id>.dkr.ecr.<region>.amazonaws.com": "ecr-login"
        }
}

```

## IAM

Ensure that you have correct permissions to push the images. The ECR helper is aware of the `AWS_PROFILE` variable; and can work under an assumed role. Here is a minimum set of privileges needed to push to ECR from Earthly:

```
{
    "Version": "2008-10-17",
    "Statement": [
        {
            "Sid": "AllowPushPull",
            "Effect": "Allow",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::<aws_account_id>:user/push-pull-user",
                ]
            },
            "Action": [
                "ecr:GetAuthorizationToken",
                "ecr:GetDownloadUrlForLayer",
                "ecr:BatchGetImage",
                "ecr:BatchCheckLayerAvailability",
                "ecr:PutImage",
                "ecr:InitiateLayerUpload",
                "ecr:UploadLayerPart",
                "ecr:CompleteLayerUpload"
            ]
        }
    ]
}
```

Additional examples for policy configuration can be found [here](https://docs.aws.amazon.com/AmazonECR/latest/userguide/repository-policy-examples.html).

## Run the target

With the helper installed, no special To build and push an image, simply execute the build target. Don't forget the `--push` flag!

```
❯ earthly --push +build
           buildkitd | Found buildkit daemon as docker container (earthly-buildkitd)
         alpine:3.11 | --> Load metadata linux/amd64
               +base | --> FROM alpine:3.11
               +base | [██████████] resolve docker.io/library/alpine:3.11@sha256:0bd0e9e03a022c3b0226667621da84fc9bf562a9056130424b5bfbd8bcb0397f ... 100%
              +build | --> RUN echo "Hello from Earthly!" > motd
              output | --> exporting outputs
              output | [██████████] exporting layers ... 100%
              output | [██████████] exporting manifest sha256:9ab4df74dafa2a71d71e39e1af133d110186698c78554ab000159cfa92081de4 ... 100%
              output | [██████████] exporting config sha256:6feef98708c14c000a6489a2a99315a5328c2c16091851ae10438b53f655d042 ... 100%
              output | [██████████] pushing layers ... 100%
              output | [██████████] pushing manifest for <aws_account_id>.dkr.ecr.<region>.amazonaws.com/hello-earthly:with-love ... 100%
              output | [██████████] sending tarballs ... 100%
=========================== SUCCESS ===========================
Loaded image: <aws_account_id>.dkr.ecr.<region>.amazonaws.com/hello-earthly:with-love
              +build | Image +build as <aws_account_id>.dkr.ecr.<region>.amazonaws.com/hello-earthly:with-love (pushed)

```

If you get a message saying `basic credentials not found`; your distribution may not have the most recent version installed. A simple workaround is to simply prepend `AWS_SDK_LOAD_CONFIG=true` to your Earthly invocation. This will force the helper to use the SDK over built-in config when executing. You can track this issue [here](https://github.com/awslabs/amazon-ecr-credential-helper/issues/232).
