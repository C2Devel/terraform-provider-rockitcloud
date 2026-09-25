---
subcategory: "EKS (Elastic Kubernetes)"
layout: "aws"
page_title: "aws_eks_cluster"
description: |-
  Provides information about an EKS cluster.
---

# Data Source: aws_eks_cluster

Provides information about an EKS cluster.

## Example Usage

```terraform
data "aws_eks_cluster" "example" {
  name = "example"
}
```

## Argument Reference

* `name` - (Required) The name of the cluster.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - Cluster ID.
* `certificate_authority` - Nested attribute containing `certificate-authority-data` for your cluster.
    * `data` - The base64 encoded certificate data required to communicate with your cluster. Add this to the `certificate-authority-data` section of the `kubeconfig` file for your cluster.
* `created_at` - The Unix epoch time stamp in seconds for when the cluster was created.
* `id` - The name of the cluster.
* `kubernetes_network_config` - The Kubernetes network configuration.
  The structure of this block is [described below](#kubernetes_network_config).
* `legacy_cluster_params` - The parameters for fine-tuning the Kubernetes cluster.
  The structure of this block is [described below](#legacy_cluster_params).
* `platform_version` - The platform version for the cluster.
* `status` - The status of the EKS cluster. One of `CLAIMED`, `CREATING`, `DELETED`, `DELETING`, `ERROR`, `MODIFYING`, `PENDING`, `PROVISIONING`, `READY`, `REPAIRING`.
* `version` - The Kubernetes server version for the cluster.
* `vpc_config` - The VPC configuration for the cluster.
  The structure of this block is [described below](#vpc_config).
* `tags` - Map of tags assigned to the cluster.

#### kubernetes_network_config

The `kubernetes_network_config` block has the following structure:

* `ip_family` - The IP family used to assign Kubernetes pod and service addresses.
* `pod_ipv4_cidr` - The CIDR block to assign Kubernetes pod IP addresses from.
* `service_ipv4_cidr` - The CIDR block to assign Kubernetes service IP addresses from.

#### legacy_cluster_params

The `legacy_cluster_params` block has the following structure:

* `cluster_autoscaler_config` – The configuration of the Cluster Autoscaler.
  The structure of this block is [described below](#cluster_autoscaler_config).
* `docker_registry_config` – The configuration of the Docker Registry.
  The structure of this block is [described below](#docker_registry_config).
* `ebs_provider_config` – The configuration of the EBS Provider.
  The structure of this block is [described below](#ebs_provider_config).
* `ingress_config` – The configuration of the Ingress controller.
  The structure of this block is [described below](#ingress_config).
* `master_config` – The configuration of the master node of the cluster.
  The structure of this block is [described below](#master_config).
* `user_data_config` - The configuration of the cluster user data.
  The structure of this block is [described below](#user_data_config).
* `nlb_provider_config` – The configuration of the NLB Provider.
  The structure of this block is [described below](#nlb_provider_config).
* `placement_config` - The placement of the cluster.
  The structure of this block is [described below](#placement_config).

##### cluster_autoscaler_config

The `cluster_autoscaler_config` block has the following structure:

* `cluster_autoscaler_required` - Whether the Cluster Autoscaler is deployed.
* `cluster_autoscaler_user` - K2 service user used by the Cluster Autoscaler.

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
    * _Valid values:_ `cloud-config`,  `x-shellscript`

#### vpc_config

The `vpc_config` block has the following structure:

* `cluster_security_group_id` - The cluster security group that was created by the cloud for the cluster.
* `security_group_ids` - List of security group IDs.
* `subnet_ids` - List of subnet IDs.
* `vpc_id` - The VPC associated with your cluster.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`enabled_cluster_log_types`, `encryption_config`, `endpoint`, `identity`, `role_arn`, `vpc_config.endpoint_private_access`, `vpc_config.endpoint_public_access`, `vpc_config.public_access_cidrs`.
