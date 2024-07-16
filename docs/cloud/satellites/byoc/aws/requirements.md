# Installing BYOC in AWS 

This page documents the requirements to install [BYOC Satellites](../byoc.md) in AWS. Additional requirements may apply depending on your chosen installation method.

## Networking
Because every configuration is different, using BYOC within your organization will require you to configure your networking to match your use case. You'll need to ensure that:

* Traffic from Earthly clients, including CI build runners and developer workstations, can reach the satellites directly.
* Traffic to any required resources (e.g. private repositories, the internet, etc) are allowed. 
* Internal AWS DNS names (or IP addresses) must resolve to an address reachable on the network.

These can all be accomplished with most VPN technologies. We recommend and have direct experience with [Tailscale](../vpn/tailscale.md). If you need help configuring other networking scenarios, [please reach out to us](https://earthly.dev/slack)!

## Cloud
Configuring BYOC in your AWS account requires:

* A VPC, 
* and a chosen subnet that Earthly will place its satellites into. *Take note of the CIDR block, you may need it later.*

Right now, BYOC is only supported in the `us-west-2` (Oregon) region.

## On Your Machine
You'll need to finish the installation at the command line. To do this, ensure that [Earthly](https://earthly.dev/get-earthly) is installed on your system.
