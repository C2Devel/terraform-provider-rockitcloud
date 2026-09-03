resource "aws_key_pair" "example" {
  key_name   = "${var.service_name}-key"
  public_key = var.ssh_public_key
}

# HA Kafka cluster with dedicated coordinator nodes.
#
# K2 Cloud Kafka always requires the coordinator role. There are two ways to
# place coordinators:
#
#   1. Dedicated coordinator nodes (this example): use the "coordinator" block.
#      Produces 3 broker nodes + 3 coordinator nodes.
#   2. Combined broker+coordinator nodes: set additional_roles = ["coordinator"].
#      Produces 3 combined nodes (HA) or 1 combined node (non-HA).
#
# The two options are mutually exclusive. See the K2 Cloud documentation for
# details: https://docs.k2.cloud/en/api/paas/actions/CreateService.html
resource "aws_paas_service" "example" {
  depends_on = [aws_route.default_route]

  name              = var.service_name
  high_availability = true
  instance_type     = var.instance_type

  root_volume {
    type = "gp2"
    size = 32
  }

  # Kafka data volumes require a minimum size of 64 GiB.
  data_volume {
    type = "gp2"
    size = 64
  }

  security_group_ids           = [aws_vpc.example.default_security_group_id]
  subnet_ids                   = [aws_subnet.example.id]
  ssh_key_name                 = aws_key_pair.example.key_name
  delete_interfaces_on_destroy = true

  coordinator {
    instance_type    = var.instance_type
    root_volume_type = "gp2"
    root_volume_size = 32
    data_volume_type = "gp2"
    data_volume_size = 64
  }

  kafka {
    version = var.kafka_version
  }
}

output "kafka_service_id" {
  description = "ID of the created Kafka PaaS service."
  value       = aws_paas_service.example.id
}

output "kafka_service_status" {
  description = "Status of the Kafka PaaS service."
  value       = aws_paas_service.example.status
}

output "kafka_endpoints" {
  description = "Broker endpoints clients can connect to."
  value       = aws_paas_service.example.endpoints
}
