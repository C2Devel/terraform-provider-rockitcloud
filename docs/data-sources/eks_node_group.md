---
subcategory: "EKS (Elastic Kubernetes)"
layout: "aws"
page_title: "aws_eks_node_group"
description: |-
  Provides information about an EKS node group.
---

# Data Source: aws_eks_node_group

Provides information about an EKS node group.

## Example usage

```terraform
data "aws_eks_node_group" "example" {
  cluster_name    = "example"
  node_group_name = "example"
}
```

## Argument reference

* `cluster_name` - (Required) The name of the cluster.
* `node_group_name` - (Required) The name of the node group.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the EKS node group.
* `disk_size` - (Integer) The volume size in GiB for worker nodes.
* `id` - (String) The EKS cluster name and EKS node group name separated by a colon (`:`).
* `instance_types` - (List of strings) The set of instance types associated with the EKS node group.
* `labels` - (Map of strings) Key-value map of Kubernetes labels. Only labels that are applied with the EKS API are managed by this argument. Other Kubernetes labels applied to the EKS node group will not be managed.
* `remote_access` - (List of objects) A configuration block with remote access settings.
    * `ec2_ssh_key` - (String) The name of the key pair that provides access for SSH communication with the worker nodes in the EKS node group.
* `resources` - (List of objects) A list of objects containing information about underlying resources.
    * `autoscaling_groups` - (List of objects) The list of objects containing information about autoscaling groups.
        * `name` - (String) The name of the autoscaling group.
* `scaling_config` - ([Block](#scaling_config)) A configuration block with scaling settings.
* `status` - (String) The status of the EKS node group.
    * _Valid values:_ `ACTIVE`, `CREATE_FAILED`, `CREATING`, `DEGRADED`, `DELETE_FAILED`, `DELETING`, `PENDING`, `UPDATING`
* `subnet_ids` - (Set of strings) Identifiers of EC2 subnets to associate with the EKS node group.
* `tags` - (Map of strings) Key-value pairs assigned to the node group.
* `taints` - ([Block](#taints)) The list of objects containing information about taints applied to the nodes in the EKS node group.
* `version` - (String) The Kubernetes version.

#### scaling_config

The `scaling_config` block has the following structure:

* `desired_size` - (Integer) The desired number of worker nodes.
* `max_size` - (Integer) The maximum number of worker nodes.
* `min_size` - (Integer) The minimum number of worker nodes.

#### taints

The `taints` block has the following structure:

* `key` - (String) The key of the taint.
* `value` - (String) The value of the taint.
* `effect` - (String) The effect of the taint.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`ami_type`, `node_role_arn`, `release_version`, `remote_access.source_security_group_ids`, `resources.remote_access_security_group_id`.
