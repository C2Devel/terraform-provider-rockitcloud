# Prometheus PaaS example

This example creates a single-node K2 Cloud Prometheus service and its
Terraform-managed configuration:

- two scrape jobs;
- one webhook notification channel;
- two alerting routes.

The VPC has an Internet gateway and a default route because K2 Cloud requires
Internet access in the Prometheus subnet.

## Run

The example takes credentials from a `c2rc.sh` file.
Get the file for your project and place it in this directory before running the example.

Generate a disposable SSH key and a unique service name:

```bash
export PROMETHEUS_EXAMPLE_KEY="${TMPDIR:-/tmp}/terraform-paas-prometheus"
ssh-keygen -q -t ed25519 -f "${PROMETHEUS_EXAMPLE_KEY}" -N ""
export TF_VAR_ssh_public_key
TF_VAR_ssh_public_key="$(tr -d '\n' < "${PROMETHEUS_EXAMPLE_KEY}.pub")"
export TF_VAR_paas_service_name="tf-prom-$(date -u +%s)"
export TF_VAR_webhook_url="https://prometheus-webhook.test/alerts"
```

The `.test` domain is reserved and cannot receive alert payloads. Set
`TF_VAR_webhook_url` to an HTTPS endpoint controlled by you only when you want
to test real notification delivery.

Run Terraform from this directory:

```bash
terraform init
terraform validate
terraform apply
terraform plan -detailed-exitcode
terraform destroy
```

An idempotent plan exits with code `0`. Always destroy the test deployment and
verify that its PaaS service and networking resources have disappeared.

## Contract notes

- Prometheus supports only `high_availability = false`.
- The service class defaults to `monitoring`.
- The provider schema requires `data_volume`.
- `remote_write_receiver` is editable in place on `paas_v4_0` and later.
- Child resource names are create-time-only and therefore force replacement.
  A scrape job name must start with a Latin letter, contain only Latin letters
  and digits, and be at most 32 characters long.
- A route's `receiver` is a notification channel name, not its ID.
- Every scrape job needs at least one target, and every route needs at least
  one matcher.
- Webhook is used here so that Telegram and SMTP secrets are unnecessary.
- A standalone Remote Write target resource is not exposed until the K2 Go SDK
  includes its API operations.
