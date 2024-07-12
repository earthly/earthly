# Installing BYOC in AWS Using Terraform

This page documents the requirements and steps required to install [BYOC Satellites](../byoc.md) in AWS using Earthly's Terraform module.

## Requirements
Before you begin to provision your BYOC configuration, ensure that you meet the [base requirements](./requirements.md) for installation within AWS.

There are also a few additional requirements you will need to make sure you meet:

* Terraform is installed, and on your `$PATH`.
* An AWS account or role that can create all the resources specified in the module.
* You have permission to list Terraform outputs.

## Installation

You can find our module in the [public Terraform Registry](https://registry.terraform.io/modules/earthly/byoc/aws/latest). If you're curious about what we're provisioning, you can look at our [source code](https://github.com/earthly/terraform-aws-byoc/blob/main/main.tf).

### Quickstart

Place the following code into a file named `byoc.tf`, in an empty directory:

```hcl
module "byoc" {
  source  = "earthly/byoc/aws"
  version = "0.0.8"

  cloud_name = "my-cloud"
  subnet = "subnet-0123456789abcde01"
}

output "my-cloud" {
  values = module.byoc.automatic_installation
}
```

Open your terminal, and navigate to the directory with `byoc.tf` in it. Run `terraform init && terraform apply`, and inspect the resources it wants to create. If they appear ok, type `yes` to create them. 

After Terraform finishes running, you can link this freshly provisioned infrastructure to an Earthly cloud by running:

```shell
earthly cloud install --via terraform --name my-cloud
```

Assuming the installation reports the status as `Green`, you should be good to go!

### Colocating With Other Terraform Code

You can use the module as explained in [Quickstart](#quickstart). However, if you would like to also enable automatic installation, some additional conditions apply:

1. The installation command is run in the same directory as the module containing your BYOC block.
2. The `automatic_installation` output is exported. The name of the output is the value used for the `--name` parameter.

We recommend that the name of the output matches the name in the `cloud_name` of the module, to ensure that naming is consistent between AWS, your tooling, and Earthly.

## Module Parameters

### Module Inputs

| Name           | Description                                                                                                                                              |
|----------------|----------------------------------------------------------------------------------------------------------------------------------------------------------|
| cloud_name     | The name to use to identify the cloud installation. Used by Earthly during automatic installation, and to mark related resources in AWS.                 |
| subnet         | The subnet Earthly will deploy satellites into.                                                                                                          |
| ssh_public_key | (Optional) The SSH key to include in provisioned satellites. If left unspecified, a new key is generated, and the private key is available as an output. |


### Module Outputs

| Name                 | Description                                                                                                                              |
|----------------------|------------------------------------------------------------------------------------------------------------------------------------------|
| installation_name    | The name to use to identify the cloud installation. Used by Earthly during automatic installation, and to mark related resources in AWS. |
| security_group_id    | The ID of the security group for new satellites.                                                                                         |
| ssh_key_name         | The name of the SSH key in AWS that is included in new satellites.                                                                       |
| ssh_private_key      | (Sensitive) The private key, if `ssh_public_key` is unspecified.                                                                         |
| instance_profile_arn | The ARN of the instance profile satellite instances will use for logging.                                                                |
| compute_role_arn     | The ARN of the role Earthly will assume to orchestrate satellites on your behalf.                                                        |
