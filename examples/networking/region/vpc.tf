data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_vpc" "example" {
  cidr_block = var.vpc_cidr

  tags = {
    Name = "terraform-network-example"
  }
}

resource "aws_subnet" "example" {
  # Create a subnet in up to two availability zones of the region
  count = min(2, length(data.aws_availability_zones.available.names))

  vpc_id            = aws_vpc.example.id
  availability_zone = data.aws_availability_zones.available.names[count.index]
  cidr_block        = cidrsubnet(var.vpc_cidr, 8, count.index)

  tags = {
    Name = "terraform-network-example"
  }
}

resource "aws_route_table" "example" {
  vpc_id = aws_vpc.example.id

  tags = {
    Name = "terraform-network-example"
  }
}

resource "aws_route_table_association" "example" {
  count = length(aws_subnet.example)

  subnet_id      = aws_subnet.example[count.index].id
  route_table_id = aws_route_table.example.id
}
