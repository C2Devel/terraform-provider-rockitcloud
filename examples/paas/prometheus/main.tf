# --------------------------------------------------------------------------
# Prometheus PaaS service — single node
#
# class defaults to "monitoring".
# The provider schema requires data_volume.
# --------------------------------------------------------------------------

resource "aws_key_pair" "example" {
  key_name   = "${var.paas_service_name}-key"
  public_key = var.ssh_public_key
}

resource "aws_paas_service" "prometheus" {
  depends_on = [aws_route.default_route]

  name              = var.paas_service_name
  high_availability = false
  instance_type     = var.instance_type

  root_volume {
    type = "gp2"
    size = 32
  }

  data_volume {
    type = "gp2"
    size = 32
  }

  delete_interfaces_on_destroy = true
  security_group_ids           = [aws_vpc.example.default_security_group_id]
  subnet_ids                   = [aws_subnet.example.id]

  ssh_key_name = aws_key_pair.example.key_name

  prometheus {
    remote_write_receiver = false
  }
}

# --------------------------------------------------------------------------
# Scrape jobs
# --------------------------------------------------------------------------

resource "aws_paas_prometheus_scrape_job" "app" {
  service_id = aws_paas_service.prometheus.id

  name = "appmetrics"
  targets = [
    "app-01.internal:9090",
    "app-02.internal:9090",
  ]
  labels = {
    env  = "qa"
    team = "platform"
  }
}

resource "aws_paas_prometheus_scrape_job" "node" {
  service_id = aws_paas_service.prometheus.id

  name = "nodeexporter"
  targets = [
    "node-01.internal:9100",
  ]
}

# --------------------------------------------------------------------------
# Notification channels (webhook — no external credentials required)
# --------------------------------------------------------------------------

resource "aws_paas_prometheus_notification_channel" "webhook" {
  service_id = aws_paas_service.prometheus.id

  name          = "opswebhook"
  type          = "webhook"
  url           = var.webhook_url
  max_alerts    = 10
  is_default    = true
  send_resolved = true
}

# --------------------------------------------------------------------------
# Alerting routes
# --------------------------------------------------------------------------

resource "aws_paas_prometheus_route" "critical" {
  service_id = aws_paas_service.prometheus.id

  name     = "criticalalerts"
  receiver = aws_paas_prometheus_notification_channel.webhook.name
  matchers = ["severity = \"critical\""]

  group_by        = ["alertname", "instance"]
  group_wait      = "30s"
  group_interval  = "5m"
  repeat_interval = "3h"
}

resource "aws_paas_prometheus_route" "warning" {
  service_id = aws_paas_service.prometheus.id

  name     = "warningalerts"
  receiver = aws_paas_prometheus_notification_channel.webhook.name
  matchers = ["severity = \"warning\""]
  continue = true
}

# --------------------------------------------------------------------------
# Outputs
# --------------------------------------------------------------------------

output "prometheus_service_id" {
  value = aws_paas_service.prometheus.id
}

output "prometheus_service_status" {
  value = aws_paas_service.prometheus.status
}

output "prometheus_service_type" {
  value = aws_paas_service.prometheus.service_type
}

output "prometheus_service_class" {
  value = aws_paas_service.prometheus.service_class
}

output "prometheus_high_availability" {
  value = aws_paas_service.prometheus.high_availability
}

output "prometheus_endpoints" {
  value = aws_paas_service.prometheus.endpoints
}

output "vpc_id" {
  value = aws_vpc.example.id
}

output "subnet_id" {
  value = aws_subnet.example.id
}

output "security_group_id" {
  value = aws_vpc.example.default_security_group_id
}

output "ssh_key_name" {
  value = aws_key_pair.example.key_name
}

output "scrape_job_app_id" {
  value = aws_paas_prometheus_scrape_job.app.job_id
}

output "scrape_job_node_id" {
  value = aws_paas_prometheus_scrape_job.node.job_id
}

output "channel_webhook_id" {
  value = aws_paas_prometheus_notification_channel.webhook.channel_id
}

output "route_critical_id" {
  value = aws_paas_prometheus_route.critical.route_id
}

output "route_warning_id" {
  value = aws_paas_prometheus_route.warning.route_id
}
