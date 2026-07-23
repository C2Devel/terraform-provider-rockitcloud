---
subcategory: "EKS (Elastic Kubernetes)"
layout: "aws"
page_title: "aws_eks_node_group"
description: |-
  Manages an EKS node group.
---

[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block
[eks-node-groups]: https://docs.k2.cloud/en/services/kubernetes/eks_cluster.html#id7
[lifecycle]: https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle
[timeouts]: https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts

# Resource: aws_eks_node_group

Manages an EKS node group, which can provision and optionally update an autoscaling group of Kubernetes worker nodes compatible with EKS.
For details about EKS node groups, see the [user documentation][eks-node-groups].

## Example usage

```terraform
resource "aws_eks_node_group" "example" {
  cluster_name    = aws_eks_cluster.example.name
  instance_types  = ["c5.large"]
  node_group_name = "example"
  subnet_ids      = aws_subnet.example[*].id

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  update_config {
    max_unavailable = 2
  }
}
```

### Specific example: Using ignore_changes to preserve external scaling

You can utilize the generic Terraform resource [lifecycle configuration block][lifecycle] with `ignore_changes` to create an EKS node group with an initial size of running instances, then ignore any changes to that count caused externally.

```terraform
resource "aws_eks_node_group" "example" {
  # ... other configurations ...

  scaling_config {
    # Example: Create EKS node group with 2 instances to start
    desired_size = 2

    # ... other configurations ...
  }

  # Optional: Allow external changes without Terraform plan difference
  lifecycle {
    ignore_changes = [scaling_config[0].desired_size]
  }
}
```

### Specific example: creating subnets for a node group

```terraform
data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_subnet" "example" {
  count = 2

  availability_zone = data.aws_availability_zones.available.names[count.index]
  cidr_block        = cidrsubnet(aws_vpc.example.cidr_block, 8, count.index)
  vpc_id            = aws_vpc.example.id

  tags = {
    "kubernetes.io/cluster/${aws_eks_cluster.example.name}" = "shared"
  }
}
```

## Argument reference

The following arguments are required:

* `cluster_name` - (Required, Forces new resource, String) The name of the EKS cluster.
    * _Constraints:_
        * From 1 to 100 characters.
        * The value can contain only Latin letters, numbers, hyphens (`-`), and underscores (`_`).
        * The value must start with a Latin letter or a number.
* `instance_types` - (Required) List of instance types associated with the EKS node group.
* `scaling_config` - (Required, [Block](#scaling_config)) The configuration block with scaling settings.
* `subnet_ids` - (Required, Forces new resource, Set of strings) The IDs of EC2 subnets to associate with the EKS node group.

The following arguments are optional:

* `ami_type` - (Optional, Forces new resource, String) The type of Amazon Machine Image (AMI) associated with the EKS node group.
* `capacity_type` - (Optional, Forces new resource, String) The type of capacity associated with the EKS node group.
    * _Valid values:_ `ON_DEMAND`
* `disk_size` - (Optional, Forces new resource, Integer) The disk size in GiB for worker nodes.
    Terraform will only perform drift detection if a configuration value is provided.
    * _Default value:_ `20`
* `force_update_version` - (Optional, Editable, Boolean) Indicates whether to force a version update of the EKS node group.
* `labels` - (Optional, Editable, Map of strings) Key-value map of Kubernetes labels.
    Only labels that are applied with the EKS API are managed by this argument.
    Other Kubernetes labels applied to the EKS node group will not be managed.
* `launch_template` - (Optional, Editable, [Block](#launch_template)) The configuration block with launch template settings.
* `node_group_name` - (Optional, Forces new resource, String) The name of the EKS node group.
    If omitted, Terraform will assign a random, unique name.
    * _Constraints:_ Conflicts with `node_group_name_prefix`.
* `node_group_name_prefix` - (Optional, Forces new resource, String) The prefix to use for generating a unique name.
    * _Constraints:_ Conflicts with `node_group_name`.
* `node_role_arn` - (Optional, Forces new resource, String) The Amazon Resource Name (ARN) of the IAM role that provides permissions for the EKS node group.
* `release_version` - (Optional, Editable, String) The AMI version of the EKS node group.
* `remote_access` - (Optional, Forces new resource, [Block](#remote_access)) The configuration block with remote access settings.
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the EKS node group.
    If the [`default_tags` configuration block][default-tags] is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.
* `taint` - (Optional, Editable, [Block](#taint)) The Kubernetes taints to apply to the nodes in the node group.
    * _Constraints:_ Maximum of 50 taints per node group.
* `update_config` - (Optional, Editable, [Block](#update_config)) A block of mutually exclusive arguments that control how many or what percent of nodes can be unavailable during a node group update.
    Use it to limit disruption while rolling out changes.
* `version` - (Optional, Editable, String) The Kubernetes version for the EKS node group.

### launch_template

* `id` - (Optional, Forces new resource, String) The ID of the launch template.
    * _Constraints:_ Conflicts with `name`.
* `name` - (Optional, Forces new resource, String) The name of the launch template.
    * _Constraints:_ Conflicts with `id`.
* `version` - (Required, String) The version number of the launch template.
    * _Constraints:_ From 1 to 255 characters.

### remote_access

* `ec2_ssh_key` - (Optional, Forces new resource, String) The EC2 key pair name that provides access for SSH communication with the worker nodes in the EKS node group.
* `source_security_group_ids` - (Optional, Forces new resource, Set of strings) The IDs of the security groups to allow SSH access from the worker nodes.

### scaling_config

* `desired_size` - (Required, Integer) The desired number of worker nodes.
* `max_size` - (Required, Integer) The maximum number of worker nodes.
* `min_size` - (Required, Integer) The minimum number of worker nodes.

### taint

* `effect` - (Required, String) The effect of the taint.
    * _Valid values:_ `NO_EXECUTE`, `NO_SCHEDULE`, `PREFER_NO_SCHEDULE`
* `key` - (Required, String) The key of the taint.
    * _Constraints:_ From 1 to 63 characters.
* `value` - (Optional, String) The value of the taint.
    * _Constraints:_ From 1 to 63 characters.

### update_config

The following arguments are mutually exclusive.

* `max_unavailable` - (Optional, Integer) The desired maximum number of unavailable worker nodes during node group update.
    * _Constraints:_ From 1 to 100.
* `max_unavailable_percentage` - (Optional, Integer) The desired maximum percentage of unavailable worker nodes during node group update.
    * _Constraints:_ From 1 to 100.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the EKS node group.
* `id` - (String) The EKS cluster name and EKS node group name separated by a colon (`:`).
* `launch_template` - (List of objects) The launch template configuration.
    * `id` - (String) The ID of the launch template.
    * `name` - (String) The name of the launch template.
    * `version` - (String) The version number of the launch template.
* `resources` - (List of objects) Information about underlying resources.
    * `autoscaling_groups` - (List of objects) Information about autoscaling groups.
        * `name` - (String) The name of the autoscaling group.
    * `remote_access_security_group_id` - (String) The ID of the security group for remote access.
* `status` - (String) The status of the EKS node group.
    * _Valid values:_ `ACTIVE`, `CREATE_FAILED`, `CREATING`, `DEGRADED`, `DELETE_FAILED`, `DELETING`, `PENDING`, `UPDATING`.
* `tags_all` - (Map of strings) Key-value pairs assigned to the EKS node group, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.
* `version` - (String) The Kubernetes version.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

* `create` - (Default `60 minutes`) How long to wait for the EKS node group to be created.
* `update` - (Default `60 minutes`) How long to wait for the EKS node group to be updated.
* `delete` - (Default `60 minutes`) How long to wait for the EKS node group to be deleted.

~> **Note** The `update` timeout is used separately for both configuration and version update operations.

## Import

EKS node groups can be imported using the `cluster_name` and `node_group_name` separated by a colon (`:`), for example:

```
$ terraform import aws_eks_node_group.my_node_group my_cluster:my_node_group
```
