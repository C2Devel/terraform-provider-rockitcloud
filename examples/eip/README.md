# Elastic IP Example

The EIP example launches the nginx web server and creates a security group.

This example assumes that an SSH key pair was created in the cloud console.

The example takes credentials from a `c2rc.sh` file.
Get the file for your project and place it in this directory before running the example.

Running the example:

```shell
$ terraform init
$ terraform apply -var="key_name=your-key-name"
```

In a few minutes, you can reach the nginx welcome page at the Elastic IP address.

Destroying the example:

```shell
$ terraform destroy -var="key_name=your-key-name"
```

Instead of using `-var`, you can copy `terraform.tfvars.example` to `terraform.tfvars` and use it to specify variable values.
