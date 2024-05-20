# Bring Your Own Cloud

"Bring Your Own Cloud" (BYOC) satellites **Experimental** are a hybrid between [self hosted](../self-hosted.md) and [earthly-managed](../../satellites.md) satellites. These satellites are managed by Earthly; just like our managed offering, but within your infrastructure. This blends the ease-of-use of managed satellites with the security benefits that self hosting can bring.

Here is how Earthly-managed satellites compare to BYOC: 

| Similarities                                                                               | Differences                                                           |
|--------------------------------------------------------------------------------------------|-----------------------------------------------------------------------|
| ✅ The ability to use `earthly satellite` commands to provision and decommission satellites | ❌ It is up to you to get traffic to and from the satellite            |
| ✅ Earthly-provided monitoring and reliability                                              | ❌ Earthly cannot access the machines for debugging or troubleshooting |
| ✅ Automatic updates                                                                        |                                                                       |
| ✅ Automatic sleep/wake                                                                     |                                                                       |

Here is how Self-hosted satellites compare to BYOC:

| Similarities                                                     | Differences                                                                            |
|------------------------------------------------------------------|----------------------------------------------------------------------------------------|
| ✅ The satellite lives in your infrastructure, next to your tools | ❌ You cannot manually provision satellites                                             |
| ✅ You pay for the usage the satellite incurs                     | ❌ You must allow Earthly access to at least some portion of an AWS account you control |
|                                                                  | ❌ You pay the cloud provider, not Earthly, for the compute usage                       |

## Installation

### Step 1: Networking

Configure your networking (usually a VPN) to ensure access to the satellites that will be created by Earthly. We have guides to enable BYOC on VPNs for the following VPN providers:
* [Tailscale](./tailscale.md)

### Step 2: Configure Your Cloud Provider

Right now, BYOC Satellites are only supported in AWS. See our [AWS CloudFormation](./aws.md) instructions to provision the required AWS resources.


### Step 3: Install In Earthly

After installing the required resources within your cloud provider; you must complete the installation by telling Earthly about your new configuration. Run:

```shell
earthly cloud install <name>
```

Where `<name>` is the name of the Installation, as specified by your cloud provider. Assuming it reports the status as `Green`, you should be good to go!


### Step 4. Test Drive

Now that you have your cloud installation configured in your cloud provider and Earthly, its time to take it for a test drive!

First, make sure you can see the cloud you just installed by running `earthly cloud list`, which lists all the cloud installations within your organization. Your output should look something like this:

```shell
❯ earthly cloud list
   NAME           SATELLITES  STATUS          
   my-new-cloud   0           Green  
*  earthly-cloud  2           Green  
```
The `*` indicates the default cloud that will be used when launching satellites within your organization, unless otherwise specified.

Note that the `earthly-cloud` installation is a special cloud present in all organizations. Satellites within this cloud are managed within Earthly's cloud, by Earthly.  

To test-launch a new satellite within the cloud, run `earthly satellite launch --cloud <name> my-byoc-sat`. This will launch a new satellite using your cloud. Assuming that works, kick the tires by trying to run one of your builds on it!

If everything looks good, you can run `earthly cloud use <name>` to set this cloud to be the default for your organization. You can change back to use Earthly-managed satellites at any time by running the `earthly cloud use earthly-cloud` installation.

Note that this setting is global, for all users within the org. This prevents people from launching satellites in the wrong cloud, and accidentially disseminating information that shouldn't be.
