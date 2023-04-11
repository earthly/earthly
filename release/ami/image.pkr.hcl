packer {
  required_plugins {
    amazon = {
      version = ">= 1.1.5"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

variable "earthly_version" {
  type = string
}

locals {
  timestamp       = replace(timestamp(), ":", "-")
  earthly_version = trimprefix(var.earthly_version, "v")
}

source "amazon-ebs" "x86_64" {
  ami_name      = "earthly-${var.earthly_version}-amzn-amd64-${local.timestamp}"
  instance_type = "t2.micro"
  region        = "us-west-2"
  source_ami_filter {
    filters = {
      name                = "amzn2-ami-kernel-5.10-hvm-2.0.*-x86_64-gp2"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["137112412989"]
  }
  ssh_username           = "ec2-user"
  ssh_read_write_timeout = "5m" # Allow reboots
}

source "amazon-ebs" "arm64" {
  ami_name      = "earthly-${var.earthly_version}-amzn-arm64-${local.timestamp}"
  instance_type = "a1.medium"
  region        = "us-west-2"
  source_ami_filter {
    filters = {
      name                = "amzn2-ami-kernel-5.10-hvm-2.0.*-arm64-gp2"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["137112412989"]
  }
  ssh_username           = "ec2-user"
  ssh_read_write_timeout = "5m" # Allow reboots
}

build {
  name    = "earthly-build"
  sources = [
    "source.amazon-ebs.x86_64",
    "source.amazon-ebs.arm64"
  ]

  # https://developer.hashicorp.com/packer/docs/debugging#issues-installing-ubuntu-packages
  provisioner "shell" {
    inline = [
      "sudo cloud-init status --wait"
    ]
  }

  provisioner "file" {
    source      = "install.sh"
    destination = "/tmp/install.sh"
    max_retries = 10
  }
  provisioner "shell" {
    environment_vars = [
      "EARTHLY_VERSION=${local.earthly_version}"
    ]
    inline = [
      "cd /tmp",
      "chmod +x install.sh && ./install.sh",
    ]
  }

  # We need to reboot since we need that to finish the docker installation, hence the sleep part
  provisioner "shell" {
    expect_disconnect = true
    inline = [
      "sudo reboot now",
    ]
    pause_after  = "60s"
  }

  provisioner "file" {
    source      = "configure.sh"
    destination = "/tmp/configure.sh"
    max_retries = 10
  }
  provisioner "file" {
    source      = "cleanup.sh"
    destination = "/tmp/cleanup.sh"
    max_retries = 10
  }
  provisioner "shell" {
    inline = [
      "cd /tmp",
      "chmod +x configure.sh && ./configure.sh",
      "chmod +x cleanup.sh && ./cleanup.sh"
    ]
  }
}
