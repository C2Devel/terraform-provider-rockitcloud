---
subcategory: "EKS (Elastic Kubernetes)"
layout: "aws"
page_title: "aws_eks_cluster"
description: |-
  Provides information about an EKS cluster.
---

# Data Source: aws_eks_cluster

Provides information about an EKS cluster.

## Example usage

```terraform
data "aws_eks_cluster" "example" {
  name = "example"
}
```

## Argument reference

The following arguments are supported:

* `name` - (Required, String) The name of the cluster.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The ARN of the cluster.
* `certificate_authority` - ([Block](#certificate_authority)) Nested attribute containing `certificate-authority-data` for your cluster.
* `created_at` - (String) The Unix epoch time stamp in seconds for when the cluster was created.
* `id` - (String) The name of the cluster.
* `kubernetes_network_config` - ([Block](#kubernetes_network_config)) The Kubernetes network configuration.
* `legacy_cluster_params` - ([Block](#legacy_cluster_params)) The parameters for fine-tuning the Kubernetes cluster.
* `platform_version` - (String) The platform version for the cluster.
* `status` - (String) The status of the EKS cluster.
    * _Valid values:_ `CLAIMED`, `CREATING`, `DELETED`, `DELETING`, `ERROR`, `MODIFYING`, `PENDING`, `PROVISIONING`, `READY`, `REPAIRING`
* `tags` - (Map of strings) Key-value pairs assigned to the cluster.
* `version` - (String) The Kubernetes server version for the cluster.
* `vpc_config` - ([Block](#vpc_config)) The VPC configuration for the cluster.

#### certificate_authority

The `certificate_authority` block has the following structure:

* `data` - (String) The base64 encoded certificate data required to communicate with your cluster.
    Add this to the `certificate-authority-data` section of the `kubeconfig` file for your cluster.

#### kubernetes_network_config

The `kubernetes_network_config` block has the following structure:

* `ip_family` - The IP family used to assign Kubernetes pod and service addresses.
* `service_ipv4_cidr` - The CIDR block to assign Kubernetes service IP addresses from.

#### legacy_cluster_params

The `legacy_cluster_params` block has the following structure:

* `docker_registry_config` - ([Block](#docker_registry_config)) The configuration of the Docker Registry.
* `ebs_provider_config` - ([Block](#ebs_provider_config)) The configuration of the EBS Provider.
* `ingress_config` - ([Block](#ingress_config)) The configuration of the Ingress controller.
* `master_config` - ([Block](#master_config)) The configuration of the master node of the cluster.
* `nlb_provider_config` - ([Block](#nlb_provider_config)) The configuration of the NLB Provider.
* `placement_config` - ([Block](#placement_config)) The placement of the cluster.
* `user_data_config` - ([Block](#user_data_config)) The configuration of the cluster user data.

##### docker_registry_config

The `docker_registry_config` block has the following structure:

* `volume_iops` - The number of read/write operations per second for the Docker Registry volume.
* `volume_size` - The size of the Docker Registry volume in GiB.
* `volume_type` - The type of the Docker Registry volume.

##### ebs_provider_config

The `ebs_provider_config` block has the following structure:

* `ebs_user` - The EBS Provider user name.

##### ingress_config

The `ingress_config` block has the following structure:

* `instance_type` - The instance type of the Ingress controller.
* `public_ip` - The public IP address at which the Ingress controller can be accessed.
* `volume_iops` - The number of read/write operations per second for the Ingress controller volume.
* `volume_size` - The size of the Ingress controller volume in GiB.
* `volume_type` - The type of the Ingress controller volume.

##### master_config

The `master_config` block has the following structure:

* `high_availability` - Indicates whether to deploy a high-availability cluster.
* `instance_type` - The instance type of the master node.
* `public_ip` - The public IP address at which the master node can be accessed.
* `volume_iops` - The number of read/write operations per second for the master node volume.
* `volume_size` - The size of the master node volume in GiB.
* `volume_type` - The type of the master node volume.

##### nlb_provider_config

The `nlb_provider_config` block has the following structure:

* `nlb_user` - The NLB Provider user name.

#### placement_config

The `placement_config` block has the following structure:

* `affinity` - The affinity setting for an instance on a dedicated host.
    * _Constraints:_ The parameter could be set to `host` only if `tenancy` is `host`
    * _Valid values:_ `default`, `host`
    * _Default value:_ `default`
* `host_id` - The ID of the dedicated host for the instance.
* `tenancy` - The tenancy of the instance (if the instance is running in a VPC).
    * _Valid values:_ `default`, `host`
    * _Default value:_ `default`

#### user_data_config

The `user_data_config` block has the following structure:

* `user_data` - User data.
* `user_data_content_type` - The type of `user_data`.
    * _Valid values:_ `cloud-config`, `x-shellscript`

#### vpc_config

The `vpc_config` block has the following structure:

* `cluster_security_group_id` - The cluster security group that was created by the cloud for the cluster.
* `endpoint_private_access` - (Boolean) Indicates whether the endpoint private access is enabled.
* `endpoint_public_access` - (Boolean) Indicates whether the endpoint public access is enabled.
* `public_access_cidrs` - (Set of strings) The list of CIDR blocks which can access the cluster endpoint.
* `security_group_ids` - List of security group IDs.
* `subnet_ids` - List of subnet IDs.
* `vpc_id` - The VPC associated with your cluster.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`enabled_cluster_log_types`, `encryption_config`, `endpoint`, `identity`, `role_arn`.