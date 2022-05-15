terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}


provider "aws" {
  region  = "cn-northwest-1"
  profile = "default"
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["cloudcasts-${var.infra_env}-1648190357-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }

  owners = ["self"]
}

variable infra_env {
  type = string
  description = "infrastructure environment"
}

variable default_region {
  type = string
  description = "default region this infrastructure is in"
  default = "cn-northwest-1"
}

variable instance_type {
  type = string
  description = "ec2 web server size"
  default = "t3.small"
}


module "ec2_app" {
  source = "./modules/ec2"

  infra_env = var.infra_env
  infra_role = "app"
  instance_type = var.instance_type
  instance_ami = data.aws_ami.ubuntu.id
  # instance_root_device_size = 50

  subnets = keys(module.vpc.vpc_public_subnets) # Note: Public subnets 
  # security_groups = [] # TODO: Create security groups
  # instance_root_device_size = 12 
  security_groups = [module.vpc.security_group_public]
  tags = {
    Name = "cloudcasts-${var.infra_env}-web"
  }
  create_eip = true
}

module "ec2_worker" {
  source = "./modules/ec2"

  infra_env = var.infra_env
  infra_role = "worker"
  instance_type = var.instance_type
  instance_ami = data.aws_ami.ubuntu.id
  instance_root_device_size = 50

  subnets = keys(module.vpc.vpc_private_subnets) # Note: Public subnets 
  # security_groups = [] # TODO: Create security groups
  # instance_root_device_size = 12 
  security_groups = [module.vpc.security_group_private]
  tags = {
    Name = "cloudcasts-${var.infra_env}-worker"
  }
  create_eip = false
}

module "vpc" {
  source = "./modules/vpc"

  infra_env = var.infra_env

  vpc_cidr = "10.0.0.0/17"
}

# THIS MANUAL SCRIPT IS REPLACED WITH MODULE
#
# resource "aws_instance" "cloudcasts_web" {
#   ami = data.aws_ami.ubuntu.id

#   instance_type = var.instance_type


#   root_block_device {
#     volume_size = 8
#     volume_type = "gp3"
#   }

#   tags = {
#     Name        = "cloudcasts-${var.infra_env}-web"
#     Project     = "cloudcasts.io"
#     Environment = var.infra_env
#     ManagedBy   = "terraform"
#   }
# }

# resource "aws_eip" "cloudcasts_web_addr" {
#   vpc = true

#   # lifecycle {
#   #   prevent_destroy = true
#   # }

#   tags = {
#     Name        = "cloudcasts-staging-web"
#     Project     = "cloudcasts.io"
#     Environment = var.infra_env
#     ManagedBy   = "terraform"
#   }
# }

# resource "aws_eip_association" "cloudcasts_web_eip_association" {
#   instance_id = aws_instance.cloudcasts_web.id
#   allocation_id = aws_eip.cloudcasts_web_addr.id
# }

