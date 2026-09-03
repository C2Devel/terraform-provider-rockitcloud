data "aws_ami" "selected" {
  most_recent = true

  name_regex = "Ubuntu*"

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["k2"]
}

resource "aws_instance" "example" {
  instance_type = "c5.large"
  ami           = data.aws_ami.selected.id
  subnet_id     = aws_subnet.example.id

  tags = {
    Name = "terraform-count-example"
  }

  # This will create 4 instances
  count = 4
}
