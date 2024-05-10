# OpenID Connect (OIDC) Authentication

Earthly can support cases where you might require access to a 3rd-party cloud provider as part of your build, without storing secrets in your CI or accessing credentials from your local environment.
This is especially useful in CI where otherwise, authentication requires MFA(multi-factor authentication).  
The OIDC protocol allows you to access the provider without storing credentials in your local environment or CI.

## Introduction

This page covers how to set up OIDC with cloud providers. 
At the moment the only AWS is supported.

## Cloud Providers

### AWS

#### Setup

1. Add the Earthly OIDC provider to AWS IAM - see the [AWS guide](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_create_oidc.html).
   1. Set https://api.earthly.dev as the provider URL.
   2. Set `sts.amazonaws.com` as the audience. 
2. Create a new IAM role (or configure an existing role you'd like to reuse) - see the [AWS guide](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_create_for-idp_oidc.html).
   1. Make sure to limit the permissions for the role (these are the actions the user can perform after assuming the role)
   2. Make sure to limit who can assume the role by specifying a trust policy such as:
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "<oidc-provider-name>"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "api.earthly.dev:aud": "sts.amazonaws.com",
          "api.earthly.dev:sub": "<earthly-org>/<earthly-project>"
        }
      }
    }
  ]
}
```

where:
* `<oidc-provider-name>` is the oidc provider's arn that was configured in step 1.
* `<earthly-org>` the earthly org the user is a member of and is set in the Earthfile or as part of the earthly build execution (see more details below).
* `<earthly-project>` the earthly project the user has access to [read secrets](./managing-permissions.md#earthly-project-access-levels) from, and is set in the Earthfile or as part of the earthly build execution (see more details below).

Note, a trust policy allows configuring different rules which you can mix and match to allow/disallow assuming the role by members of your team:
* To allow access to all members of the org:
```json
"Condition": {
    "StringLike": {
      "api.earthly.dev:sub": "<earthly-org>/*"
    }
}
```
* To allow access only to a specific user:
```json
"Condition": {
    "StringEquals": {
      "api.earthly.dev:email": "<user-email>"
    }
}
```
where `<user-email>` is the email address associated with the earthly account.

#### Usage

Once OIDC is configured, you can access AWS resources from your build.
Here is an example Earthfile to list S3 objects:
```dockerfile
VERSION --run-with-aws --run-with-aws-oidc 0.8

PROJECT <your-org>/<your-project>

aws:
    FROM amazon/aws-cli
    LET OIDC="role-arn=arn:aws:iam::1234567890:role/your-oidc-role,session-name=my-session,region=us-east-1"
    RUN --aws --oidc=$OIDC aws s3 ls
```

For more information on the `RUN --aws --oidc` flags, see [here](../earthfile/earthfile.md#--oidc-oidc-spec-experimental) 
