resource "aws_vpc" "example" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "terraform-count-example"
  }
}

resource "aws_subnet" "example" {
  cidr_block = "10.0.0.0/24"
  vpc_id     = aws_vpc.example.id

  tags = {
    Name = "terraform-count-example"
  }
}
