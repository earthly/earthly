# Installing BYOC in AWS Using Terraform

This page documents the requirements and steps required to install [BYOC Satellites](../byoc.md) in AWS using Earthly's Terraform module.

## Requirements
Before you begin to provision your BYOC configuration, ensure that you meet the [base requirements](./requirements.md) for installation within AWS.

There are also a few additional requirements you will need to make sure you meet:

* An AWS account or role that can create all the resources specified in the module.
* You have permission to list Terraform outputs.

## Using the BYOC Terraform Module
You can find the module in the [public Terraform Registry](https://registry.terraform.io/modules/earthly/byoc/aws/latest). Here is an example of using the module:

```hcl
module "byoc" {
  source  = "earthly/byoc/aws"
  version = "0.0.8"

  cloud_name = "my-cloud"
  subnet = "subnet-0123456789abcde01"
}
```
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

## Installation

Earthly is able to automatically install BYOC when provisioned via Terraform by running:

```shell
earthly cloud install --via terraform --name <output-name>
```

Automatic installation requires some conditions to be met. If both of these conditions cannot be met, you will need to use the [manual](./manual.md) installation method with the infrastructure provisioned by Terraform.

1. Terraform is in your `$PATH`.
2. The install command is run in the same directory as the module containing your BYOC block.
2. The `automatic_installation` output is exported under the name specified by `<output-name>`. You could enable automatic installation for the example earlier by adding the following output:

```hcl
output "my-cloud" {
    values = module.byoc.automatic_installation
}
```

 We recommend that `<output-name>` match the name in the `cloud_name` of the module, to ensure that naming is consistent between AWS, your tooling, and Earthly.

Assuming the installation reports the status as `Green`, you should be good to go!
