# PostgreSQL databases managed by the provider.
#
# Notes:
#   * Refer to users with resource attributes, so that Terraform creates users
#     before databases and deletes databases before their owners.
#   * users must not contain the owner: the owner already has access to the database.
#   * Changing encoding or locale deletes the database with all its data
#     and creates a new one.

resource "aws_paas_pgsql_database" "app" {
  service_id = aws_paas_pgsql.example.id
  name       = "app"
  owner      = aws_paas_pgsql_user.app_owner.name
  encoding   = "UTF8"
  locale     = "ru_RU.UTF-8"
  extensions = ["pgcrypto", "pg_trgm", "uuid-ossp"]
  users      = [aws_paas_pgsql_user.app_reader.name]
}

# Database without additional users: only the owner has access.
resource "aws_paas_pgsql_database" "reports" {
  service_id = aws_paas_pgsql.example.id
  name       = "reports"
  owner      = aws_paas_pgsql_user.app_owner.name
}

output "database_ids" {
  description = "IDs of the databases in the service."
  value = {
    app     = aws_paas_pgsql_database.app.database_id
    reports = aws_paas_pgsql_database.reports.database_id
  }
}
