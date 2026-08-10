---
subcategory: "EKS (Elastic Kubernetes)"
layout: "aws"
page_title: "aws_eks_cluster_auth"
description: |-
  Provides information about an authentication token to communicate with an EKS cluster.
---

# Data Source: aws_eks_cluster_auth

Provides information about an authentication token to communicate with an EKS cluster.

## Example usage

```terraform
data "aws_eks_cluster_auth" "example" {
  name = "example"
}
```

## Argument reference

* `name` - (Required, String) The name of the cluster.

## Attribute reference

* `id` - (String) The name of the cluster.
* `token` - (String) The authentication token for the EKS cluster.
