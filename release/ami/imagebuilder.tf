terraform {
  required_version = ">= 0.13"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.12"
    }
  }
  backend "s3" {
    bucket               = "earthly-terraform-state"
    key                  = "state"
    region               = "us-west-2"
    workspace_key_prefix = "imagebuilder"
  }
}

data "aws_ami" "ecs_ami" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-ecs-hvm-*-x86_64-ebs"]
  }
}

data "aws_imagebuilder_component" "docker" {
  arn = "arn:aws:imagebuilder:us-west-2:aws:component/docker-ce-linux/1.0.0"
}

resource "aws_imagebuilder_component" "earthly" {
  for_each = toset(split("\n", trimspace(file("./versions"))))

  data = yamlencode({
    phases = [{
      name = "build"
      steps = [
        {
          action = "WebDownload"
          inputs: [{
            source: "https://github.com/earthly/earthly/releases/download/v${each.key}/earthly-linux-amd64"
            destination: "/usr/local/bin/earthly"
          }]
          name      = "download_earthly"
          onFailure = "Abort"
        },
        {
        action = "ExecuteBash"
        inputs = {

          commands = [
            "chmod +x /usr/local/bin/earthly",
            "/usr/local/bin/earthly bootstrap"
          ]
        }
        name      = "install_earthly"
        onFailure = "Abort"
      }]
    }]
    schemaVersion = 1.0
  })
  name        = "Install Earthly"
  description = "A component to install Earthly"
  platform    = "Linux"
  version     = replace(each.key, "v", "") # Not the Earthly version, this is for the YAML document
}

resource "aws_imagebuilder_image_recipe" "earthly" {
  for_each = aws_imagebuilder_component.earthly

  name         = "Earthly ${replace(each.key, ".", "_")}"
  parent_image = data.aws_ami.ecs_ami.image_id
  version      = each.value.version
  component {
    component_arn = data.aws_imagebuilder_component.docker.arn
  }
  component {
    component_arn = each.value.arn
  }
}

resource "aws_imagebuilder_infrastructure_configuration" "earthly" {
  instance_profile_name = "EC2InstanceProfileForImageBuilder"
  name                  = "Earthly Image Build Configuration"
}

resource "aws_imagebuilder_distribution_configuration" "earthly" {
  for_each = aws_imagebuilder_image_recipe.earthly

  name = "earthly-${replace(each.key, ".", "_")}"
  description = "Distribution settings for Earthly ${each.key} AMI"
  distribution {
    region = "us-west-2"
    ami_distribution_configuration {
      name = "earthly-${replace(each.key, ".", "_")}-{{ imagebuilder:buildDate }}"
    }
  }
}

resource "aws_imagebuilder_image_pipeline" "earthly" {
  for_each = aws_imagebuilder_image_recipe.earthly

  image_recipe_arn                 = each.value.arn
  infrastructure_configuration_arn = aws_imagebuilder_infrastructure_configuration.earthly.arn
  distribution_configuration_arn   = aws_imagebuilder_distribution_configuration.earthly[each.key].arn
  name                             = "earthly-${replace(each.key, ".", "_")}"
  description                      = "Builds an AMI with Earthly ${each.key}"
}

output "pipelines" {
  value       = aws_imagebuilder_image_pipeline.earthly
  description = "The pipelines we might build; used so we can trigger them manually from another Earthfile."
}