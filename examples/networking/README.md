# Networking Example

This example creates a network in each of the two regions:

* VPC;
* one or two subnets, depending on the number of availability zones in region;
* route table associated with every subnet;
* security group that allows all traffic within VPC, ICMP, SSH and HTTP(S) from the internet and all outbound traffic.

The network resources are described in the `region` child module.
The example instantiates the module once per region, which demonstrates how to create several copies
of the same resource set with different arguments.

The example takes credentials from a `c2rc.sh` file.
Get the file for your project and place it in this directory before running the example.

Instead of using `-var`, you can copy `terraform.tfvars.example` to `terraform.tfvars` and use it to specify variable values.

Running the example:

```shell
$ terraform init
$ terraform apply
```

Destroying the example:

```shell
$ terraform destroy
```
