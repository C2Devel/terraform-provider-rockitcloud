---
subcategory: "EKS (Elastic Kubernetes)"
layout: "aws"
page_title: "aws_eks_node_groups"
description: |-
  Returns a list of EKS node group names associated with an EKS cluster.
---

# Data Source: aws_eks_node_groups

Returns a list of EKS node group names associated with an EKS cluster.

## Example usage

```terraform
data "aws_eks_node_groups" "example" {
  cluster_name = "example"
}

data "aws_eks_node_group" "example" {
  for_each = data.aws_eks_node_groups.example.names

  cluster_name    = "example"
  node_group_name = each.value
}
```

## Argument reference

* `cluster_name` - (Required) The name of the cluster.

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The cluster name.
* `names` - (Set of strings) The set of all node group names in an EKS cluster.
