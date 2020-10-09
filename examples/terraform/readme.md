# Terraform Example

This is to demonstrate how you might combine Earthly and Terraform.

## +localstack

`earthly -P +localstack`

This target runs localstack and applies the Terraform to it. This is easy to do because it doesn't require a cloud to actually test applying your Terraform.

Requires priveleged mode to spin up the Localstack using DIND.

## +plan

`earthly --build_arg AWS_ACCESS_KEY_ID --build_arg AWS_SECRET_ACCESS_KEY +plan`

This target actually plans the Terraform against your cloud, if you pass in valid credentials. The region is optional, and can be overridden. Defaults to `us-east-1`.

## +apply

`earthly --push --build_arg AWS_ACCESS_KEY_ID --build_arg AWS_SECRET_ACCESS_KEY +apply`

This target actually applies the Terraform against your cloud, if you pass in valid credentials. Like in `+plan`, the region is optional, and can be overridden.

Requires `--push` to actually run the apply. Saves the `.tfstate` files as artifacts. If it is your first run, you will need to comment out saving the `.tfstate.backup`, since you will not have one yet.