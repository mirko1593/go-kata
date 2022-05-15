module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  version = "~> 3.0.0"

  name = "cloudcasts-${var.infra_env}-vpc"
  cidr = var.vpc_cidr

  azs = var.azs

  enable_nat_gateway = true
  single_nat_gateway = true
  one_nat_gateway_per_az = false

  private_subnets = var.private_subnets
  public_subnets = var.public_subnets

  tags = {
    Name = "cloudcasts-${var.infra_env}-vpc"
    Project = "cloudcasts.io"
    Environment = var.infra_env
    ManangedBy = "terraform"
  }
}
