# Single-node PostgreSQL service.
#
# aws_paas_pgsql manages only PostgreSQL services. Users and databases are always
# managed by the standalone aws_paas_pgsql_user and aws_paas_pgsql_database resources
# (see pgsql-users.tf and pgsql-databases.tf), so the pgsql block has no user and
# database blocks and no flags are needed.
resource "aws_paas_pgsql" "example" {
  depends_on = [aws_route.default_route]

  name          = var.service_name
  instance_type = var.instance_type

  root_volume {
    type = "gp2"
    size = 32
  }

  data_volume {
    type = "gp2"
    size = 32
  }

  security_group_ids           = [aws_vpc.example.default_security_group_id]
  subnet_ids                   = [aws_subnet.example.id]
  delete_interfaces_on_destroy = true

  pgsql {
    version = var.pgsql_version

    # PostgreSQL parameters can be updated in place.
    max_connections = 200
  }
}

output "pgsql_service_id" {
  description = "ID of the created PostgreSQL PaaS service."
  value       = aws_paas_pgsql.example.id
}

output "pgsql_service_status" {
  description = "Status of the PostgreSQL PaaS service."
  value       = aws_paas_pgsql.example.status
}

output "pgsql_endpoints" {
  description = "Endpoints clients can connect to."
  value       = aws_paas_pgsql.example.endpoints
}
