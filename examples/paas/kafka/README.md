# PaaS Kafka Example

This example deploys a fully managed Apache Kafka cluster on K2 Cloud PaaS
using the `aws_paas_service` and `aws_paas_kafka_topic` resources.

It creates:

* a VPC with a single subnet, an internet gateway and a NAT gateway (required by PaaS);
* an SSH key pair used for node access;
* an HA Kafka cluster (3 brokers + 3 coordinators) running the configured
  Kafka version;
* three Kafka topics (`events`, `logs`, `state`) showcasing different
  retention and compression configurations.

Layout:

| File | Purpose |
|------|---------|
| `providers.tf`     | Terraform/provider versions and provider config |
| `variables.tf`     | Input variables (region, instance type, SSH key, ...) |
| `vpc.tf`           | VPC + subnet + internet gateway + NAT gateway + default route |
| `kafka-cluster.tf` | SSH key pair, Kafka PaaS service, outputs |
| `kafka-topics.tf`  | `aws_paas_kafka_topic` resources, outputs |
| `terraform.tfvars.example` | Template for `terraform.tfvars` |

## Running the example

The example takes credentials from a `c2rc.sh` file.
Get the file for your project and place it in this directory before running the example.

```bash
cp terraform.tfvars.example terraform.tfvars
# edit terraform.tfvars and at least set ssh_public_key

terraform init
terraform apply
```

Provisioning a Kafka cluster typically takes 10-20 minutes. While it runs,
Terraform will block until the platform reports the service as `READY`.

When the apply is complete the following outputs are exposed:

```bash
terraform output kafka_service_id
terraform output kafka_endpoints
terraform output topic_events_id
```

## Destroying the example

```bash
terraform destroy
```

## Cluster variants

The example deploys an HA cluster with dedicated coordinator nodes (the
`coordinator { ... }` block). Two more variants are documented in the
provider source tree under `examples/paas-kafka/`:

* `service-ha-combined.tf.example`: HA cluster where broker and coordinator
  roles share the same nodes (`additional_roles = ["coordinator"]`).
* `service-non-ha.tf.example`: single-node cluster
  (`high_availability = false` plus `additional_roles = ["coordinator"]`).

K2 Cloud Kafka always requires the coordinator role; the
[CreateService docs](https://docs.k2.cloud/en/api/paas/actions/CreateService.html)
and `aws_paas_service` validation enforce this on the client side.

## Notes

* `data_volume.size` must be at least 64 GiB for Kafka (validated locally).
* `partitions` and `replication_factor` for a topic must not exceed the
  number of broker nodes (3 for HA, 1 for non-HA). The platform rejects
  invalid values with `InvalidParameterValue`.
* Supported Kafka versions are listed at
  https://docs.k2.cloud/en/api/paas/parameters/kafka.html
  and validated locally by the provider.
