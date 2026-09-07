# Default security group to access
# the instances over SSH and HTTP
resource "aws_security_group" "example" {
  name        = "terraform-elb-example"
  description = "Used in the terraform"
  vpc_id      = aws_vpc.example.id

  # SSH access from anywhere
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # HTTP access from anywhere
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Outbound internet access
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "terraform-elb-example"
  }
}

resource "aws_lb" "example" {
  name               = "terraform-elb-example"
  internal           = false
  load_balancer_type = "application"
  subnets            = [aws_subnet.example.id]

  # Ensure the VPC has an internet gateway with a configured route or this step will fail
  depends_on = [aws_route.default_route]

  tags = {
    Name = "terraform-elb-example"
  }
}

resource "aws_lb_target_group" "example" {
  name = "terraform-elb-example"

  target_type = "instance"
  port        = 80
  protocol    = "HTTP"
  vpc_id      = aws_vpc.example.id

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    interval            = 30
  }

  tags = {
    Name = "terraform-elb-example"
  }
}

resource "aws_lb_target_group_attachment" "example" {
  target_group_arn = aws_lb_target_group.example.arn
  target_id        = aws_instance.example.id
}

resource "aws_lb_listener" "example" {
  load_balancer_arn = aws_lb.example.arn

  port     = 80
  protocol = "HTTP"

  default_action {
    type = "forward"

    forward {
      target_group {
        arn = aws_lb_target_group.example.arn
      }
    }
  }

  tags = {
    Name = "terraform-elb-example"
  }
}

data "aws_ami" "selected" {
  most_recent = true

  name_regex = "Ubuntu 24.04 [Cloud Image]*"

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["k2"]
}

resource "aws_instance" "example" {
  depends_on = [aws_route.default_route]

  instance_type = "c5.large"
  ami           = data.aws_ami.selected.id

  # The name of the SSH keypair created in the cloud console
  key_name = var.key_name

  # Once the instance is created, a remote provisioner is launched on it.
  # In this case, it installs nginx and starts it. By default,
  # this should be on port 80
  user_data = file("userdata.sh")

  subnet_id              = aws_subnet.example.id
  vpc_security_group_ids = [aws_security_group.example.id]

  tags = {
    Name = "terraform-elb-example"
  }
}
