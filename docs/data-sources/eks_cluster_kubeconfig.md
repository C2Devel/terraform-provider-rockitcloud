---
subcategory: "EKS (Elastic Kubernetes)"
layout: "aws"
page_title: "aws_eks_cluster_kubeconfig"
description: |-
  Returns the kubeconfig for EKS cluster.
---

# Data Source: aws_eks_cluster_kubeconfig

Returns the kubeconfig string for an existing EKS cluster.
The kubeconfig can be used to configure `kubectl` or other Kubernetes tooling.

!> **Warning** The `kubeconfig` attribute is marked as sensitive because it
contains credentials for the cluster.

## Example usage

```terraform
data "aws_eks_cluster_kubeconfig" "example" {
  name = "example"
}
```

## Argument reference

* `name` - (Required, String) The name of the cluster.

## Attribute reference

* `id` - (String) The name of the cluster.
* `kubeconfig` - (String) The kubeconfig for the cluster.
