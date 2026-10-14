---
subcategory: "PaaS"
layout: "aws"
page_title: "aws_paas_service"
description: |-
  Manages a PaaS service.
---

[doc-effective_cache_size]: https://postgresqlco.nf/doc/en/param/effective_cache_size/
[doc-innodb_flush_log_at_trx_commit]: https://dev.mysql.com/doc/refman/5.7/en/innodb-parameters.html#sysvar_innodb_flush_log_at_trx_commit
[doc-innodb_strict_mode]: https://dev.mysql.com/doc/refman/5.7/en/innodb-parameters.html#sysvar_innodb_strict_mode
[doc-mariadb-charset-collate]: https://mariadb.com/kb/en/supported-character-sets-and-collations/
[doc-mysql-charset-collate]: https://dev.mysql.com/doc/refman/8.0/en/charset-charsets.html
[doc-pxc_strict_mode]: https://docs.percona.com/percona-xtradb-cluster/5.7/features/pxc-strict-mode.html
[doc-shared_buffers]: https://postgresqlco.nf/doc/en/param/shared_buffers/
[doc-transaction_isolation]: https://dev.mysql.com/doc/refman/5.7/en/server-system-variables.html#sysvar_transaction_isolation

[elk-version]: https://docs.k2.cloud/en/api/paas/parameters/elk.html#version
[elasticsearch-version]: https://docs.k2.cloud/en/api/paas/parameters/elasticsearch.html#version
[mongodb-version]: https://docs.k2.cloud/en/api/paas/parameters/mongodb.html#version
[mysql-version]: https://docs.k2.cloud/en/api/paas/parameters/mysql.html#version
[paas]: https://docs.k2.cloud/en/services/paas/index.html
[pgsql-version]: https://docs.k2.cloud/en/api/paas/parameters/pgsql.html#version
[kafka-version]: https://docs.k2.cloud/en/api/paas/parameters/kafka.html#version
[rabbitmq-version]: https://docs.k2.cloud/en/api/paas/parameters/rabbitmq.html#version
[redis-version]: https://docs.k2.cloud/en/api/paas/parameters/redis.html#version
[technical support]: https://support.k2int.ru/app/#/project/CS
[timeouts]: https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts

[ClickHouse]: #clickhouse-argument-reference
[ELK]: #elk-argument-reference
[Elasticsearch]: #elasticsearch-argument-reference
[GitLab]: #gitlab-argument-reference
[Kafka]: #kafka-argument-reference
[Memcached]: #memcached-argument-reference
[MongoDB]: #mongodb-argument-reference
[MySQL]: #mysql-argument-reference
[PostgreSQL]: #postgresql-argument-reference
[Prometheus]: #prometheus-argument-reference
[RabbitMQ]: #rabbitmq-argument-reference
[Redis]: #redis-argument-reference

# Resource: aws_paas_service

Manages a PaaS service. For details about PaaS, see the [user documentation][paas].

~> **Important** A PaaS service requires the VPC to have an internet gateway, a NAT gateway and a default route (`0.0.0.0/0`) to the NAT gateway; otherwise creating the service fails.
It's recommended to specify the route as an explicit dependency via `depends_on`.

## Example usage

### Elasticsearch service

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
  vpc_id            = aws_vpc.example.id
  cidr_block        = cidrsubnet(aws_vpc.example.cidr_block, 4, 1)
  availability_zone = data.aws_availability_zones.available.names[0]

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

resource "aws_paas_service" "elasticsearch" {
  depends_on = [aws_route.default_route]

  name          = "tf-service"
  instance_type = "c5.large"

  root_volume {
    type = "gp2"
    size = 32
  }

  data_volume {
    type = "gp2"
    size = 32
  }

  delete_interfaces_on_destroy = true
  security_group_ids           = [aws_vpc.example.default_security_group_id]
  subnet_ids                   = [aws_subnet.example.id]

  ssh_key_name = "<name>"

  elasticsearch {
    version = "8.12"
    kibana  = true
  }
}
```

### ClickHouse service

~> **Note** This example uses the VPC, subnet, internet gateway, NAT gateway and default route defined in the [Elasticsearch service example](#elasticsearch-service).

```terraform
resource "aws_paas_service" "clickhouse" {
  depends_on = [aws_route.default_route]

  name          = "tf-service"
  instance_type = "c5.large"

  additional_roles = ["coordinator"]

  root_volume {
    type = "gp2"
    size = 32
  }

  data_volume {
    type = "gp2"
    size = 32
  }

  delete_interfaces_on_destroy = true
  security_group_ids           = [aws_vpc.example.default_security_group_id]
  subnet_ids                   = [aws_subnet.example.id]

  ssh_key_name = "<name>"

  clickhouse {
    version = "25.3"

    database {
      name = "test_db"
    }

    user {
      name     = "test_user"
      password = "********"
      grants   = ["GRANT SELECT ON test_db.*"]

      networks {
        ip = ["172.16.0.0/16"]
      }
    }
  }
}
```

### ClickHouse cluster with dedicated keepers

~> **Note** This example uses the VPC, subnet, internet gateway, NAT gateway and default route defined in the [Elasticsearch service example](#elasticsearch-service).

```terraform
resource "aws_paas_service" "clickhouse" {
  depends_on = [aws_route.default_route]

  name          = "tf-service"
  instance_type = "c5.large"

  high_availability = true
  instance_qty      = 6

  root_volume {
    type = "gp2"
    size = 32
  }

  data_volume {
    type = "gp2"
    size = 32
  }

  coordinator {
    instance_type    = "c5.large"
    root_volume_type = "gp2"
    root_volume_size = 32
  }

  delete_interfaces_on_destroy = true
  security_group_ids           = [aws_vpc.example.default_security_group_id]
  subnet_ids                   = [aws_subnet.example.id]

  ssh_key_name = "<name>"

  clickhouse {
    version = "25.3"
  }
}
```

### GitLab service

GitLab requires an instance with at least 2 vCPUs and 8 GiB of memory.

~> **Note** This example uses the VPC, subnet, internet gateway, NAT gateway and default route defined in the [Elasticsearch service example](#elasticsearch-service).

```terraform
resource "aws_paas_service" "gitlab" {
  depends_on = [aws_route.default_route]

  name          = "tf-service"
  instance_type = "m5.large"

  root_volume {
    type = "gp2"
    size = 32
  }

  data_volume {
    type = "gp2"
    size = 64
  }

  delete_interfaces_on_destroy = true
  security_group_ids           = [aws_vpc.example.default_security_group_id]
  subnet_ids                   = [aws_subnet.example.id]

  ssh_key_name = "<name>"

  gitlab {
    password = "********"

    registry {
      storage = "filesystem"
    }

    smtp {
      address        = "smtp.gitlab.invalid"
      port           = 587
      username       = "gitlab@gitlab.invalid"
      password       = "********"
      authentication = "login"
      email_from     = "gitlab@gitlab.invalid"
    }
  }
}
```

### ELK service

~> **Note** This example uses the VPC, subnet, internet gateway, NAT gateway and default route defined in the [Elasticsearch service example](#elasticsearch-service).

```terraform
resource "aws_paas_service" "elk" {
  depends_on = [aws_route.default_route]

  name          = "tf-service"
  instance_type = "m5.large"

  root_volume {
    type = "gp2"
    size = 32
  }

  data_volume {
    type = "gp2"
    size = 32
  }

  delete_interfaces_on_destroy = true
  security_group_ids           = [aws_vpc.example.default_security_group_id]
  subnet_ids                   = [aws_subnet.example.id]

  ssh_key_name = "<name>"

  elk {
    version         = "8.17"
    password        = "********"
    allow_anonymous = false
  }
}
```

### Memcached service with enabled monitoring

~> **Note** This example uses the VPC, subnet, internet gateway, NAT gateway and default route defined in the [Elasticsearch service example](#elasticsearch-service).

```terraform
resource "aws_paas_service" "memcached" {
  depends_on = [aws_route.default_route]

  name          = "tf-service"
  instance_type = "c5.large"

  root_volume {
    type = "gp2"
    size = 32
  }

  data_volume {
    type = "gp2"
    size = 32
  }

  delete_interfaces_on_destroy = true
  security_group_ids           = [aws_vpc.example.default_security_group_id]
  subnet_ids                   = [aws_subnet.example.id]

  ssh_key_name = "<name>"

  memcached {
    monitoring {
      monitor_by = "fm-cluster-12345678"
      monitoring_labels = {
        key1 = "value1"
        key3 = "value3"
      }
    }
  }
}
```

### MongoDB service

~> **Note** This example uses the VPC, subnet, internet gateway, NAT gateway and default route defined in the [Elasticsearch service example](#elasticsearch-service).

```terraform
resource "aws_paas_service" "mongodb" {
  depends_on = [aws_route.default_route]

  name          = "tf-service"
  instance_type = "c5.large"

  root_volume {
    type = "gp2"
    size = 32
  }

  data_volume {
    type = "gp2"
    size = 32
  }

  delete_interfaces_on_destroy = true
  security_group_ids           = [aws_vpc.example.default_security_group_id]
  subnet_ids                   = [aws_subnet.example.id]

  ssh_key_name = "<name>"

  mongodb {
    version = "5.0"

    journal_commit_interval = 301
    maxconns                = 16
    profile                 = "all"
    slowms                  = 3600001

    quiet          = false
    verbositylevel = "vvvv"

    user {
      name     = "user1"
      password = "********"
    }

    database {
      name = "test_db1"

      user {
        name  = "user1"
        roles = ["readWrite", "dbAdmin"]
      }
    }
  }
}
```

### MySQL service

~> **Note** This example uses the VPC, subnet, internet gateway, NAT gateway and default route defined in the [Elasticsearch service example](#elasticsearch-service).

```terraform
resource "aws_paas_service" "mysql" {
  depends_on = [aws_route.default_route]

  name          = "tf-service"
  instance_type = "c5.large"

  root_volume {
    type = "gp2"
    size = 32
  }

  data_volume {
    type = "gp2"
    size = 32
  }

  delete_interfaces_on_destroy = true
  security_group_ids           = [aws_vpc.example.default_security_group_id]
  subnet_ids                   = [aws_subnet.example.id]

  ssh_key_name = "<name>"

  mysql {
    vendor  = "mariadb"
    version = "10.11"

    user {
      name     = "user1"
      host     = "127.0.0.1"
      password = "********"
    }

    user {
      name     = "user2"
      host     = "127.0.0.1"
      password = "********"
    }

    database {
      backup_enabled = false
      name           = "test_db1"

      user {
        name       = "user1"
        privileges = ["INSERT"]
        options    = ["GRANT"]
      }

      user {
        name = "user2"
      }
    }
  }
}
```

### PostgreSQL service with arbitrator

```terraform
data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_vpc" "example" {
  cidr_block = "172.33.0.0/16"

  tags = {
    Name = "tf-vpc"
  }
}

resource "aws_subnet" "example" {
  # A service is deployed in either one or three availability zones
  count = length(data.aws_availability_zones.available.names) >= 3 ? 3 : 1

  vpc_id            = aws_vpc.example.id
  cidr_block        = cidrsubnet(aws_vpc.example.cidr_block, 8, count.index)
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

resource "aws_s3_bucket" "example" {
  bucket = "tf-paas-backup"
}

resource "aws_paas_service" "pgsql" {
  depends_on = [aws_route.default_route]

  name = "tf-service"

  arbitrator_required = true
  high_availability   = true

  instance_type = "c5.large"

  root_volume {
    type = "gp2"
    size = 32
  }

  data_volume {
    type = "gp2"
    size = 32
  }

  delete_interfaces_on_destroy = true
  security_group_ids           = [aws_vpc.example.default_security_group_id]
  subnet_ids                   = aws_subnet.example[*].id

  ssh_key_name = "<name>"

  backup_settings {
    enabled            = true
    expiration_days    = 5
    notification_email = "example@mail.com"
    start_time         = "15:10"
    bucket_name        = aws_s3_bucket.example.id
    user_login         = "user@company"
  }

  pgsql {
    version = "16"

    autovacuum_analyze_scale_factor = 0.3
    work_mem                        = 4 * 1024 * 1024
    maintenance_work_mem            = 1024 * 1024
    replication_mode                = "synchronous"

    user {
      name     = "user1"
      password = "********"
    }

    user {
      name     = "user2"
      password = "********"
    }

    database {
      name           = "test_db1"
      owner          = "user1"
      backup_enabled = true
      extensions     = ["bloom", "dict_int"]
      user {
        name = "user2"
      }
    }

    options = {
      logDestination = "csvlog"
    }
  }
}
```

### Kafka service

~> **Note** This example uses the VPC, subnet, internet gateway, NAT gateway and default route defined in the [Elasticsearch service example](#elasticsearch-service).

~> **Note** Kafka always requires a coordinator role. For HA clusters, use either a dedicated `coordinator` block (shown below) or `additional_roles = ["coordinator"]` for combined broker+coordinator nodes.

```terraform
resource "aws_paas_service" "kafka" {
  depends_on = [aws_route.default_route]

  name              = "tf-service"
  high_availability = true
  instance_type     = "c5.large"

  root_volume {
    type = "gp2"
    size = 32
  }

  data_volume {
    type = "gp2"
    size = 64
  }

  delete_interfaces_on_destroy = true
  security_group_ids           = [aws_vpc.example.default_security_group_id]
  subnet_ids                   = [aws_subnet.example.id]

  ssh_key_name = "<name>"

  coordinator {
    instance_type    = "c5.large"
    root_volume_type = "gp2"
    root_volume_size = 32
    data_volume_type = "gp2"
    data_volume_size = 64
  }

  kafka {
    version = "3.7.0"
  }
}
```

To create topics in this service, use the [aws_paas_kafka_topic](paas_kafka_topic.md) resource.

### Prometheus service

~> **Note** This example uses the VPC, subnet, internet gateway, NAT gateway and default route defined in the [Elasticsearch service example](#elasticsearch-service).

~> **Note** K2 Cloud supports Prometheus only as a single-node service, so `high_availability` must be `false`.

```terraform
resource "aws_paas_service" "prometheus" {
  depends_on = [aws_route.default_route]

  name              = "tf-service"
  high_availability = false
  instance_type     = "c5.large"

  root_volume {
    type = "gp2"
    size = 32
  }

  data_volume {
    type = "gp2"
    size = 32
  }

  delete_interfaces_on_destroy = true
  security_group_ids           = [aws_vpc.example.default_security_group_id]
  subnet_ids                   = [aws_subnet.example.id]

  ssh_key_name = "<name>"

  prometheus {
    class                 = "monitoring"
    remote_write_receiver = true
  }
}
```

The `remote_write_receiver` argument is supported for Prometheus environment versions `paas_v4_0` and later.

To configure alert delivery and metric collection, use the
[aws_paas_prometheus_notification_channel](paas_prometheus_notification_channel.md),
[aws_paas_prometheus_route](paas_prometheus_route.md), and
[aws_paas_prometheus_scrape_job](paas_prometheus_scrape_job.md) resources.

A complete runnable configuration is available in the
[Prometheus PaaS example](https://github.com/C2Devel/terraform-provider-rockitcloud/tree/develop/examples/paas/prometheus).

### RabbitMQ service

~> **Note** This example uses the VPC, subnet, internet gateway, NAT gateway and default route defined in the [Elasticsearch service example](#elasticsearch-service).

```terraform
resource "aws_paas_service" "rabbitmq" {
  depends_on = [aws_route.default_route]

  name          = "tf-service"
  instance_type = "c5.large"

  root_volume {
    type = "gp2"
    size = 32
  }

  data_volume {
    type = "gp2"
    size = 40
  }

  delete_interfaces_on_destroy = true
  security_group_ids           = [aws_vpc.example.default_security_group_id]
  subnet_ids                   = [aws_subnet.example.id]

  ssh_key_name = "<name>"

  rabbitmq {
    version  = "3.10"
    password = "********"
  }
}
```

### Redis service with logging enabled

~> **Note** This example uses the VPC, subnet, internet gateway, NAT gateway and default route defined in the [Elasticsearch service example](#elasticsearch-service).

```terraform
resource "aws_paas_service" "redis" {
  depends_on = [aws_route.default_route]

  name          = "tf-service"
  instance_type = "c5.large"

  root_volume {
    type = "gp2"
    size = 32
  }

  data_volume {
    type = "gp2"
    size = 32
  }

  delete_interfaces_on_destroy = true
  security_group_ids           = [aws_vpc.example.default_security_group_id]
  subnet_ids                   = [aws_subnet.example.id]

  ssh_key_name = "<name>"

  redis {
    class   = "database"
    version = "7.2"

    password = "********"

    persistence_rdb = false
    persistence_aof = true

    databases     = 1
    timeout       = 50
    tcp_backlog   = 300
    tcp_keepalive = 600

    logging {
      log_to       = "fm-cluster-87654321"
      logging_tags = ["tag1", "tag2", "tag3"]
    }
  }
}
```

## Argument reference

~> **Note** Arguments are not editable (changes lead to a new resource) except for the blocks with service parameters and `backup_settings`.

The following arguments are required:

* `instance_type` - (Required) The instance type.
* `name` - (Required) The service name. The value must start and end with a Latin letter or number and
  can only contain lowercase Latin letters, numbers, periods (.) and hyphens (-).
* `root_volume` - (Required) The root volume parameters for the service. The structure of this block is [described below](#root_volume).
* `security_group_ids` - (Required) List of security group IDs.
* `ssh_key_name` - (Required) The name of the SSH key for accessing instances.

The following arguments are optional:

* `arbitrator_required` - (Optional) Indicates whether to create a cluster with an arbitrator.
    * _Default value:_ `false`
    * _Constraints:_ The parameter can be set to `true` only if `high_availability` is `true`
      The parameter is supported only for [ELK], [Elasticsearch], [MongoDB], [MySQL] and [PostgreSQL] services.
* `additional_roles` - (Optional, ForceNew) Additional roles for the main service nodes (broker nodes for [Kafka]).
    * _Valid values:_ `coordinator`
    * _Constraints:_ Supported only for [ClickHouse] and [Kafka]. Cannot be used together with `coordinator`.
      For [ClickHouse], set the value to `["coordinator"]` to create a standalone service or a cluster of combined nodes.
* `backup_settings` - (Optional) The backup settings for the service. The structure of this block is [described below](#backup_settings).
  The parameter is supported only for [MySQL] and [PostgreSQL] services.
* `coordinator` - (Optional, ForceNew) Dedicated coordinator node parameters. The structure of this block is [described below](#coordinator).
    * _Constraints:_ Supported only for [ClickHouse] and [Kafka]. Cannot be used together with `additional_roles`.
      For [ClickHouse], the block is used to create a cluster with dedicated keepers.
* `data_volume` - (Optional) The data volume parameters for the service. The structure of this block is [described below](#data_volume).
  The provider schema requires this parameter for [ClickHouse], [ELK], [Elasticsearch], [GitLab], [Kafka], [Memcached], [MongoDB],
  [MySQL], [PostgreSQL], [Prometheus], [RabbitMQ] and [Redis] services.
* `delete_interfaces_on_destroy` - (Optional) Indicates whether to delete the instance network interfaces when the service is destroyed.
    * _Default value:_ `false`
* `high_availability` - (Optional) Indicates whether to create a high availability service.
  The parameter is supported only for [ClickHouse], [ELK], [Elasticsearch], [Kafka], [MongoDB], [MySQL], [PostgreSQL], [RabbitMQ]
  and [Redis] services. For [Prometheus], the value must remain `false`.
    * _Default value:_ `false`
* `instance_qty` - (Optional, ForceNew) The total number of service instances including coordinator instances.
  The parameter can be set only if `high_availability` is `true`.
  If the parameter is not specified, the default number of instances for the service type is used.
* `network_interface_ids` - (Optional) List of network interface IDs.
    * _Constraints:_ Required if `subnet_ids` is not specified
* `subnet_ids` - (Optional) List of subnet IDs.
    * _Constraints:_ Required if `network_interface_ids` is not specified
* `user_data` - (Optional) User data.
    * _Constraints:_ Required if `user_data_content_type` is specified
* `user_data_content_type` - (Optional) The type of `user_data`.
    * _Valid values:_ `cloud-config`, `x-shellscript`
    * _Constraints:_ Required if `user_data` is specified

One of the following blocks with service parameters must be specified:

* `clickhouse` - ClickHouse parameters. The structure of this block is [described below](#clickhouse-argument-reference).
* `elk` - ELK parameters. The structure of this block is [described below](#elk-argument-reference).
* `elasticsearch` - Elasticsearch parameters. The structure of this block is [described below](#elasticsearch-argument-reference).
* `gitlab` - GitLab parameters. The structure of this block is [described below](#gitlab-argument-reference).
* `kafka` - Kafka parameters. The structure of this block is [described below](#kafka-argument-reference).
* `memcached` - Memcached parameters. The structure of this block is [described below](#memcached-argument-reference).
* `mongodb` - MongoDB parameters. The structure of this block is [described below](#mongodb-argument-reference).
* `mysql` - MySQL parameters. The structure of this block is [described below](#mysql-argument-reference).
* `pgsql` - PostgreSQL parameters. The structure of this block is [described below](#postgresql-argument-reference).
* `prometheus` - Prometheus parameters. The structure of this block is [described below](#prometheus-argument-reference).
* `rabbitmq` - RabbitMQ parameters. The structure of this block is [described below](#rabbitmq-argument-reference).
* `redis` - Redis parameters. The structure of this block is [described below](#redis-argument-reference).

### backup_settings

~> **Note** All the parameters in the `backup_settings` block are editable.

The `backup_settings` block has the following structure:

* `bucket_name` - (Optional) The name of the bucket in object storage where the service backup is saved.
  The parameter must be set if `enabled` is `true`.
* `enabled` - (Optional) Indicates whether backup is enabled for the service.
    * _Default value:_ `false`
* `expiration_days` - (Optional) The backup retention period in days.
    * _Valid values:_ From 1 to 3650
* `notification_email` - (Optional) The email address to which a notification that backup was created is sent.
* `start_time` - (Optional) The time when the daily backup process starts. It is set as a string in the HH:MM format Moscow time.
  The parameter must be set if `enabled` is `true`.
* `user_login` - (Optional) The login of a user with write permissions to the bucket in object storage (e.g. `user@company`).
  The parameter must be set if `enabled` is `true`.

### data_volume

The `data_volume` block has the following structure:

* `iops` - (Optional) The number of read/write operations per second for the data volume.
  The parameter must be set if `type` is `io2`.
* `size` - (Optional) The size of the data volume in GiB.
    * _Default value:_ `32`
    * _Constraints:_ For [Kafka], the minimum value is `64`
* `type` - (Optional) The type of the data volume.
    * _Default value:_ `st2`

### coordinator

The `coordinator` block has the following structure:

* `data_volume_iops` - (Optional) The number of read/write operations per second for the coordinator data volume.
  The parameter must be set if `data_volume_type` is `io2`.
* `data_volume_size` - (Optional) The size of the coordinator data volume in GiB.
    * _Constraints:_ For [Kafka], set this to at least `64`
* `data_volume_type` - (Optional) The type of the coordinator data volume.
* `instance_type` - (Required) The instance type for coordinator nodes.
* `root_volume_iops` - (Optional) The number of read/write operations per second for the coordinator root volume.
  The parameter must be set if `root_volume_type` is `io2`.
* `root_volume_size` - (Required) The size of the coordinator root volume in GiB.
* `root_volume_type` - (Required) The type of the coordinator root volume.

### root_volume

The `root_volume` block has the following structure:

* `iops` - (Optional) The number of read/write operations per second for the root volume.
  The parameter must be set if `type` is `io2`.
* `size` - (Optional) The size of the root volume in GiB.
    * _Default value:_ `32`
* `type` - (Optional) The type of the root volume.
    * _Default value:_ `st2`

## ClickHouse argument reference

In addition to the common arguments for all services [described above](#argument-reference),
the `clickhouse` block can contain the following arguments:

* `class` - (Optional) The service class.
    * _Valid values:_ `database`
    * _Default value:_ `database`
* `database` - (Optional, Editable) List of ClickHouse databases with parameters. The structure of this block is [described below](#clickhouse-database).
* `logging` - (Optional, Editable) The logging settings for the service. The structure of this block is [described below](#logging).
* `monitoring` - (Optional, Editable) The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `options` - (Optional, Editable) Map containing advanced ClickHouse server settings. Values are strings.
* `user` - (Optional, Editable) List of ClickHouse users with parameters. The structure of this block is [described below](#clickhouse-user).
* `version` - (Required) The version to install. Only LTS versions are available.
  The list of supported versions is available in the [user documentation][paas].

ClickHouse profiles, quotas and shards are managed by the
[`aws_paas_clickhouse_profile`](paas_clickhouse_profile.md), [`aws_paas_clickhouse_quota`](paas_clickhouse_quota.md)
and [`aws_paas_clickhouse_shard`](paas_clickhouse_shard.md) resources.

### ClickHouse database

~> **Note** All the parameters in the `database` block are editable.

The `database` block has the following structure:

* `name` - (Required) The database name. The value must start with a Latin letter, contain at most 63 characters,
  and can only contain Latin letters, numbers and underscores (_).
  The value must not be `default`, `readonly`, `system` and `information_schema`.

### ClickHouse user

~> **Note** All the parameters in the `user` block are editable.

The `user` block has the following structure:

* `grants` - (Optional) List of `GRANT` statements without grantees (e.g. `GRANT SELECT ON db.*`).
* `name` - (Required) The ClickHouse user name. The value must start with a Latin letter, contain at most 63 characters,
  and can only contain Latin letters, numbers and underscores (_).
  The value must not be `default` and `cluster_user`.
* `networks` - (Optional) Networks from which the user can connect. The structure of this block is [described below](#clickhouse-user-networks).
* `password` - (Optional) The ClickHouse user password.
* `password_type` - (Optional) The password encoding.
    * _Valid values:_ `double_sha1`
    * _Default value:_ `double_sha1`
* `profile` - (Optional) The name of the settings profile assigned to the user.
  It can be a built-in profile (`default`, `readonly`) or an existing profile created with the `aws_paas_clickhouse_profile` resource.
    * _Default value:_ `default`
* `quota` - (Optional) The name of the quota assigned to the user.
  It can be a built-in quota (`default`) or an existing quota created with the `aws_paas_clickhouse_quota` resource.
    * _Default value:_ `default`

~> **Note** An omitted or empty `grants` list selects the backend defaults, which can include broad privileges
for [ClickHouse XML users](https://github.com/ClickHouse/ClickHouse/blob/25.3/src/Access/UsersConfigAccessStorage.cpp).
To give a user no table privileges, use an explicit non-empty list such as `grants = ["GRANT USAGE ON *.*"]`.
`USAGE` is an alias of the [`NONE` privilege](https://clickhouse.com/docs/reference/statements/grant#none).

~> **Note** Custom profile and quota resources depend on their parent service through `service_id`.
First apply the service with built-in user assignments and create the child resources.
Then set `profile` and `quota` to their existing names as literal strings and apply again.
Referencing child resource attributes from the same parent service creates a Terraform dependency cycle.

~> **Note** K2 Cloud rejects deletion of a profile or quota while it is assigned to a user.
Before removing these child resources or running `terraform destroy`, change all affected user assignments
to existing alternatives such as `default` and apply that change first.

#### ClickHouse user networks

The `networks` block has the following structure:

* `host` - (Optional) List of host names from which the user can connect.
* `host_regexp` - (Optional) List of regular expressions for host names from which the user can connect.
* `ip` - (Optional) List of IP addresses or subnets from which the user can connect.

## Elasticsearch argument reference

In addition to the common arguments for all services [described above](#argument-reference),
the `elasticsearch` block can contain the following arguments as described below.

The following arguments are required:

* `version` - (Required) The version to install. The list of supported versions is available in the [user documentation][elasticsearch-version].

The following arguments are optional:

* `allow_anonymous` - (Optional) Indicates whether anonymous access is enabled.
  The parameter can be set only if `kibana` is `true` and `password` is specified.
* `anonymous_role` - (Optional) The role for anonymous access.
    * _Valid values:_ `viewer`, `editor`
  The parameter can be set only if `allow_anonymous` is `true`.
* `class` - (Optional) The service class.
    * _Valid values:_ `search`
    * _Default value:_ `search`
* `kibana` - (Optional) Indicates whether the Kibana deployment is enabled.
    * _Default value:_ `false`
* `logging` - (Optional, Editable) The logging settings for the service. The structure of this block is [described below](#logging).
* `monitoring` - (Optional, Editable) The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `options` - (Optional, Editable) Map containing other Elasticsearch parameters.
  Parameter names must be in camelCase. Values are strings.

~> **Note** If a parameter name includes a dot, it cannot be passed in `options`.
If you need to use such a parameter, contact [technical support].

* `password` - (Optional) The Elasticsearch user password.
  The value must not contain `-`, `!`, `:`, `;`, `%`, `'`, `"`, `` ` `` and `\`.

## ELK argument reference

In addition to the common arguments for all services [described above](#argument-reference),
the `elk` block can contain the following arguments:

* `allow_anonymous` - (Optional) Indicates whether anonymous access to Kibana is enabled.
  When `password` is specified, the API defaults this value to `false`.
* `anonymous_role` - (Optional) List containing at most one role for anonymous access.
    * _Valid values:_ `viewer`, `editor`
    * _Default value:_ `viewer`
    * _Constraints:_ The parameter is applicable only if `allow_anonymous` is `true`.
  Terraform uses the documented list-shaped configuration but sends its one
  element as the scalar value required by the live API.
* `class` - (Optional) The service class.
    * _Valid values:_ `logging`
    * _Default value:_ `logging`
* `monitoring` - (Optional, Editable) The monitoring settings for the service.
  The structure of this block is [described below](#monitoring).
* `options` - (Optional, Forces new resource) Map containing other ELK parameters. Values are strings.
* `password` - (Optional) The Elasticsearch user password. The value must be 8 to 128 characters long
  and must not contain `-`, `!`, `:`, `;`, `%`, `'`, `"`, `` ` `` or `\`. If the password is
  omitted, Kibana allows anonymous access by default.
* `version` - (Required) The version to install.
  The current list of supported versions is available in the [K2 Cloud API documentation][elk-version].

~> **Note** ELK does not expose `kibana` or `logging` arguments. Kibana is an integral component of
the K2 Cloud ELK service, and the service is itself a logging destination.

~> **Note** Of the ELK-specific parameters, only `monitoring` is updated in
place. The live API can omit `password` and `options` from DescribeService.
Terraform preserves them in an existing resource state, but a fresh import
cannot recover values that the API does not return.

## GitLab argument reference

K2 Cloud provides GitLab CE (Community Edition) with a single master node.
`high_availability` must be `false`. The subnet must have internet access for
service installation. See the [K2 Cloud GitLab documentation](https://docs.k2.cloud/ru/services/paas/development/index.html).

In addition to the common arguments for all services [described above](#argument-reference),
the `gitlab` block can contain the following arguments:

* `class` - (Optional) The service class.
    * _Valid values:_ `development`
    * _Default value:_ `development`
* `ldap` - (Optional, Editable) The LDAP settings. If the block is specified, LDAP is enabled.
  The structure of this block is [described below](#gitlab-ldap).
* `logging` - (Optional, Editable) The logging settings for the service. The structure of this block is [described below](#logging).
* `monitoring` - (Optional, Editable) The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `oidc` - (Optional, Editable) The OIDC settings. If the block is specified, OIDC is enabled.
  The structure of this block is [described below](#gitlab-oidc).
* `password` - (Required) The initial GitLab root user password. Changing this argument recreates the service.
  The value must be 8 to 128 characters long and must not contain `'`, `"`, `` ` `` and `\`.
* `registry` - (Optional, Editable) The Container Registry settings. If the block is specified, Container Registry is enabled.
  The structure of this block is [described below](#gitlab-registry).
* `smtp` - (Optional, Editable) The SMTP settings. If the block is specified, SMTP is enabled.
  The structure of this block is [described below](#gitlab-smtp).

~> **Note** The `ldap`, `oidc` and `smtp` settings cannot be set on the service creation,
so they are applied by modifying the service right after it is created.

### GitLab ldap

~> **Note** All the parameters in the `ldap` block are editable.

The `ldap` block has the following structure:

* `active_directory` - (Optional) Indicates whether the LDAP server is Active Directory.
* `allow_username_or_email_login` - (Optional) Indicates whether login by username or email is allowed.
* `base` - (Required) The search base.
* `bind_dn` - (Optional) The bind distinguished name.
* `ca_cert` - (Optional) The CA certificate content.
* `ciphers` - (Optional) The SSL ciphers.
* `client_cert` - (Optional) The client certificate content.
* `client_key` - (Optional) The client private key content.
* `encryption` - (Required) The encryption method.
    * _Valid values:_ `start_tls`, `simple_tls`, `plain`
* `host` - (Required) The LDAP server host.
* `label` - (Required) The label of the login button.
* `lowercase_usernames` - (Optional) Indicates whether usernames are converted to lowercase.
* `password` - (Optional) The bind password.
* `port` - (Required) The LDAP server port.
    * _Valid values:_ From 1 to 65535
* `retry_empty_result_with_codes` - (Optional) List of response codes on which the empty result query is retried.
* `ssl_version` - (Optional) The SSL version.
    * _Valid values:_ `TLSv1_1`, `TLSv1_2`, `TLSv1_3`
* `timeout` - (Optional) The query timeout in seconds.
    * _Valid values:_ From 0 to 60
    * _Default value:_ `10`
* `uid` - (Required) The non-empty LDAP attribute used to identify a user, such as `uid` for OpenLDAP or `sAMAccountName` for Active Directory.
  Set the attribute used by your directory explicitly; the provider does not infer a default.
* `user_filter` - (Optional) The user filter.
* `verify_certificates` - (Optional) Indicates whether SSL certificates are verified.

### GitLab oidc

~> **Note** All the parameters in the `oidc` block are editable.

The `oidc` block has the following structure:

* `client_auth_method` - (Optional) The client authentication method.
    * _Valid values:_ `query`, `basic`
* `client_id` - (Optional) The client ID.
* `client_secret` - (Optional) The client secret.
* `issuer` - (Required) The URL of the OIDC provider.
* `label` - (Optional) The label of the login button.
    * _Default value:_ `OpenID Connect`
* `pkce` - (Optional) Indicates whether PKCE is enabled.
* `send_scope_to_token_endpoint` - (Optional) Indicates whether scope is sent to the token endpoint.
    * _Default value:_ `true`
* `uid_field` - (Optional) The name of the field for UID.
    * _Default value:_ `sub`

### GitLab registry

~> **Note** All the parameters in the `registry` block are editable.

The `registry` block has the following structure:

* `bucket` - (Optional) The name of the bucket in object storage where the Container Registry data is saved.
  The parameter must be set if `storage` is `s3`.
* `s3_user` - (Optional) The login of a user with permissions to the bucket in object storage.
  The parameter must be set if `storage` is `s3`.
* `storage` - (Optional) The storage type. The value cannot be changed while Container Registry is enabled.
  To change the storage type, first remove the `registry` block and apply the configuration, then add the block with the new storage type and apply again.
  Changing the storage type does not copy existing images. Migrate Registry data separately as described in the [GitLab documentation](https://docs.gitlab.com/administration/packages/container_registry/#migrate-to-object-storage-without-downtime).
    * _Valid values:_ `filesystem`, `s3`
    * _Default value:_ `filesystem`
* `url` - (Optional) The object storage endpoint URL.
  If the parameter is not specified and `storage` is `s3`, the cloud object storage endpoint is used.

### GitLab smtp

~> **Note** All the parameters in the `smtp` block are editable.

The `smtp` block has the following structure:

* `address` - (Required) The SMTP server address.
* `authentication` - (Required) The SMTP authentication method. Required by the K2 Cloud API even when credentials are omitted.
    * _Valid values:_ `plain`, `login`
* `domain` - (Optional) The SMTP domain.
* `email_from` - (Optional) The email sender address.
* `email_reply_to` - (Optional) The email reply-to address.
* `force_tls` - (Optional) Indicates whether TLS is forced. The parameter can be set only if `tls` is `false`.
* `password` - (Optional) The SMTP password.
* `port` - (Required) The SMTP server port.
    * _Valid values:_ From 1 to 65535
* `starttls_auto` - (Optional) Indicates whether STARTTLS is enabled automatically.
  The parameter can be set only if `tls` is `false`.
* `tls` - (Optional) Indicates whether TLS is used immediately.
    * _Default value:_ `false`
* `tls_verify_mode` - (Optional) The TLS certificate verification mode.
    * _Valid values:_ `none`, `peer`
* `username` - (Optional) The SMTP username.

## Memcached argument reference

In addition to the common arguments for all services [described above](#argument-reference),
the `memcached` block can contain the following arguments:

* `class` - (Optional) The service class.
    * _Valid values:_ `cacher`
    * _Default value:_ `cacher`
* `logging` - (Optional, Editable) The logging settings for the service. The structure of this block is [described below](#logging).
* `monitoring` - (Optional, Editable) The monitoring settings for the service. The structure of this block is [described below](#monitoring).

## MongoDB argument reference

In addition to the common arguments for all services [described above](#argument-reference),
the `mongodb` block can contain arguments as described below.

The following arguments are required:

* `version` - (Required) The version to install. The list of supported versions is available in the [user documentation][mongodb-version].

The following arguments are optional:

* `class` - (Optional) The service class.
    * _Valid values:_ `database`
    * _Default value:_ `database`
* `database` - (Optional, Editable) List of MongoDB databases with parameters. The structure of this block is [described below](#mongodb-database).
* `journal_commit_interval` - (Optional, Editable) The maximum interval in milliseconds between saving log data.
    * _Valid values:_ From 1 to 500
* `logging` - (Optional, Editable) The logging settings for the service. The structure of this block is [described below](#logging).
* `maxconns` - (Optional, Editable) The maximum number of concurrent connections allowed for _mongos_ or _mongod_.
    * _Valid values:_ From 10 to 51200
* `monitoring` - (Optional, Editable) The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `options` - (Optional, Editable) Map containing other MongoDB parameters.
  Parameter names must be in camelCase. Values are strings.

~> **Note** If a parameter name includes a dot, it cannot be passed in `options`.
If you need to use such a parameter, contact [technical support].

* `profile` - (Optional, Editable) Indicates which operations to profile.
    * _Valid values:_ `off`, `slowOp`, `all`
* `slowms` - (Optional, Editable) The operation time threshold in milliseconds, above which the operation is considered slow.
    * _Valid values:_ From 0 to 36000000
* `storage_engine_cache_size` - (Optional, Editable) The maximum size of internal cache in GiB used to store all data.
  A floating-point number.
    * _Valid values:_ Greater or equal to `0.25`
* `user` - (Optional, Editable) List of MongoDB users with parameters. The structure of this block is [described below](#mongodb-user).
* `quiet` - (Optional, Editable) Indicates whether the quiet mode of _mongos_ or _mongod_ is enabled.
    * _Default value:_ `false`
* `verbositylevel` - (Optional, Editable) The level of message detail in the message log.
    * _Valid values:_ `v`, `vv`, `vvv`, `vvvv`, `vvvvv`

### MongoDB database

~> **Note** All the parameters in the `database` block are editable.

The `database` block has the following structure:

* `name` - (Required) The database name.
* `user` - (Optional) List of database users with parameters. The structure of this block is [described below](#mongodb-database-user).

### MongoDB database user

~> **Note** All the parameters in the `user` block are editable.

* `name` - (Required) The MongoDB user name.
* `roles` - (Optional) List of user roles.
    * _Valid values:_ `read`, `readWrite`, `dbAdmin`, `dbOwner`

### MongoDB user

~> **Note** All the parameters in the `user` block are editable.

* `name` - (Required) The MongoDB user name.
* `password` - (Required) The MongoDB user password. The value must not contain `'`, `"`, `` ` `` and `\`.


## MySQL argument reference

In addition to the common arguments for all services [described above](#argument-reference),
the `mysql` block can contain arguments as described below.

The following arguments are required:

* `vendor` - (Required) The engine vendor.
    * _Valid values:_ `mariadb`, `percona`, `mysql`
* `version` - (Required) The version to install. The list of supported versions is available in the [user documentation][mysql-version].

The following arguments are optional:

* `class` - (Optional) The service class.
    * _Valid values:_ `database`
    * _Default value:_ `database`
* `connect_timeout` - (Optional) The number of seconds that the _mysqld_ server waits for a connect packet before responding with **Bad handshake**.
    * _Valid values:_ From 2 to 31536000
* `database` - (Optional, Editable) List of MySQL databases with parameters. The maximum number of databases is 1000.
  The structure of this block is [described below](#mysql-database).
* `galera_options` - (Optional) Map containing other Galera parameters.
  Parameter names must be in camelCase. Values are strings.

~> **Note** If a parameter name includes a dot, it cannot be passed in `galera_options`.
If you need to use such a parameter, contact [technical support].

* `gcache_size` - (Optional) A Galera parameter. The size of GCache circular buffer storage preallocated on startup in bytes.
    * _Valid values:_ From 128 MiB
    * _Constraints:_ The parameter can be set only if `high_availability` is `true`
* `gcs_fc_factor` - (Optional) A Galera parameter. The fraction of `gcs_fc_limit` at which replication is resumed
  when the recv queue length falls below this value.
    * _Valid values:_ From 0.0 to 1.0
    * _Constraints:_ The parameter can be set only if `high_availability` is `true`
* `gcs_fc_limit` - (Optional) A Galera parameter. The number of writesets. If the recv queue length exceeds it replication is suspended.
  Replication will resume according to the `gcs_fc_factor` setting.
    * _Valid values:_ From 1 to 2147483647
    * _Constraints:_ The parameter can be set only if `high_availability` is `true`
* `gcs_fc_master_slave` - (Optional) A Galera parameter. Indicates whether the cluster has only one source node.
    * _Constraints:_ The parameter can be set only if `high_availability` is `true`

~> **Note** `gcs_fc_master_slave` is deprecated. This parameter is relevant for Percona 5.7.
Use `gcs_fc_single_primary` instead.

* `gcs_fc_single_primary` - (Optional) A Galera parameter. Indicates whether there is more than one replication source.
    * _Constraints:_ The parameter can be set only if `high_availability` is `true`

~> **Note** `gcs_fc_single_primary` replaces the deprecated `gcs_fc_master_slave` parameter.
This parameter is relevant for Percona 8.0, MySQL 8.0, and MariaDB 10.4, 10.5, 10.6 and 10.11.

* `innodb_buffer_pool_instances` - (Optional) The number of regions that `innodb_buffer_pool_size` is divided into
  when `innodb_buffer_pool_size` > 1 GiB. This parameter is relevant for Percona 5.7, 8.0 и MariaDB 10.2, 10.3, 10.4.
    * _Valid values:_ From 1 to 64
* `innodb_buffer_pool_size` - (Optional) The size in bytes of the buffer pool used to cache table data and indexes.
    * _Valid values:_ From 128 MiB
* `innodb_change_buffering` - (Optional) Operations for which change buffering optimization is enabled.
    * _Valid values:_ `inserts`, `deletes`, `changes`, `purges`, `all`, `none`
* `innodb_flush_log_at_trx_commit` - (Optional) The value of the parameter controls the behaviour for transaction commit operations.
    * _Valid values:_ From 0 to 2
  For more information about the parameter, see the [MySQL documentation][doc-innodb_flush_log_at_trx_commit].
* `innodb_io_capacity` - (Optional) The number of I/O operations per second (IOPS) available to InnoDB background tasks.
    * _Valid values:_ From 100 to 9223372036854775807
* `innodb_io_capacity_max` - (Optional) The maximum number of IOPS that InnoDB background tasks can perform.
    * _Valid values:_ From 100 to 9223372036854775807
* `innodb_log_file_size` - (Optional) The size of a single file in bytes in the redo system log.
    * _Valid values:_ From 4 MiB to 4 GiB
* `innodb_log_files_in_group` - (Optional) The number of system log files in a log group.
    * _Valid values:_ From 2 to 100
* `innodb_purge_threads` - (Optional) The number of background threads allocated for the InnoDB purge operation.
    * _Valid values:_ From 1 to 32
* `innodb_thread_concurrency` - (Optional) The maximum number of threads permitted inside of InnoDB.
  This parameter is relevant for Percona 5.7, 8.0 and MariaDB 10.2, 10.3, 10.4.
    * _Valid values:_ From 0 to 1000
* `innodb_strict_mode` - (Optional) The MySQL operation mode.
    * _Valid values:_ `ON`, `OFF`
    * _Default value:_ `OFF`
  For more information about the parameter, see the [MySQL documentation][doc-innodb_strict_mode].
* `innodb_sync_array_size` - (Optional) The size of the mutex/lock wait array.
  This parameter is relevant for Percona 5.7, 8.0 and MariaDB 10.2, 10.3, 10.4.
    * _Valid values:_ From 1 to 1024
* `max_allowed_packet` - (Optional) The maximum size of one packet, any generated/intermediate string  or any parameter sent by the _mysql_stmt_send_long_data()_ C API function.
    * _Valid values:_ From 16 MiB to 1 GiB
    * _Default value:_ `16777216` (16 MiB)
* `max_connect_errors` - (Optional) The maximum number of connection errors, at which the server blocks the host from further connections.
    * _Valid values:_ From 1 to 9223372036854775807
* `max_connections` - (Optional) The maximum permitted number of simultaneous client connections that a host can handle.
    * _Valid values:_ From 10 to 100000
* `max_heap_table_size` - (Optional) The maximum size in bytes, to which user-created `MEMORY` tables are permitted to grow.
    * _Valid values:_ From 16 KiB to 4 GiB
* `logging` - (Optional, Editable) The logging settings for the service. The structure of this block is [described below](#logging).
* `monitoring` - (Optional, Editable) The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `options` - (Optional) Map containing other MySQL parameters.
  Parameter names must be in camelCase. Values are strings.

~> **Note** If a parameter name includes a dot, it cannot be passed in `options`.
If you need to use such a parameter, contact [technical support].

* `pxc_strict_mode` - (Optional) PXC mode. For more information about the parameter, see the [Percona documentation][doc-pxc_strict_mode].
    * _Valid values:_ `DISABLED`, `PERMISSIVE`, `ENFORCING`, `MASTER`
    * _Constraints:_ The parameter can be set only if `high_availability` is `true`
* `table_open_cache` - (Optional) The number of open tables for all threads.
    * _Valid values:_ From 1 to 1048576
* `thread_cache_size` - (Optional) The number of threads that the server caches to establish new network connections.
    * _Valid values:_ From 0 to 16 KiB
* `tmp_table_size` - (Optional) The maximum size of internal in-memory temporary tables in bytes.
    * _Valid values:_ From 1 KiB to 4 GiB
* `transaction_isolation` - (Optional) The transaction isolation level.
  For more information about the parameter, see the [MySQL documentation][doc-transaction_isolation].
    * _Valid values:_ `READ-UNCOMMITTED`, `READ-COMMITTED`, `REPEATABLE-READ`, `SERIALIZABLE`
* `user` - (Optional, Editable) List of MySQL users with parameters. The maximum number of users is 1000.
  The structure of this block is [described below](#mysql-user).
* `wait_timeout` - (Optional) The number of seconds the server waits for activity on a noninteractive connection before closing it.
    * _Valid values:_ From 1 to 31536000

### MySQL database

~> **Note** All the parameters in the `database` block are editable.

The `database` block structure is described below.

The following arguments are required:

* `name` - (Required) The database name.

The following arguments are optional:

* `backup_enabled` - (Optional) Indicates whether backup is enabled for the database.
    * _Default value:_ `false`
* `backup_id` - (Optional) The database backup ID.
* `backup_db_name` - (Optional) The name of a database from the backup specified in the `backup_id` parameter.
* `charset` - (Optional) The database charset.
    * _Valid values:_ Depend on `vendor`.
      For `mariadb` see the [MariaDB documentation][doc-mariadb-charset-collate].
      For `percona` and `mysql` see the [MySQL documentation][doc-mysql-charset-collate].
* `collate` - (Optional) The database collation.
    * _Valid values:_ Depend on `vendor`.
      For `mariadb` see the [MariaDB documentation][doc-mariadb-charset-collate].
      For `percona` and `mysql` see the [MySQL documentation][doc-mysql-charset-collate].
* `user` - (Optional) List of database users with parameters. The maximum number of users is 1000.
  The structure of this block is [described below](#mysql-database-user).

### MySQL database user

~> **Note** All the parameters in the `user` block are editable.

The `user` block structure is described below.

The following arguments are required:

* `name` - (Required) The MySQL user name.

The following arguments are optional:

* `options` - (Optional) List of user options.
    * _Valid values:_ `ALTER`, `ALTER ROUTINE`, `CREATE`, `CREATE ROUTINE`, `CREATE TEMPORARY TABLES`, `CREATE VIEW`, `DELETE`, `DROP`, `EVENT`, `EXECUTE`, `INDEX`, `INSERT`, `LOCK TABLES`, `SELECT`, `SHOW VIEW`, `TRIGGER`, `UPDATE`
* `privileges` - (Optional) List of user privileges.
    * _Valid values:_ `GRANT`, `NONE`

### MySQL user

~> **Note** All the parameters in the `user` block are editable.

The `user` block structure is described below.

The following arguments are required:

* `name` - (Required) The MySQL user name.
* `password` - (Required) The MySQL user password. The value must not contain `'`, `"`, `` ` `` and `\`.

The following arguments are optional:

* `host` - (Optional) The hostname or IP address. The value must be 1 to 60 characters long.

## PostgreSQL argument reference

In addition to the common arguments for all services [described above](#argument-reference),
the `pgsql` block can contain arguments described below.

The following arguments are required:

* `version` - (Required) The version to install. The list of supported versions is available in the [user documentation][pgsql-version].

The following arguments are optional. Except for `replication_mode`, they can be updated without recreating the service:
only changed arguments are sent with `PatchServiceParameters`. An argument removed from the configuration keeps its current value
(`max_parallel_maintenance_workers` and `wal_keep_segments` keep showing a diff until set explicitly).

* `autovacuum` - (Optional) Indicates whether the server must run the autovacuum launcher daemon.
    * _Valid values:_ `ON`, `OFF`
* `autovacuum_max_workers` - (Optional) The maximum number of autovacuum processes (other than the autovacuum launcher)
  that can run simultaneously.
    * _Valid values:_ From 1 to 262143
* `autovacuum_vacuum_cost_delay` - (Optional) The cost delay value in milliseconds used in automatic `VACUUM` operations.
    * _Valid values:_ `-1`, from 1 to 100
* `autovacuum_vacuum_cost_limit` - (Optional) The cost limit value used in automatic `VACUUM` operations.
    * _Valid values:_ `-1`, from 1 to 10000
* `autovacuum_analyze_scale_factor` - (Optional) The fraction of the table size to add to `autovacuum_analyze_threshold`
  when deciding whether to trigger `ANALYZE`.
    * _Valid values:_ From 0 to 100
* `autovacuum_vacuum_scale_factor` - (Optional) The fraction of the table size to add to `autovacuum_vacuum_threshold`
  when deciding whether to trigger `VACUUM`.
    * _Valid values:_ From 0 to 100
* `class` - (Optional) The service class.
    * _Valid values:_ `database`
    * _Default value:_ `database`
* `database` - (Optional, Editable) List of PostgreSQL databases with parameters. The maximum number of databases is 1000.
  The structure of this block is [described below](#postgresql-database).
* `effective_cache_size` - (Optional) The planner’s assumption about the effective size of the disk cache (in bytes, multiple of 1 KiB),
  that is available to a single query.
    * _Valid values:_ From 8 to 17179869176 KiB
  For more information about the parameter, see the [PostgreSQL documentation][doc-effective_cache_size].
* `effective_io_concurrency` - (Optional) The number of concurrent disk I/O operations.
    * _Valid values:_ From 0 to 1000
* `logging` - (Optional, Editable) The logging settings for the service. The structure of this block is [described below](#logging).
* `maintenance_work_mem` - (Optional) The maximum amount of memory (in bytes, multiple of 1 KiB) used by maintenance operations, such as `VACUUM`, `CREATE INDEX`, and `ALTER TABLE ADD FOREIGN KEY`.
    * _Valid values:_ From 1 MiB to 2 GiB
* `max_connections` - (Optional) The maximum number of simultaneous connections to the database server.
    * _Valid values:_ From 1 to 262143
* `max_wal_size` - (Optional, **Deprecated**) The maximum size (in bytes, multiple of 1 MiB) that WAL can reach at automatic checkpoints.
    * _Valid values:_ From 2 to 2147483647 MiB

~> **Note** The parameter `max_wal_size` is marked as deprecated since it is not supported for PostgreSQL services
starting with the _paas_v4_0_ environment version.

* `max_parallel_maintenance_workers` - (Optional) The maximum number of parallel workers that a single utility command can start.
    * _Valid values:_ From 0 to 1024
    * _Constraints:_ This parameter is relevant only for PostgreSQL versions 11 and higher
* `max_parallel_workers` - (Optional) The maximum number of workers that the system can support for parallel operations.
* `max_parallel_workers_per_gather` - (Optional) The maximum number of workers that a single _Gather_ node can start.
    * _Valid values:_ From 0 to 1024
* `max_worker_processes` - (Optional) The maximum number of background processes that the system can support.
    * _Valid values:_ From 0 to 262143
* `min_wal_size` - (Optional, **Deprecated**) The minimum size (in bytes, multiple of 1 MiB) to shrink the WAL to. As long as WAL disk usage stays below this setting, old WAL files are always recycled for future use at a checkpoint, rather than removed.
    * _Valid values:_ From 32 to 2147483647 MiB

~> **Note** The parameter `min_wal_size` is marked as deprecated since it is not supported for PostgreSQL services
starting with the _paas_v4_0_ environment version.

* `monitoring` - (Optional, Editable) The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `options` - (Optional) Map containing other PostgreSQL parameters.
  Parameter names must be in camelCase. Values are strings.

~> **Note** If a parameter name includes a dot, it cannot be passed in `options`.
If you need to use such a parameter, contact [technical support].

* `replication_mode` - (Optional, Forces new resource) The replication mode in the _Patroni_ cluster.
    * _Valid values:_ `asynchronous`, `synchronous`, `synchronous_strict`
    * _Constraints:_ The parameter must be set if `high_availability` is `true`
* `shared_buffers` - (Optional) The amount of memory (in bytes, multiple of 1 KiB) the database server uses for shared memory buffers.
    * _Valid values:_ From 128 to 8589934584 KiB
  For more information about the parameter, see the [PostgreSQL documentation][doc-shared_buffers].
* `user` - (Optional, Editable) List of PostgreSQL users with parameters. The maximum number of users is 1000.
  The structure of this block is [described below](#postgresql-user).
* `wal_buffers` - (Optional) The amount of shared memory in 8 KiB pages used for WAL data not yet written to a volume.
    * _Valid values:_ From 8 to 262143
* `wal_keep_segments` - (Optional) The minimum number of log files segments that must be kept in the _pg_xlog_ directory, in case a standby server needs to fetch them for streaming replication.
    * _Valid values:_ From 0 to 2147483647
    * _Constraints:_ This parameter is relevant only for PostgreSQL version 12
* `work_mem` - (Optional) The base maximum amount of memory (in bytes, multiple of 1 KiB) to be used by a query operation (such as a sort or hash table) before writing to temporary disk files.
    * _Valid values:_ From 64 to 2147483647 KiB

### PostgreSQL database

~> **Note** All the parameters in the `database` block are editable.

The `database` block structure is described below.

The following arguments are required:

* `name` - (Required) The database name.

The following arguments are optional:

* `backup_enabled` - (Optional) Indicates whether backup is enabled for the database.
    * _Default value:_ `false`
* `backup_id` - (Optional) The database backup ID.
* `backup_db_name` - (Optional) The name of a database from the backup specified in the `backup_id` parameter.
* `encoding` - (Optional) The database encoding.
* `extensions` - (Optional) List of extensions for the database.
    * _Valid values:_ `address_standardizer`, `address_standardizer_data_us`, `amcheck`, `autoinc`, `bloom`, `btree_gin`, `btree_gist`, `citext`,`cube`, `dblink`, `dict_int`, `dict_xsyn`, `earthdistance`, `fuzzystrmatch`, `hstore`, `intarray`, `isn`, `lo`, `ltree`, `moddatetime`, `pg_buffercache`, `pg_trgm`, `pg_visibility `, `pgcrypto`, `pgrowlocks`, `pgstattuple`, `postgis`, `postgis_tiger_geocoder`, `postgis_topology`, `postgres_fdw`, `seg`, `tablefunc`, `tcn`, `timescaledb`,  `tsm_system_rows`, `tsm_system_time`, `unaccent`, `uuid-ossp`, `xml2`
* `locale` - (Optional) The database locale.
* `owner` - (Required) The name of the user who is the database owner. This must be one of the existing users.
  Such a user cannot be deleted as long as it is the database owner.
* `user` - (Optional) List of PostgreSQL users with parameters. The maximum number of users is 1000.
  The structure of this block is [described below](#postgresql-database-user).

### PostgreSQL database user

~> **Note** All the parameters in the `user` block are editable.

The `user` block has the following structure:

* `name` - (Required) The PostgreSQL user name.

### PostgreSQL user

~> **Note** All the parameters in the `user` block are editable.

The `user` block has the following structure:

* `name` - (Required) The PostgreSQL user name.
* `password` - (Required) The PostgreSQL user password. The value must not contain `'`, `"`, `` ` `` and `\`.

## Kafka argument reference

In addition to the common arguments for all services [described above](#argument-reference),
the `kafka` block can contain the following arguments:

* `class` - (Optional) The service class.
    * _Valid values:_ `message_broker`
    * _Default value:_ `message_broker`
* `logging` - (Optional, Editable) The logging settings for the service. The structure of this block is [described below](#logging).
* `monitoring` - (Optional, Editable) The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `options` - (Optional, Editable) Map containing other Kafka parameters.
  Parameter names must be in camelCase. Values can be strings, numbers, or booleans.

~> **Note** If a parameter name includes a dot, it cannot be passed in the `options`.
If you need to use such a parameter, contact [technical support].

* `version` - (Required) The version to install.
  The list of supported versions is available in the [user documentation][kafka-version].

### Kafka options example

```terraform
kafka {
  version = "3.7.0"

  options = {
    autoCreateTopicsEnable      = true
    uncleanLeaderElectionEnable = false
    logRetentionHours           = 168
  }
}
```

## Prometheus argument reference

In addition to the common arguments for all services [described above](#argument-reference),
the `prometheus` block can contain the following arguments:

* `class` - (Optional) The service class.
    * _Valid values:_ `monitoring`
    * _Default value:_ `monitoring`
* `remote_write_receiver` - (Optional, Editable) Whether the Prometheus service accepts metrics through the Remote Write protocol.
    * _Default value:_ `false`
    * _Constraints:_ Supported for Prometheus environment versions `paas_v4_0` and later

~> **Note** Prometheus supports only a single-node configuration. Set the common `high_availability` argument to `false`.

## RabbitMQ argument reference

In addition to the common arguments for all services [described above](#argument-reference),
the `rabbitmq` block can contain the following arguments:

* `password` - (Required, Editable) The RabbitMQ admin password.
  The value must be 8 to 32 characters long and must not contain `-`, `|`, `[`, `]`, `'`, `"`, `;` and `\`.
* `version` - (Required) The version to install. The list of supported versions is available in the [user documentation][rabbitmq-version].
* `class` - (Optional) The service class.
    * _Valid values:_ `message_broker`
    * _Default value:_ `message_broker`
* `logging` - (Optional, Editable) The logging settings for the service. The structure of this block is [described below](#logging).
* `monitoring` - (Optional, Editable) The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `options` - (Optional, Editable) Map containing other RabbitMQ parameters.
  Parameter names must be in camelCase. Values are strings.

~> **Note** If a parameter name includes a dot, it cannot be passed in the `options`.
If you need to use such a parameter, contact [technical support].

## Redis argument reference

In addition to the common arguments for all services [described above](#argument-reference),
the `redis` block can contain arguments as described below.

The following arguments are required:

* `version` - (Required) The version to install. The list of supported versions is available in the [user documentation][redis-version].

The following arguments are optional:

* `class` - (Optional) The service class.
    * _Valid values:_ `cacher`, `database`
    * _Default value:_ `cacher`
* `cluster_type` - (Optional) The clustering option.
  The parameter must be set if `high_availability` is `true`.
    * _Valid values:_ `native`, `sentinel`
* `databases` - (Optional, Editable) The number of databases.
    * _Valid values:_ From 1 to 2147483647
* `logging` - (Optional, Editable) The logging settings for the service. The structure of this block is [described below](#logging).
* `maxmemory_policy` - (Optional, Editable) The memory management mode.
    * _Valid values:_ `noeviction`, `allkeys-lru`, `allkeys-lfu`, `volatile-lru`, `volatile-lfu`, `allkeys-random`, `volatile-random`, `volatile-ttl`
* `monitoring` - (Optional, Editable) The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `options` - (Optional, Editable) Map containing other Redis parameters.
  Parameter names must be in camelCase. Values are strings.

~> **Note** If a parameter name includes a dot, it cannot be passed in `options`.
If you need to use such a parameter, contact [technical support].

* `password` - (Optional) The Redis user password. The value must not contain `'`, `"`, `` ` `` and `\`.
* `persistence_aof` - (Optional, Editable) Indicates whether the AOF storage mode is enabled.
    * _Default value:_ `false`
* `persistence_rdb` - (Optional, Editable) Indicates whether the RDB storage mode is enabled.
    * _Default value:_ `true`
* `timeout` - (Optional, Editable) The time in seconds during which the connection to an inactive client is retained.
    * _Valid values:_ From 0 to 2147483647
* `tcp_backlog` - (Optional, Editable) The size of a connection queue.
    * _Valid values:_ From 1 to 4096
* `tcp_keepalive` - (Optional, Editable) The time in seconds during which the service sends ACKs to detect dead peers (unreachable clients).
  The value must be non-negative.

## Common service argument reference

### logging

~> **Note** All the parameters in the `logging` block are editable.

The `logging` block structure is described below.

The following arguments are required:

* `log_to` - (Required) The ID of the logging service. It must run in the same VPC as the service.

The following arguments are optional:

* `logging_tags` - (Optional) List of tags that are assigned to the log records of the service.
  Each value in the list must be 1 to 256 characters long.

### monitoring

~> **Note** All the parameters in the `monitoring` block are editable.

The `monitoring` block structure is described below.

The following arguments are required:

* `monitor_by` - (Required) The ID of the monitoring service. It must run in the same VPC as the service.

The following arguments are optional:

* `monitoring_labels` - (Optional) Map containing labels that are assigned to the metrics of the service.
  Keys must be 1 to 64 characters long.

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `auto_created_security_group_ids` - List of security group IDs that the cloud created for the service.
* `available_environment_versions` - The environment versions to which the current version can be updated.
* `endpoints` - List of endpoints for connecting to the service. The structure of this block is [described below](#endpoints).
* `environment_version` - The current version of the service environment.
* `error_code` - The service error code.
* `error_description` - The detailed description of the service error.
* `id` - The ID of the PaaS service.
* `instances` - List of instances that refers to the service. The structure of this block is [described below](#instances).
* `service_class` - The service class. The value matches the `class` parameter of the specified block with service parameters.
* `service_type` - The service type. The value matches the name of the specified block with service parameters.
* `status` - The current status of the service.
* `supported_features` - List of service features.
* `total_cpu_count` - Total number of CPU cores in use.
* `total_memory` - Total RAM in use in MiB.
* `vpc_id` - The ID of the VPC.

For `backup_settings` the following attribute is also exported:

* `user_id` - The ID of the user whose login is set to `backup_settings.user_login`.

For `*.database` the following attribute is also exported:

* `id` - The ID of the database.

For `*.user` the following attribute is also exported:

* `id` - The ID of the user.

For `gitlab` the following attribute is also exported:

* `external_url` - The URL for accessing GitLab.

### endpoints

The `endpoints` block has the following structure:

* `addresses` - List of addresses for connecting to the service.
* `name` - The name of the endpoint.

### instances

* `endpoints` - List of service endpoints on the instance. The structure of this block is [described below](#instance-endpoints).
* `index` - The instance index.
* `instance_id` - The ID of the instance.
* `interface_id` - The ID of the instance network interface.
* `name` - The instance name.
* `private_ip` - The private IP address of the instance.
* `role` - The instance role.
* `status` - The current status of the instance.

#### instance endpoints

The `endpoints` block has the following structure:

* `address` - The address of the endpoint.
* `name` - The name of the endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

* `create` - (Default `30 minutes`) How long to wait for the service to be created.
* `update` - (Default `60 minutes`) How long to wait for the service to be updated.
* `delete` - (Default `15 minutes`) How long to wait for the service to be deleted.

## Import

PaaS service can be imported using `id`, e.g.,

```
$ terraform import aws_paas_service.example fm-cluster-12345678
```
