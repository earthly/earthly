# Bring Your Own Cloud

"Bring Your Own Cloud" (BYOC) satellites **Experimental** are a hybrid between [self hosted](self-hosted.md) and [earthly-managed](../satellites.md) satellites. These satellites are managed by Earthly; just like our managed offering, but within your infrastructure. This blends the ease-of-use of managed satellites with the security benefits that self hosting can bring.

Here is how Earthly-managed satellites compare to BYOC: 

| Similarities                                                                               | Differences                                                           |
|--------------------------------------------------------------------------------------------|-----------------------------------------------------------------------|
| ✅ The ability to use `earthly satellite` commands to provision and decommission satellites | ❌ Networking must be user-provided                                    |
| ✅ Earthly-provided monitoring and reliability                                              | ❌ Earthly cannot access the machines for debugging or troubleshooting |
| ✅ Automatic updates                                                                        |                                                                       |
| ✅ Automatic sleep/wake                                                                     |                                                                       |

Here is how Self-hosted satellites compare to BYOC:

| Similarities                                                     | Differences                                                                            |
|------------------------------------------------------------------|----------------------------------------------------------------------------------------|
| ✅ The satellite lives in your infrastructure, next to your tools | ❌ You cannot manually provision satellites                                             |
| ✅ You pay for the usage the satellite incurs                     | ❌ You must allow Earthly access to at least some portion of an AWS account you control |
|                                                                  | ❌ You pay the cloud provider, not Earthly, for the compute usage                       |

Right now, BYOC is only supported in AWS.

## Getting Started

### Requirements
Before you begin to provision your BYOC configurations, there are a few requirements you need to make sure you meet first:

* Within AWS:
  * You need to have permissions to create a new CloudFormation stack in AWS, and install our provided template (see link below).
  * You need to have permissions to describe an existing CloudFormation stack in AWS.
  * A VPC, and *single* subnet that Earthly will place its satellites into. *Take note of the CIDR block, you will need it later.*
  * Any needed networking is ready - this includes things like NAT Gateways for internet access, or access to other resources with in the VPC.
* Within your VPN:
  * Allow ingress and egress between the satellites and client machines on your VPN.
  * Allow DNS resolution of the internal AWS domain names.
* On your machine:
  * Earthly must be installed.
  * You must have [AWS Credentials configured](https://docs.aws.amazon.com/cli/v1/userguide/cli-configure-files.html) properly.

### 1. Install AWS Components

Begin by installing our CloudFormation Template:
region=us-west-2#/stacks/create
[<img src="img/cloudformation.png" alt="Launch Stack In us-west-2" title="Launch CloudFormation Stack Quick-Create link" />](https://console.aws.amazon.com/cloudformation/home?region=us-west-2#/stacks/create/review?templateURL=https://production-byoc-installation.s3.us-west-2.amazonaws.com/cloudformation-byoc-installation.yaml)

If you need help installing a Cloudformation Template, you can reference [this guide from AWS](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/GettingStarted.Walkthrough.html). Once the installation is complete, continue to step 2.

Some important things to note:
* To ensure consistency across systems, Earthly will use the stack name as the name of the cloud in your CLI.
* Right now, only the `us-west-2` region is supported.
* All fields marked as "Required" *must* be filled in; even if you are electing to create new resources. This is due to CloudFormation template restrictions.

### 2. Add To Earthly

After installing our template; you must complete the installation by telling Earthly about your new configuration. Run:

```shell
earthly cloud install <stack-name>
```

Where `<stack-name>` is the name of the stack you provided to AWS in step 1. Earthly will automatically gather all the details it needs, validate them, and create the new installation.


### 3. Test

Now that you have your cloud installation configured in AWS and Earthly, its time to take it for a test drive! First, lets make sure we have all our clouds available, by running `earthly cloud list`. Your output should look something like this:

```shell
❯ earthly cloud list
   NAME           SATELLITES  STATUS          
   my-new-cloud   0           Green  
*  earthly-cloud  2           Green  
```

To use the new cloud, you can run `earthly cloud use my-new-cloud`. This will change the default cloud for the current organization. You can change back to use Earthly-managed satellites at any time by using the `earthly-cloud` installation.

Finally, launch a satellite by `earthly satellite launch my-satellite`. This should launch the new satellite in your new cloud!

## VPN Guides

### Tailscale

Tailscale is a super-simple VPN that is easy to set up, and works fairly well with BYOC satellites. However, there are a few requirements:

* Use a single [subnet router](https://tailscale.com/kb/1021/install-aws) per subnet.
  * If you have multiple cloud installations sharing a single subnet, the single subnet router will suffice.
* If you are running Earthly from within a Kubernetes pod, or GHA runner; you may need to make use of the userspace networking mode.
  * When using userspace networking, you need to add a Global nameserver to your DNS settings.
* It is required to add a [Split DNS](https://tailscale.com/learn/why-split-dns) entry for the `us-west-2.compute.internal` TLD, to point to all DNS addresses in your connected VPCs. This is usually the `x.x.0.2` address within a VPC.


