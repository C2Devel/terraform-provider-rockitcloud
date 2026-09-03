[Terraform http provider]: https://www.terraform.io/docs/providers/http/index.html

# EKS Example

The EKS (Elastic Kubernetes Service) example launches an EKS cluster and EKS node group with necessary network resources in all availability zones in the region.

The example takes credentials from a `c2rc.sh` file.
Get the file for your project and place it in this directory before running the example.

Running the example:

```shell
$ terraform init
$ terraform apply
```

Get kubeconfig for the created cluster:

```shell
$ terraform output -raw kubeconfig > ~/.kube/config
```

Destroying the example:

```shell
$ terraform destroy
```

This example uses [Terraform http provider] to send a request to https://icanhazip.com/ to determine the local workstation external IP for the security group configuration.
This request is optional and can be replaced by specifying the IP address manually in the [workstation-external-ip.tf](workstation-external-ip.tf).
