---
subcategory: "EKS (Elastic Kubernetes)"
layout: "aws"
page_title: "aws_eks_clusters"
description: |-
  Provides a list of EKS cluster names.
---

# Data Source: aws_eks_clusters

Provides a list of EKS cluster names.

## Example usage

```terraform
data "aws_eks_clusters" "example" {}

data "aws_eks_cluster" "example" {
  for_each = toset(data.aws_eks_clusters.example.names)
  name     = each.value
}
```

## Argument reference

This data source has no arguments.

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The region configured in the provider.
* `names` - (Set of strings) Set of EKS cluster names.
