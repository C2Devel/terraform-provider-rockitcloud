---
subcategory: "EKS (Elastic Kubernetes)"
layout: "aws"
page_title: "aws_eks_cluster"
description: |-
  Manages an EKS cluster.
---

[default-tags]: https://developer.hashicorp.com/terraform/plugin/framework/resources/default-tags
[eks-clusters]: https://docs.k2.cloud/en/services/kubernetes/eks_cluster.html
[ha-clusters]: https://docs.k2.cloud/en/services/kubernetes/overview.html#ha-control-plane
[timeouts]: https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts
[eks-kubeconfig-ds]: ../data-sources/eks_cluster_kubeconfig.md

# Resource: aws_eks_cluster

Manages an EKS cluster. For details about EKS clusters, see the [user documentation][eks-clusters].

## Example usage

### EKS High-Availability Cluster

~> **Note** By default, Terraform creates [high availability clusters][ha-clusters].

```terraform
resource "aws_vpc" "example" {
  cidr_block = "172.16.0.0/16"

  tags = {
    Name = "tf-vpc"
  }
}

resource "aws_subnet" "example" {
  vpc_id            = aws_vpc.example.id
  cidr_block        = cidrsubnet(aws_vpc.example.cidr_block, 4, 1)
  availability_zone = "ru-msk-vol52"

  tags = {
    Name = "tf-subnet"
  }
}

resource "aws_eks_cluster" "example" {
  name    = "tf-cluster-ha"
  version = "1.30.2"

  vpc_config {
    subnet_ids = [aws_subnet.example.id]
  }
}
```

### EKS Cluster with High-Availability Disabled

~> **Note** This example uses the same VPC and subnet as in the [EKS high-availability cluster example](#eks-high-availability-cluster).

```terraform
resource "aws_eks_cluster" "example" {
  name    = "tf-cluster-disabled-ha"
  version = "1.30.2"

  legacy_cluster_params {
    master_config {
      high_availability = false
      instance_type     = "c5.large"
      volume_type       = "gp2"
      volume_size       = 64
    }
  }

  vpc_config {
    subnet_ids = [aws_subnet.example.id]
  }
}
```

### EKS Cluster with extra services

~> **Note** This example uses the same VPC and subnet as in the [EKS High-Availability Cluster example](#eks-high-availability-cluster).

```terraform
resource "aws_eks_cluster" "example" {
  name    = "tf-cluster-extra-services"
  version = "1.30.2"

  legacy_cluster_params {
    docker_registry_config {
      volume_type = "gp2"
      volume_size = 32
    }

    ebs_provider_config {
      ebs_user = "ebs"
    }

    ingress_config {
      instance_type = "c5.large"
      volume_type   = "gp2"
      volume_size   = 32
    }

    master_config {
      high_availability = false
      instance_type     = "c5.large"
      volume_type       = "gp2"
      volume_size       = 64
    }

    nlb_provider_config {
      nlb_user = "nlb"
    }

    placement_config {
      tenancy = "host"
    }
  }

  vpc_config {
    subnet_ids = [aws_subnet.example.id]
  }
}
```

### Get kubeconfig for the cluster

Use the [aws_eks_cluster_kubeconfig data source][eks-kubeconfig-ds] to generate a kubeconfig
for the created cluster and export it for `kubectl`:

```terraform
resource "aws_eks_cluster" "example" {
  name    = "tf-cluster-ha"
  version = "1.30.2"

  vpc_config {
    subnet_ids = [aws_subnet.example.id]
  }
}

data "aws_eks_cluster_kubeconfig" "example" {
  name = aws_eks_cluster.example.name
}

output "kubeconfig" {
  value     = data.aws_eks_cluster_kubeconfig.example.kubeconfig
  sensitive = true
}
```

Then save it locally:

```bash
terraform output -raw kubeconfig > ~/.kube/config
kubectl get nodes
```

## Argument reference

The following arguments are required:

* `name` - (Required, Forces new resource, String) The name of the cluster.
    * _Constraints:_
        * From 1 to 100 characters.
        * The value can contain only Latin letters, numbers, hyphens (`-`), and underscores (`_`).
        * The value must start with a Latin letter or a number.
* `version` - (Required, Forces new resource, String) The Kubernetes server version for the cluster.
* `vpc_config` - (Required, Forces new resource, [Block](#vpc_config)) Configuration block for the VPC associated with your cluster.

The following arguments are optional:

* `enabled_cluster_log_types` - (Optional, Editable, Set of strings) The list of the desired control plane logging to enable.
    * _Valid values:_ `api`, `audit`, `authenticator`, `controllerManager`, `scheduler`
* `encryption_config` - (Optional, Editable, [Block](#encryption_config)) The configuration block for encryption for the cluster.
* `kubernetes_network_config` - (Optional, Editable, [Block](#kubernetes_network_config)) Configuration block with Kubernetes network configuration for the cluster. If removed, Terraform will only perform drift detection if a configuration value is provided.
* `legacy_cluster_params` - (Optional, Editable, [Block](#legacy_cluster_params)) The parameters for fine-tuning the Kubernetes cluster.
* `role_arn` - (Optional, Forces new resource, String) The ARN of the IAM role that provides permissions for the Kubernetes cluster.
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the cluster. If the [`default_tags` configuration block][default-tags] is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.

### kubernetes_network_config

The following arguments are optional:

* `ip_family` - (Optional, Forces new resource, String) The IP family used to assign Kubernetes pod and service addresses.
    * _Valid values:_ `ipv4`
* `service_ipv4_cidr` - (Optional, Forces new resource, String) The CIDR block to assign Kubernetes service IP addresses from. If you don't specify a block, Kubernetes assigns addresses from the 10.96.0.0/12 CIDR block.
    * _Constraints:_
        * Must be within one of the following private IP address blocks: 10.0.0.0/8, 172.16.0.0/12, or 192.168.0.0/16.
        * Must not overlap with any CIDR block assigned to the selected VPC.
        * Must have a prefix length between /12 and /24 (inclusive).

### encryption_config

The following arguments are required:

* `provider` - (Required, Editable, [Block](#provider)) The configuration block for the encryption provider.
* `resources` - (Required, Editable, Set of strings) The resources to encrypt.

#### provider

The following arguments are required:

* `key_arn` - (Required, Editable, String) The ARN of the KMS key.

### legacy_cluster_params

The `legacy_cluster_params` block has the following structure:

* `docker_registry_config` – (Optional, Editable, [Block](#docker_registry_config)) The configuration of the Docker Registry.
* `ebs_provider_config` – (Optional, Editable, [Block](#ebs_provider_config)) The configuration of the EBS Provider.
* `ingress_config` – (Optional, Editable, [Block](#ingress_config)) The configuration of the Ingress controller.
* `master_config` - (Optional, Editable, [Block](#master_config)) The configuration of the master node of the cluster.
* `nlb_provider_config` – (Optional, Editable, [Block](#nlb_provider_config)) The configuration of the NLB Provider.
* `placement_config` - (Optional, Editable, [Block](#placement_config)) The placement of the cluster.
* `user_data_config` - (Optional, Editable, [Block](#user_data_config)) The configuration of the cluster user data.

#### docker_registry_config

The following arguments are required:

* `volume_size` - (Required, Forces new resource, Integer) The size of the Docker Registry volume in GiB.
* `volume_type` - (Required, Forces new resource, String) The type of the Docker Registry volume.

The following arguments are optional:

* `volume_iops` - (Optional, Forces new resource, Integer) The number of read/write operations per second for the Docker Registry volume.
    * _Constraints:_ Required only when `volume_type` is `io2`

#### ebs_provider_config

The following arguments are required:

* `ebs_user` - (Required, Forces new resource, String) The EBS Provider user name.

#### ingress_config

The following arguments are required:

* `instance_type` - (Required, Forces new resource, String) The instance type of the Ingress controller.
* `volume_size` - (Required, Forces new resource, Integer) The size of the Ingress controller volume in GiB.
* `volume_type` - (Required, Forces new resource, String) The type of the Ingress controller volume.

The following arguments are optional:

* `public_ip` - (Optional, Forces new resource, String) The public IP address at which the Ingress controller can be accessed.
* `volume_iops` - (Optional, Forces new resource, Integer) The number of read/write operations per second for the Ingress controller volume.
    ** _Constraints:_ Required only when `volume_type` is `io2`
#### master_config

The following arguments are required:

* `high_availability` - (Required, Forces new resource, Boolean) Indicates whether to deploy a high-availability cluster.
* `instance_type` - (Required, Forces new resource, String) The instance type of the master node.
* `volume_size` - (Required, Forces new resource, Integer) The size of the master node volume in GiB.
* `volume_type` - (Required, Forces new resource, String) The type of the master node volume.

The following arguments are optional:

* `public_ip` - (Optional, Forces new resource, String) The public IP address at which the master node can be accessed.
* `volume_iops` - (Optional, Forces new resource, Integer) The number of read/write operations per second for the master node volume.
    * * _Constraints:_ Required only when `volume_type` is `io2`#### nlb_provider_config

The following arguments are required:

* `nlb_user` - (Required, Forces new resource, String) The NLB Provider user name.

#### placement_config

The following arguments are optional:

* `affinity` - (Optional, Forces new resource, String) The affinity setting for an instance on a dedicated host.
    * _Default value:_ `default`
    * _Valid values:_ `default`, `host`
    * _Constraints:_ The parameter could be set to `host` only if `tenancy` is `host`
* `host_id` - (Optional, Editable, String) The ID of the dedicated host for the instance.
* `tenancy` - (Optional, Forces new resource, String) The tenancy of the instance (if the instance is running in a VPC).
    * _Default value:_ `default`
    * _Valid values:_ `default`, `host`

#### user_data_config

The following arguments are required:

* `user_data` - (Required, Forces new resource, String) User data.
* `user_data_content_type` - (Required, Forces new resource, String) The type of `user_data`.
    * _Valid values:_ `cloud-config`, `x-shellscript`

### vpc_config

The following arguments are required:

* `subnet_ids` - (Required, Forces new resource, Set of strings) The list of subnet IDs.

The following arguments are optional:

* `endpoint_private_access` - (Optional, Editable, Boolean) Indicates whether the endpoint private access is enabled.
* `endpoint_public_access` - (Optional, Editable, Boolean) Indicates whether the endpoint public access is enabled.
* `public_access_cidrs` - (Optional, Editable, Set of strings) The list of CIDR blocks which can access the cluster endpoint.
* `security_group_ids` - (Optional, Forces new resource, Set of strings) The list of security group IDs.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The ARN of the cluster.
* `certificate_authority` - (List) Nested attribute containing `certificate-authority-data` for your cluster.
    * `data` - (String) The base64 encoded certificate data required to communicate with your cluster.
* `created_at` - (String) The Unix epoch time stamp in seconds for when the cluster was created.
* `id` - (String) The name of the cluster.
* `platform_version` - (String) The platform version for the cluster.
* `status` - (String) The status of the EKS cluster.
    * _Valid values:_ `CLAIMED`, `CREATING`, `DELETED`, `DELETING`, `ERROR`, `MODIFYING`, `PENDING`, `PROVISIONING`, `READY`, `REPAIRING`
* `tags_all` - (Map of strings) Key-value pairs assigned to the cluster, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.
* `vpc_config` - (List) Nested list containing VPC configuration for the cluster.
    * `cluster_security_group_id` - (String) The cluster security group that was created for the cluster.
    * `vpc_id` - (String) The VPC associated with your cluster.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`endpoint`, `identity`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

* `create` - (Default `30 minutes`) How long to wait for the EKS cluster to be created.
* `update` - (Default `60 minutes`) How long to wait for the EKS cluster to be updated.
Note that the `update` timeout is used separately for both `version` and `vpc_config` update timeouts.
* `delete` - (Default `15 minutes`) How long to wait for the EKS cluster to be deleted.

## Import

EKS clusters can be imported using the `name`, for example:

```
$ terraform import aws_eks_cluster.my_cluster my_cluster
```
