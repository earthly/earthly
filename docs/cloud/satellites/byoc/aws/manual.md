# Configuring BYOC Manually

This page documents the requirements and steps required to install [BYOC Satellites](../byoc.md) in AWS by manually configuring them.

## Requirements

Manual installation only requires that you meet the [base requirements](./requirements.md) for installation within AWS. No additional permissions within AWS are required.

## Manual Installation

Manual Installation requires you to provide all the information that Earthly would otherwise gather automatically. To do this, run:

```shell
earthly cloud install \
  --name <name> \
  --aws-account-id <aws-account-id> \
  --aws-region <aws-region> \
  --aws-security-group-id <aws-security-group-id> \
  --aws-ssh-key-name <aws-ssh-key-id> \
  --aws-subnet-id <aws-subnet-id> \
  --aws-instance-profile-arn <aws-instance-profile-arn> \
  --aws-earthly-access-role-arn <aws-earthly-access-role-arn>
```
### Parameters

| Name                                         | Description                                                                                                                                                   |
|----------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------|
| name                                         | The name to use to identify the cloud installation. This can be any string.                                                                                   |
| aws-account-id                               | The AWS account ID the required BYOC infrastructure was provisioned into.                                                                                     |
| aws-region                                   | The AWS region the required BYOC infrastructure was provisioned into. Currently, only `us-west-2` is supported.                                               |
| [aws-security-group-id](#security-group)     | The ID (`sg-0123456789abcde01`) of the security group used for new satellites.                                                                                |
| [aws-ssh-key-name](#ssh-key)                 | The name of the SSH key to be included in new satellites.                                                                                                     |
| [aws-subnet-id](#subnet)                     | The ID (`subnet-0123456789abcde01`) of the subnet that satellites will be launched into.                                                                      |
| [aws-instance-profile-arn](#instance-role)   | The ARN (`arn:aws:iam::012345678901:instance-profile/earthly/satellites/name/profile-name`) of the instance profile satellite instances will use for logging. |
| [aws-earthly-access-role-arn](#compute-role) | The ARN (`arn:aws:iam::012345678901:role/earthly/satellites/name/role-name`) of the role Earthly will assume to orchestrate satellites on your behalf.        |

## Manually Configuring BYOC Infrastructure

If you can't provision BYOC infrastructure using CloudFormation or Terraform, this should give you enough information to recreate what they do yourself. If you need help, you can [contact us](https://earthly.dev/slack).

### Subnet

You will need to create a Subnet (and a VPC, if needed) within the desired AWS account. The Subnet should have internet access, and have a CIDR block or DNS that is resolvable from within your network (VPN or otherwise).

### Security Group

Each satellite has one security group associated with it. Each satellite gets the following ingress rules by default:

| Protocol | CIDR             | From Port | To Port | Description                                                                                        |
|----------|------------------|-----------|---------|----------------------------------------------------------------------------------------------------|
| TCP      | Satellite Subnet | 22        | 22      | Allow SSH access from within the ingress subnet. Used for debugging satellite issues.              |
| TCP      | Satellite Subnet | 8372      | 8372    | Allow Buildkit access.                                                                             |
| TCP      | Satellite Subnet | 9000      | 9000    | Allow Prometheus scraping for monitoring your satellites. Metrics are exported by `node_exporter`. |

Satellite egress defaults to allowing general, unrestricted outbound traffic to the general internet.

### SSH Key

Any SSH key will do. Follow [AWS's guide](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/create-key-pairs.html) to create or upload an existing keypair.

### Cloudwatch Logs

Create a new Cloudwatch log group named `/earthly/satellites/<cloud-name>`, where `<cloud-name>` is the same value you will provide to Earthly via the `--name` parameter. The default class is `STANDARD`. [For more information on log group classes, see AWS documentation](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/CloudWatch_Logs_Log_Classes.html).

### Instance Role

Each satellite is configured to put relevant Buildkit logs in Cloudwatch. Earthly relies on an instance role to provide the relevant permissions.

Start by creating a new IAM policy:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "logs:PutLogEvents",
                "logs:CreateLogStream"
            ],
            "Resource": [
                "arn:aws:logs:us-west-2:012345678901:log-group:/earthly/satellites/name:log-stream:*",
                "arn:aws:logs:us-west-2:012345678901:log-group:/earthly/satellites/name"
            ],
            "Effect": "Allow"
        }
    ]
}
```

Follow the [AWS guide to create a new instance role](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-roles-for-amazon-ec2.html#create-iam-role), associating this new policy with your role. 

### Compute Role

Earthly uses a role within your AWS account to enable orchestration. This means that Earthly will only time-limited, user-revocable access to your cloud account.

Start by creating the policy needed:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "tag:GetResources",
                "iam:PassRole",
                "ec2:DescribeSubnets",
                "ec2:DescribeInstances",
                "ec2:DescribeInstanceTypes",
                "ec2:DescribeImages"
            ],
            "Resource": "*",
            "Effect": "Allow"
        },
        {
            "Action": [
                "ec2:RunInstances",
                "ec2:ModifyInstanceAttribute"
            ],
            "Resource": [
                "arn:aws:ec2:us-west-2::image/*",
                "arn:aws:ec2:us-west-2:012345678901:volume/*",
                "arn:aws:ec2:us-west-2:012345678901:security-group/sg-012345678901",
                "arn:aws:ec2:us-west-2:012345678901:network-interface/*",
                "arn:aws:ec2:us-west-2:012345678901:key-pair/name-satellite-key",
                "arn:aws:ec2:us-west-2:012345678901:instance/*",
                "arn:aws:ec2:us-west-2:012345678901:subnet/subnet-012345678901"
            ],
            "Effect": "Allow"
        },
        {
            "Action": [
                "ec2:TerminateInstances",
                "ec2:StopInstances",
                "ec2:StartInstances"
            ],
            "Resource": "arn:aws:ec2:us-west-2:012345678901:instance/*",
            "Effect": "Allow"
        },
        {
            "Action": [
                "ec2:DetachVolume",
                "ec2:DeleteVolume",
                "ec2:CreateTags",
                "ec2:AttachVolume"
            ],
            "Resource": [
                "arn:aws:ec2:us-west-2:012345678901:volume/*",
                "arn:aws:ec2:us-west-2:012345678901:instance/*"
            ],
            "Effect": "Allow"
        }
    ]
}
```

Follow [AWS's guide for creating a new IAM role with a custom trust policy](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_create_for-custom.html). Disregard all optional steps. Use the following trust policy:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::404851345508:role/compute-production"
            },
            "Action": "sts:AssumeRole"
        }
    ]
}
```
This trust policy enables Earthly to use the newly created role to orchestrate satellites on your behalf.
