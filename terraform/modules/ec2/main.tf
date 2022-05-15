  resource "random_shuffle" "subnets" {
    input = var.subnets
    result_count = 1
  }

resource "aws_instance" "cloudcasts_web" {
  ami = var.instance_ami
  instance_type = var.instance_type

  root_block_device {
    volume_size = var.instance_root_device_size
    volume_type = "gp3"
  }

  subnet_id = random_shuffle.subnets.result[0]
  vpc_security_group_ids = var.security_groups

  tags = merge(
    {
      Name = "cloudcasts-${var.infra_env}-web"
      Role = var.infra_role
      Project = "cloudcasts"
      Environment = var.infra_env
      ManagedBy = "terraform"
    },
    var.tags
  )
}

resource "aws_eip" "cloudcasts_addr" {
  count = (var.create_eip) ? 1 : 0

  vpc = true

  # lifecycle {
  #   prevent_destroy = true
  # }

  tags = {
    Name = "cloudcasts-${var.infra_env}-web"
    Role = var.infra_role
    Project = "cloudcasts"
    Environment = var.infra_env
    ManagedBy = "terraform"
  }
}

resource "aws_eip_association" "cloudcast_web_eip_assoc" {
  count = (var.create_eip) ? 1 : 0

  instance_id = aws_instance.cloudcasts_web.id

  allocation_id = aws_eip.cloudcasts_addr[0].id
}
