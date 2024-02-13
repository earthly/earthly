# Guide Template

## Overview

The headers in this document provide a skeleton to fill in when adding documentation for a new CI platform. This section should include an overview of how Earthly might fit into this CI. If the CI offers multiple modes, we should mention them even if they are not all documented.

### Compatibility

It should include any relevant compatibility information, such as versions of plugins or runtimes.

### Resources

It should also include links to the CI systems relevant documentation for setup and configuration.

## Setup

This section should contain all special tweaks needed that are different than the ones detailed in [overview](../overview.md). If there are different tweaks for different modes, then each mode should have its own Setup header.

### Dependencies

Any special dependencies needed for this CI, w.r.t. Earthly. Probably a rare section.

### Installation

Any special installation instructions needed for this CI, w.r.t. Earthly. Special configuration or example targets for configuring a container-based CI would go here.

### Configuration

Most common section. Any special, or recommended configuration (for Earthly or dependencies) for this CI.

## Additional Notes

Any extra notes regarding the CI that aren't per-setup-section.

## Example

{% hint style='danger' %}
##### Note

This example is not production ready, and is intended to showcase configuration needed to get Earthly off the ground. If you run into any issues, or need help, [don't hesitate to reach out](https://github.com/earthly/earthly/issues/new)!

{% endhint %}

This should walk through (and/or link to) an example using Earthly with this CI.

### Notes

Commentary on the example to make it production ready.
