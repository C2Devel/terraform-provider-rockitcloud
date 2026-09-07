output "vpc_id" {
  value = aws_vpc.example.id
}

output "subnet_ids" {
  value = aws_subnet.example[*].id
}

output "security_group_id" {
  value = aws_security_group.example.id
}
