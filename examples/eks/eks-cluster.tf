resource "aws_security_group" "example" {
  name        = "terraform-eks-example"
  description = "Cluster communication with worker nodes"
  vpc_id      = aws_vpc.example.id

  ingress {
    description = "Allow workstation to communicate with the cluster API Server"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [local.workstation-external-cidr]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "terraform-eks-example"
  }
}

# resource "aws_security_group_rule" "example" {
#   cidr_blocks       = [local.workstation-external-cidr]
#   description       = "Allow workstation to communicate with the cluster API Server"
#   from_port         = 443
#   protocol          = "tcp"
#   security_group_id = aws_security_group.example.id
#   to_port           = 443
#   type              = "ingress"
# }

resource "aws_eks_cluster" "example" {
  depends_on = [aws_route.default_route]

  name    = "terraform-eks-example"
  version = "1.33.1"

  vpc_config {
    security_group_ids = [aws_security_group.example.id]
    subnet_ids         = aws_subnet.example[*].id
  }

  tags = {
    Name = "terraform-eks-example"
  }
}

data "aws_eks_cluster_kubeconfig" "example_config" {
  name = "terraform-eks-example"
}

output "kubeconfig" {
  value     = data.aws_eks_cluster_kubeconfig.example_config.kubeconfig
  sensitive = true
}
