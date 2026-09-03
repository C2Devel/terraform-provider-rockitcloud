resource "aws_ec2_transit_gateway" "example" {
  provider = aws.first

  tags = {
    Name = "terraform-tgw-example"
  }
}

resource "aws_ec2_transit_gateway_project_access" "example" {
  provider = aws.first

  transit_gateway_id = aws_ec2_transit_gateway.example.id
  account_id         = local.second_account_id
}

resource "aws_vpc" "example" {
  provider = aws.second

  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "terraform-tgw-example"
  }
}

data "aws_availability_zones" "available" {
  provider = aws.second

  state = "available"
}

resource "aws_subnet" "example" {
  provider = aws.second

  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = "10.0.0.0/24"
  vpc_id            = aws_vpc.example.id

  tags = {
    Name = "terraform-tgw-example"
  }
}

resource "aws_ec2_transit_gateway_vpc_attachment" "example" {
  provider = aws.second

  depends_on = [aws_ec2_transit_gateway_project_access.example]

  subnet_ids         = [aws_subnet.example.id]
  transit_gateway_id = aws_ec2_transit_gateway.example.id
  vpc_id             = aws_vpc.example.id

  tags = {
    Name = "terraform-tgw-example"
  }
}
