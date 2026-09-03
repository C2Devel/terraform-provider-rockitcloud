resource "aws_vpc" "example" {
  cidr_block = "172.16.0.0/16"

  tags = {
    Name = "terraform-paas-prometheus-example"
  }
}

resource "aws_subnet" "example" {
  vpc_id            = aws_vpc.example.id
  cidr_block        = cidrsubnet(aws_vpc.example.cidr_block, 4, 1)
  availability_zone = var.availability_zone

  tags = {
    Name = "terraform-paas-prometheus-example"
  }
}

resource "aws_internet_gateway" "example" {
  vpc_id = aws_vpc.example.id

  tags = {
    Name = "terraform-paas-prometheus-example"
  }
}

resource "aws_nat_gateway" "example" {
  depends_on = [aws_internet_gateway.example]

  vpc_id = aws_vpc.example.id

  tags = {
    Name = "terraform-paas-prometheus-example"
  }
}

resource "aws_route" "default_route" {
  route_table_id         = aws_vpc.example.main_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = aws_nat_gateway.example.id
}
