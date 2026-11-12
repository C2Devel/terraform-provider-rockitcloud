# PostgreSQL users managed by the provider.
#
# Notes:
#   * Each operation switches the service to UPDATING for about 45 seconds.
#     Operations on users and databases of one service are performed one at a time,
#     even if Terraform starts them in parallel.
#   * roles requires PaaS environment paas_v4_3 or later.
#   * Passwords are stored in the Terraform state.

# Owner of the databases in pgsql-databases.tf.
resource "aws_paas_pgsql_user" "app_owner" {
  service_id = aws_paas_pgsql.example.id
  name       = "app_owner"
  password   = var.app_owner_password
}

# Application user with access to the app database.
resource "aws_paas_pgsql_user" "app_reader" {
  service_id = aws_paas_pgsql.example.id
  name       = "app_reader"
  password   = var.app_reader_password
}

# Service-wide monitoring user with predefined PostgreSQL roles.
resource "aws_paas_pgsql_user" "monitoring" {
  service_id = aws_paas_pgsql.example.id
  name       = "monitoring"
  password   = var.monitoring_password
  roles      = ["pg_monitor", "pg_read_all_stats"]
}

output "user_ids" {
  description = "IDs of the users in the service."
  value = {
    app_owner  = aws_paas_pgsql_user.app_owner.user_id
    app_reader = aws_paas_pgsql_user.app_reader.user_id
    monitoring = aws_paas_pgsql_user.monitoring.user_id
  }
}
