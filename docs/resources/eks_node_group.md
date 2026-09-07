---
subcategory: "EKS (Elastic Kubernetes)"
layout: "aws"
page_title: "aws_eks_node_group"
description: |-
  Manages an EKS node group.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[eks-node-groups]: https://docs.k2.cloud/en/services/kubernetes/eks_cluster.html#id7
[lifecycle]: https://www.terraform.io/docs/configuration/meta-arguments/lifecycle.html
[timeouts]: https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts

# Resource: aws_eks_node_group

Manages an EKS node group, which can provision and optionally update an autoscaling group of Kubernetes worker nodes compatible with EKS.
For details about EKS node groups, see the [user documentation][eks-node-groups].

~> **Important** A node group can only be deployed in subnets with internet access, so the VPC must have an internet gateway, a NAT gateway and a default route (`0.0.0.0/0`) to the NAT gateway.
It's recommended to specify the route as an explicit dependency via `depends_on`.

## Example Usage

```terraform
data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_vpc" "example" {
  cidr_block = "172.16.0.0/16"

  tags = {
    Name = "tf-vpc"
  }
}

resource "aws_subnet" "example" {
  count = length(data.aws_availability_zones.available.names)

  vpc_id            = aws_vpc.example.id
  cidr_block        = "172.16.${count.index}.0/24"
  availability_zone = data.aws_availability_zones.available.names[count.index]

  tags = {
    Name = "tf-subnet"
  }
}

resource "aws_internet_gateway" "example" {
  vpc_id = aws_vpc.example.id

  tags = {
    Name = "tf-igw"
  }
}

resource "aws_nat_gateway" "example" {
  depends_on = [aws_internet_gateway.example]

  vpc_id = aws_vpc.example.id

  tags = {
    Name = "tf-nat-gw"
  }
}

resource "aws_route" "default_route" {
  route_table_id         = aws_vpc.example.main_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = aws_nat_gateway.example.id
}

resource "aws_eks_cluster" "example" {
  depends_on = [aws_route.default_route]

  name    = "tf-cluster"
  version = "1.30.2"

  vpc_config {
    subnet_ids = aws_subnet.example[*].id
  }
}

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

### Ignoring Changes to Desired Size

~> **Note**
This example uses the VPC, subnets, internet gateway, NAT gateway, default route and cluster defined in the [example above](#example-usage).

You can utilize the generic Terraform resource [lifecycle configuration block][lifecycle] with `ignore_changes` to create an EKS node group with an initial size of running instances, then ignore any changes to that count caused externally.

```terraform
resource "aws_eks_node_group" "example" {
  cluster_name    = aws_eks_cluster.example.name
  instance_types  = ["c5.large"]
  node_group_name = "example"
  subnet_ids      = aws_subnet.example[*].id

  scaling_config {
    # Example: Create EKS node group with 2 instances to start
    desired_size = 2
    max_size     = 2
    min_size     = 1
  }

  # Optional: Allow external changes without Terraform plan difference
  lifecycle {
    ignore_changes = [scaling_config[0].desired_size]
  }
}
```

## Argument Reference

The following arguments are required:

* `cluster_name` - (Required) Name of the EKS cluster.
    * _Value length:_ From 1 to 100 symbols
    * _Constraints:_
        * The value can contain only Latin letters, numbers, hyphens (`-`), and underscores (`_`)
        * The value must start with a Latin letter or a number
* `instance_types` - (Required) List of instance types associated with the EKS node group.
* `scaling_config` - (Required) Configuration block with scaling settings. Detailed below.
* `subnet_ids` - (Required) IDs of EC2 subnets to associate with the EKS node group.

The following arguments are optional:

* `capacity_type` - (Optional) Type of capacity associated with the EKS node group.  Terraform will only perform drift detection if a configuration value is provided.
    * _Valid values:_ `ON_DEMAND`
* `disk_size` - (Optional) Disk size in GiB for worker nodes. Terraform will only perform drift detection if a configuration value is provided.
    * _Default value:_ `20`
* `labels` - (Optional) Key-value map of Kubernetes labels. Only labels that are applied with the EKS API are managed by this argument. Other Kubernetes labels applied to the EKS node group will not be managed.
* `node_group_name` - (Optional) Name of the EKS node group. If omitted, Terraform will assign a random, unique name. Conflicts with `node_group_name_prefix`.
* `node_group_name_prefix` - (Optional) Creates a unique name beginning with the specified prefix. Conflicts with `node_group_name`.
* `remote_access` - (Optional) Configuration block with remote access settings. Detailed below.
* `tags` - (Optional) Map of tags to assign to the node group. If a provider [`default_tags` configuration block][default-tags] is used, tags with matching keys will overwrite those defined at the provider level.
* `taint` - (Optional) The Kubernetes taints to be applied to the nodes in the node group. Maximum of 50 taints per node group. Detailed below.
* `update_config` - (Optional) A block of mutually exclusive arguments that control how many or what percent of nodes can be unavailable during a node group update. Use it to limit disruption while rolling out changes.

### remote_access

* `ec2_ssh_key` - (Optional) EC2 key pair name that provides access for SSH communication with the worker nodes in the EKS node group.

### scaling_config

* `desired_size` - (Required) Desired number of worker nodes.
* `max_size` - (Required) Maximum number of worker nodes.
* `min_size` - (Required) Minimum number of worker nodes. Must be at least `1`.

### taint

* `effect` - (Required) The effect of the taint.
    * _Valid values:_ `NO_SCHEDULE`, `NO_EXECUTE`, `PREFER_NO_SCHEDULE`
* `key` - (Required) The key of the taint.
    * _Value length:_ From 1 to 63 symbols
* `value` - (Optional) The value of the taint.
    * _Value length:_ From 1 to 63 symbols

### update_config

The following arguments are mutually exclusive.

* `max_unavailable` - (Optional) Desired max number of unavailable worker nodes during node group update.
* `max_unavailable_percentage` - (Optional) Desired max percentage of unavailable worker nodes during node group update.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - EKS node group ID.
* `id` - EKS cluster name and EKS node group name separated by a colon (`:`).
* `launch_template` - Computed configuration for the platform-managed launch template. Custom launch templates cannot be configured.
    * `id` - EC2 launch template ID.
    * `name` - Name of the EC2 launch template.
    * `version` - EC2 launch template version number.
* `resources` - List of objects containing information about underlying resources.
    * `autoscaling_groups` - List of objects containing information about autoscaling groups.
        * `name` - Name of the autoscaling group.
* `tags_all` - Map of tags assigned to the node group, including those inherited from the provider [`default_tags` configuration block][default-tags].
* `status` - Status of the EKS node group. One of `CREATING`, `ACTIVE`, `PENDING`, `UPDATING`, `DELETING`, `CREATE_FAILED`, `DELETE_FAILED`, `DEGRADED`.
* `version` - Kubernetes version.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`ami_type`, `force_update_version`, `node_role_arn`, `release_version`, configurable `version`, `remote_access.source_security_group_ids`, `resources.remote_access_security_group_id`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

* `create` - (Default `60 minutes`) How long to wait for the EKS node group to be created.
* `update` - (Default `60 minutes`) How long to wait for the EKS node group to be updated.
* `delete` - (Default `60 minutes`) How long to wait for the EKS node group to be deleted.

~> **Note** K2 node groups inherit the cluster Kubernetes version; version-update arguments are not supported.

## Import

EKS node groups can be imported using the `cluster_name` and `node_group_name` separated by a colon (`:`), e.g.,

```
$ terraform import aws_eks_node_group.my_node_group my_cluster:my_node_group
```
