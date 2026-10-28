variable "region" {
  description = "Region to create the infrastructure in."
  type        = string
  default     = "ru-msk"
}

variable "availability_zone" {
  description = "Availability zone for the subnet (must be inside region)."
  type        = string
  default     = "ru-msk-vol52"
}

variable "service_name" {
  description = "PaaS service name (must be unique in the project)."
  type        = string
  default     = "terraform-pgsql-resource-example"
}

variable "pgsql_version" {
  description = "PostgreSQL version to deploy."
  type        = string
  default     = "16"
}

variable "instance_type" {
  description = "Instance type for PostgreSQL nodes."
  type        = string
  default     = "c5.large"
}

variable "app_owner_password" {
  description = "Password of the app_owner user. Must not contain ', \", ` and \\."
  type        = string
  sensitive   = true
}

variable "app_reader_password" {
  description = "Password of the app_reader user. Must not contain ', \", ` and \\."
  type        = string
  sensitive   = true
}

variable "monitoring_password" {
  description = "Password of the monitoring user. Must not contain ', \", ` and \\."
  type        = string
  sensitive   = true
}
