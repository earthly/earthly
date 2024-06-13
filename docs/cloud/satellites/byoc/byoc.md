# Bring Your Own Cloud

"Bring Your Own Cloud" (BYOC) satellites **Experimental** are a hybrid between [self hosted](../self-hosted.md) and [earthly-managed](../../satellites.md) satellites. These satellites are managed by Earthly; just like our managed offering, but within your infrastructure. This blends the ease-of-use of managed satellites with the security benefits that self hosting can bring.

BYOC satellites are [available with an Enterprise plan](https://earthly.dev/pricing).

### BYOC vs Cloud

|                                                                                    | Earthly Cloud                     | Earthly BYOC                               |
|------------------------------------------------------------------------------------|-----------------------------------|--------------------------------------------|
| Who is responsible for monitoring and reliability of the Satellite                 | ✅ Earthly                         | ✅ Earthly                                  |
| Satellites are deployed within your internal network                               | ❌ No                              | ✅ Yes                                      |
| Earthly Cloud and Earthly staff are prevented from accessing your internal network | ✅ N/A                             | ✅ Yes                                      |
| How is compute billed                                                              | ✅ Zero-margin compute via Earthly | ✅ Supported by you via your cloud provider |
| Automatic updates                                                                  | ✅ Yes                             | ✅ Yes                                      |
| Auto-sleep to drastically reduce compute cost                                      | ✅ Yes                             | ✅ Yes                                      |
| Automatic management and GCing of cache volumes                                    | ✅ Yes                             | ✅ Yes                                      |
| Users can launch and remove satellites via the `earthly sat` CLI                   | ✅ Yes                             | ✅ Yes                                      |
| Requires access to a set of limited AWS capabilities                               | ✅ No                              | ❌ Yes                                      |


### BYOC vs Self-Hosted

|                                                                                    | Earthly Self-Hosted | Earthly BYOC |
|------------------------------------------------------------------------------------|---------------------|--------------|
| Who is responsible for monitoring and reliability of the Satellite                 | ❌ You               | ✅ Earthly    |
| Satellites are deployed within your internal network                               | ✅ Yes               | ✅ Yes        |
| Earthly Cloud and Earthly staff are prevented from accessing your internal network | ✅ Yes               | ✅ Yes        |
| Compute is billed directly to you from your cloud provider                         | ✅ Yes               | ✅ Yes        |
| Automatic updates                                                                  | ❌ No                | ✅ Yes        |
| Auto-sleep to drastically reduce compute cost                                      | ❌ No                | ✅ Yes        |
| Automatic management and GCing of cache volumes                                    | ❌ No                | ✅ Yes        |
| Users can launch and remove satellites via the `earthly sat` CLI                   | ❌ No                | ✅ Yes        |
| Requires access to a set of limited AWS capabilities                               | ✅ No                | ❌ Yes        |
## Installation

### Step 1: Configure Your Cloud Provider

Follow the instructions on the [AWS CloudFormation](./aws.md) page to provision the required AWS resources. Right now, BYOC Satellites are only supported in AWS.


### Step 2: Install In Earthly

After installing the required resources within your cloud provider; you must complete the installation by telling Earthly about your new configuration. Run:

```shell
earthly cloud install <name>
```

Where `<name>` is the name of the Installation, as specified by your cloud provider. Assuming it reports the status as `Green`, you should be good to go!


### Step 3: Networking

To use a satellite created by BYOC, you'll need to configure your networking (usually a VPN) to ensure access. We have guides to enable BYOC on VPNs for the following VPN providers:
* [Tailscale](./tailscale.md)

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

{% hint style='info' %}
Note that the `earthly-cloud` installation is a special cloud present in all organizations. Satellites within this cloud are managed within Earthly's cloud, by Earthly. You can change back to use Earthly-managed satellites at any time by running the `earthly cloud use earthly-cloud` installation.
{% endhint %}

### Step 5. Launch A Satellite

To test-launch a new satellite within the cloud, run:

```shell
earthly satellite launch --cloud <name> my-byoc-sat
```
This will launch a new satellite using your newly created cloud. Assuming that works, kick the tires by trying to run one of your builds on it!


### Step 6. Use The Cloud

If everything looks good, you're done! If you would like this new cloud to be the default, simply run:

```shell
earthly cloud use <name>
```

This makes the new cloud you created be the default for _the entire organization_.

{% hint style='warning' %}
This setting is global for all users within the org. This prevents people from launching satellites in the wrong cloud, and accidentally disseminating information that shouldn't be.
{% endhint %}
