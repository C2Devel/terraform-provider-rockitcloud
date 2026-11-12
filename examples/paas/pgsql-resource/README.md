# PaaS PostgreSQL Example with aws_paas_pgsql

This example deploys a managed PostgreSQL service on K2 Cloud PaaS
using the PostgreSQL-only `aws_paas_pgsql` resource,
and manages its users and databases with the `aws_paas_pgsql_user` and `aws_paas_pgsql_database` resources.

The same service can be created with `aws_paas_service`, see the [pgsql](../pgsql) example.
`aws_paas_pgsql` has only the `pgsql` block, and users and databases are always managed by the standalone resources,
so the `users_managed_externally` and `databases_managed_externally` flags are not needed.

It creates:

* a VPC with a single subnet, an internet gateway and a NAT gateway (required by PaaS);
* a single-node PostgreSQL service;
* three users: `app_owner` (owner of the databases), `app_reader` (application user)
  and `monitoring` (service-wide user with the `pg_monitor` and `pg_read_all_stats` roles);
* two databases: `app` (with extensions and access for `app_reader`) and `reports` (accessible only by its owner).

Layout:

| File | Purpose |
|------|---------|
| `providers.tf`        | Terraform/provider versions and provider config |
| `variables.tf`        | Input variables (region, instance type, passwords, ...) |
| `vpc.tf`              | VPC + subnet + internet gateway + NAT gateway + default route |
| `pgsql-service.tf`    | `aws_paas_pgsql` service, outputs |
| `pgsql-users.tf`      | `aws_paas_pgsql_user` resources, outputs |
| `pgsql-databases.tf`  | `aws_paas_pgsql_database` resources, outputs |
| `terraform.tfvars.example` | Template for `terraform.tfvars` |

## Running the example

The example takes credentials from a `c2rc.sh` file.
Get the file for your project and place it in this directory before running the example.

```bash
cp terraform.tfvars.example terraform.tfvars
# edit terraform.tfvars and set the passwords

terraform init
terraform apply
```

**Warning:** `terraform.tfvars` contains passwords. Do not commit it, or pass the passwords
with `TF_VAR_*` environment variables instead.

Provisioning the service typically takes 10-15 minutes. After that, the users and databases are created
one at a time, about 45 seconds each, so the whole apply takes a few more minutes.
Terraform creates the users before the databases, because the databases refer to them.

When the apply is complete the following outputs are exposed:

```bash
terraform output pgsql_service_id
terraform output pgsql_endpoints
terraform output user_ids
terraform output database_ids
```

## Destroying the example

```bash
terraform destroy
```

Terraform deletes the databases first, then the users and then the service.

## High availability

To deploy an HA cluster, change `aws_paas_pgsql`. `replication_mode` is required for an HA cluster:

```terraform
  high_availability   = true
  arbitrator_required = true # two nodes and an arbitrator; omit and set instance_qty = 3 for three nodes

  pgsql {
    version          = var.pgsql_version
    replication_mode = "synchronous"
  }
```

The nodes can be placed in one subnet, or in three subnets in different availability zones
(`subnet_ids = aws_subnet.example[*].id`).

## Moving a service from aws_paas_service

An existing service managed by `aws_paas_service` can be moved to `aws_paas_pgsql` without recreation
with the `removed` and `import` blocks. See
[Moving a service from aws_paas_service](../../../docs/resources/paas_pgsql.md#moving-a-service-from-aws_paas_service).

## Notes

* A database owner must not be listed in the `users` of the database.
* Changing `encoding` or `locale` of a database deletes the database with all its data and creates a new one.
* With the default timeouts, approximately 35–40 users and databases of one service can be created,
  modified or deleted in a single apply. Apply larger changes in batches or increase the `timeouts`.
* Supported PostgreSQL versions are listed at
  https://docs.k2.cloud/en/api/paas/parameters/pgsql.html
