---
subcategory: "PaaS"
layout: "aws"
page_title: "aws_paas_service"
description: |-
  Provides information about a PaaS service.
---

# Data Source: aws_paas_service

Provides information about a PaaS service.

## Example Usage

```terraform
data "aws_paas_service" "selected" {
  id = "fm-cluster-12345678"
}

output "paas_service_name" {
  value = data.aws_paas_service.selected.name
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the PaaS service (e.g. `fm-cluster-12345678`).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `auto_disk_expansion` - The auto volume expansion settings for the service. The structure of this block is [described below](#auto_disk_expansion).
  The block is empty if auto volume expansion has never been configured for the service.
* `available_environment_versions` - The environment versions to which the current version can be updated.
* `auto_created_security_group_ids` - List of security group IDs that cloud created for the service.
* `backup_settings` - The backup settings for the service. The structure of this block is [described below](#backup_settings).
* `data_volume` - The data volume parameters for the service. The structure of this block is [described below](#data_volume).
* `endpoints` - List of endpoints for connecting to the service. The structure of this block is [described below](#endpoints).
* `environment_version` - The current version of the service environment.
* `error_code` - The service error code.
* `error_description` - The detailed description of the service error.
* `high_availability` - Indicates whether this is a high-availability service.
* `instances` - List of instances that refers to the service. The structure of this block is [described below](#instances).
* `instance_qty` - The total number of service instances including coordinator instances.
* `instance_type` - The instance type.
* `name` - The service name.
* `network_interface_ids` - List of network interface IDs.
* `nodes` - The service node topology. The structure of this block is [described below](#nodes).
* `root_volume` - The root volume parameters for the service. The structure of this block is [described below](#root_volume).
* `security_group_ids` - List of security group IDs that were specified for the service.
* `service_class` - The service class.
* `service_type` - The service type. This value determines which service parameters are included in the corresponding block.
    * `clickhouse` - ClickHouse parameters. The structure of this block is [described below](#clickhouse-attribute-reference).
    * `elk` - ELK parameters. The structure of this block is [described below](#elk-attribute-reference).
    * `elasticsearch` - Elasticsearch parameters. The structure of this block is [described below](#elasticsearch-attribute-reference).
    * `gitlab` - GitLab parameters. The structure of this block is [described below](#gitlab-attribute-reference).
    * `kafka` - Kafka parameters. The structure of this block is [described below](#kafka-attribute-reference).
    * `memcached` - Memcached parameters. The structure of this block is [described below](#memcached-attribute-reference).
    * `mongodb` - MongoDB parameters. The structure of this block is [described below](#mongodb-attribute-reference).
    * `mysql` - MySQL parameters. The structure of this block is [described below](#mysql-attribute-reference).
    * `pgsql` - PostgreSQL parameters. The structure of this block is [described below](#postgresql-attribute-reference).
    * `prometheus` - Prometheus parameters. The structure of this block is [described below](#prometheus-attribute-reference).
    * `rabbitmq` - RabbitMQ parameters. The structure of this block is [described below](#rabbitmq-attribute-reference).
    * `redis` - Redis parameters. The structure of this block is [described below](#redis-attribute-reference).
* `ssh_key_name` - The name of the SSH key for accessing instances.
* `status` - The current status of the service.
* `subnet_ids` - List of subnet IDs.
* `supported_features` - List of service features.
* `total_cpu_count` - Total number of CPU cores in use.
* `total_memory` - Total RAM in use in MiB.
* `vpc_id` - The ID of the VPC.

### auto_disk_expansion

The `auto_disk_expansion` block has the following structure:

* `condition_value` - The utilization level of the data volume (in percent), upon exceeding which,
  within 10 minutes on at least one node, automatic expansion is triggered.
* `enabled` - Indicates whether auto volume expansion is enabled.
* `max_size_gb` - The limit in GiB up to which the volume can be increased.
* `status` - The current status of auto volume expansion (`active`, `warning`, `running` or `failed`).

### backup_settings

The `backup_settings` block has the following structure:

* `bucket_name` - The name of the bucket in object storage where the service backup is saved.
* `enabled` - Indicates whether backup is enabled for the service.
* `expiration_days` - The backup retention period in days.
* `notification_email` - The email address to which a notification that backup was created is sent.
* `start_time` - The time when the daily backup process starts. It is set as a string in the HH:MM format Moscow time.
* `user_id` - The ID of a user with write permissions to the bucket in object storage.
* `user_login` - The login of a user with write permissions to the bucket in object storage.

### data_volume

The `data_volume` block has the following structure:

* `iops` - The number of read/write operations per second for the data volume.
* `size` - The size of the data volume in GiB.
* `type` - The type of the data volume.

### endpoints

The `endpoints` block has the following structure:

* `addresses` - List of addresses for connecting to the service.
* `name` - The name of the endpoint.

### instances

The `instances` block has the following structure:

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

### nodes

The `nodes` block has the following structure:

* `coordinator` - The coordinator node parameters.
    * `role` - The coordinator node role.
* `main` - The main node parameters.
    * `role` - The main node role.

### root_volume

The `root_volume` block has the following structure:

* `iops` - The number of read/write operations per second for the root volume.
* `size` - The size of the root volume in GiB.
* `type` - The type of the root volume.

## ClickHouse Attribute Reference

~> **Note** The following attributes contain default parameter values or user-defined values used when the service is created.

In addition to the common attributes for all services [described above](#attribute-reference),
the following attributes are exported only for a ClickHouse service:

* `class` - The service class.
* `database` - List of ClickHouse databases with parameters. The structure of this block is [described below](#clickhouse-database).
* `logging` - The logging settings for the service. The structure of this block is [described below](#logging).
* `monitoring` - The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `options` - Advanced ClickHouse server settings.
* `user` - List of ClickHouse users with parameters. The structure of this block is [described below](#clickhouse-user).
* `version` - The installed version.

### ClickHouse database

The `database` block has the following structure:

* `id` - The ID of the database.
* `name` - The database name.

### ClickHouse user

The `user` block has the following structure:

* `grants` - Configured `GRANT` statements for the user.
* `id` - The ID of the user.
* `name` - The ClickHouse user name.
* `networks` - Networks from which the user can connect. The structure of this block is [described below](#clickhouse-user-networks).
* `password` - The ClickHouse user password.
* `password_type` - The password encoding.
* `profile` - The name of the settings profile assigned to the user.
* `quota` - The name of the quota assigned to the user.

~> **Note** An empty `grants` value represents the backend defaults, which can include broad privileges.
See the [resource's grants configuration](../resources/paas_service.md#clickhouse-user) for explicit privilege restrictions.

#### ClickHouse user networks

The `networks` block has the following structure:

* `host` - List of host names from which the user can connect.
* `host_regexp` - List of regular expressions for host names from which the user can connect.
* `ip` - List of IP addresses or subnets from which the user can connect.

## Elasticsearch Attribute Reference

~> **Note** The following attributes contain default parameter values or user-defined values used when the service is created.

In addition to the common attributes for all services [described above](#attribute-reference),
the following attributes are exported only for an Elasticsearch service:

* `allow_anonymous` - Indicates whether anonymous access is enabled.
* `anonymous_role` - The role for anonymous access.
* `class` - The service class.
* `kibana` - Indicates whether the Kibana deployment is enabled.
* `logging` - The logging settings for the service. The structure of this block is [described below](#logging).
* `monitoring` - The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `options` - Other Elasticsearch parameters.
* `password` - The Elasticsearch user password.
* `version` - The installed version.

## ELK Attribute Reference

~> **Note** The following attributes contain default parameter values or user-defined values used when the service is created.

In addition to the common attributes for all services [described above](#attribute-reference),
the following attributes are exported only for an ELK service:

* `allow_anonymous` - Indicates whether anonymous access to Kibana is enabled.
* `anonymous_role` - List containing zero or one role for anonymous access.
* `class` - The service class (`logging`).
* `monitoring` - The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `options` - Other ELK parameters when they are returned by the API.
* `password` - The Elasticsearch user password when it is returned by the API. This attribute is sensitive.
* `version` - The installed version.

~> **Note** The live DescribeService API can omit `password` and `options`.
Unlike an existing resource state, a data source cannot reconstruct create
inputs that the API does not return.

## GitLab Attribute Reference

~> **Note** The following attributes contain default parameter values or user-defined values used when the service is created.

In addition to the common attributes for all services [described above](#attribute-reference),
the following attributes are exported only for a GitLab service:

* `class` - The service class.
* `external_url` - The URL for accessing GitLab.
* `ldap` - The LDAP settings. The structure of this block is [described below](#gitlab-ldap).
* `logging` - The logging settings for the service. The structure of this block is [described below](#logging).
* `monitoring` - The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `oidc` - The OIDC settings. The structure of this block is [described below](#gitlab-oidc).
* `password` - The GitLab root user password.
* `registry` - The Container Registry settings. The structure of this block is [described below](#gitlab-registry).
* `smtp` - The SMTP settings. The structure of this block is [described below](#gitlab-smtp).

### GitLab ldap

The `ldap` block has the following structure:

* `active_directory` - Indicates whether the LDAP server is Active Directory.
* `allow_username_or_email_login` - Indicates whether login by username or email is allowed.
* `base` - The search base.
* `bind_dn` - The bind distinguished name.
* `ca_cert` - The CA certificate content.
* `ciphers` - The SSL ciphers.
* `client_cert` - The client certificate content.
* `client_key` - The client private key content.
* `encryption` - The encryption method.
* `host` - The LDAP server host.
* `label` - The label of the login button.
* `lowercase_usernames` - Indicates whether usernames are converted to lowercase.
* `password` - The bind password.
* `port` - The LDAP server port.
* `retry_empty_result_with_codes` - List of response codes on which the empty result query is retried.
* `ssl_version` - The SSL version.
* `timeout` - The query timeout in seconds.
* `uid` - The UID attribute.
* `user_filter` - The user filter.
* `verify_certificates` - Indicates whether SSL certificates are verified.

### GitLab oidc

The `oidc` block has the following structure:

* `client_auth_method` - The client authentication method.
* `client_id` - The client ID.
* `client_secret` - The client secret.
* `issuer` - The URL of the OIDC provider.
* `label` - The label of the login button.
* `pkce` - Indicates whether PKCE is enabled.
* `send_scope_to_token_endpoint` - Indicates whether scope is sent to the token endpoint.
* `uid_field` - The name of the field for UID.

### GitLab registry

The `registry` block has the following structure:

* `bucket` - The name of the bucket in object storage where the Container Registry data is saved.
* `s3_user` - The login of a user with permissions to the bucket in object storage.
* `storage` - The storage type.
* `url` - The object storage endpoint URL.

### GitLab smtp

The `smtp` block has the following structure:

* `address` - The SMTP server address.
* `authentication` - The SMTP authentication method.
* `domain` - The SMTP domain.
* `email_from` - The email sender address.
* `email_reply_to` - The email reply-to address.
* `force_tls` - Indicates whether TLS is forced.
* `password` - The SMTP password.
* `port` - The SMTP server port.
* `starttls_auto` - Indicates whether STARTTLS is enabled automatically.
* `tls` - Indicates whether TLS is used immediately.
* `tls_verify_mode` - The TLS certificate verification mode.
* `username` - The SMTP username.

## Memcached Attribute Reference

~> **Note** The following attributes contain default parameter values or user-defined values used when the service is created.

In addition to the common attributes for all services [described above](#attribute-reference),
the following attributes are exported only for a Memcached service:

* `class` - The service class.
* `logging` - The logging settings for the service. The structure of this block is [described below](#logging).
* `monitoring` - The monitoring settings for the service. The structure of this block is [described below](#monitoring).

## MongoDB Attribute Reference

~> **Note** The following attributes contain default parameter values or user-defined values used when the service is created.

In addition to the common attributes for all services [described above](#attribute-reference),
the following attributes are exported only for a Memcached service:

* `class` - The service class.
* `database` - List of MongoDB databases with parameters. The structure of this block is [described below](#mongodb-database).
* `journal_commit_interval` - The maximum interval in milliseconds between saving log data.
* `logging` - The logging settings for the service. The structure of this block is [described below](#logging).
* `maxconns` - The maximum number of concurrent connections allowed for _mongos_ or _mongod_.
* `monitoring` - The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `options` - Other MongoDB parameters.
* `profile` - Indicates which operations to profile.
* `slowms` - The operation time threshold in milliseconds, above which the operation is considered slow.
* `storage_engine_cache_size` - The maximum size of internal cache in GiB used to store all data.
* `user` - List of MongoDB users with parameters. The structure of this block is [described below](#mongodb-user).
* `quiet` - Indicates whether the quiet mode of _mongos_ or _mongod_ is enabled.
* `verbositylevel` - The level of message detail in the message log.
* `version` - The installed version.

### MongoDB database

The `database` block has the following structure:

* `backup_enabled` - Indicates whether backup is enabled for the database.
* `id` - The ID of the database.
* `name` - The database name.
* `user` - List of database users with parameters. The structure of this block is [described below](#mongodb-database-user).

### MongoDB database user

The `user` block has the following structure:

* `id` - The ID of the user.
* `name` - The MongoDB user name.
* `roles` - List of user roles.

### MongoDB user

The `user` block has the following structure:

* `id` - The ID of the user.
* `name` - The MongoDB user name.
* `password` - The MongoDB user password.

## MySQL Attribute Reference

~> **Note** The following attributes contain default parameter values or user-defined values used when the service is created.

In addition to the common attributes for all services [described above](#attribute-reference),
the following attributes are exported only for a MySQL service:

* `bin_log_enabled` - Indicates whether the binary log is enabled.
* `binlog_retention_days` - The binary log retention period in days.
* `class` - The service class.
* `connect_timeout` - The number of seconds that the _mysqld_ server waits for a connect packet before responding with **Bad handshake**.
* `database` - List of MySQL databases with parameters. The structure of this block is [described below](#mysql-database).
* `galera_options` - Other Galera parameters.
* `gcache_size` - A Galera parameter. The size of GCache circular buffer storage preallocated on startup in bytes.
* `gcs_fc_factor` - A Galera parameter. The fraction of `gcs_fc_limit` at which replication is resumed
  when the recv queue length falls below this value.
* `gcs_fc_limit` - A Galera parameter. The number of writesets. If the recv queue length exceeds it replication is suspended.
* `gcs_fc_master_slave` - A Galera parameter. Indicates whether the cluster has only one source node.
* `gcs_fc_single_primary` - A Galera parameter. Indicates whether there is more than one replication source.
* `innodb_buffer_pool_instances` - The number of regions that `innodb_buffer_pool_size` is divided into when `innodb_buffer_pool_size` > 1 GiB.
* `innodb_buffer_pool_size` - The size in bytes of the buffer pool used to cache table data and indexes.
* `innodb_change_buffering` - Operations for which change buffering optimization is enabled.
* `innodb_flush_log_at_trx_commit` - The value of the parameter controls the behaviour for transaction commit operations.
* `innodb_io_capacity` - The number of I/O operations per second (IOPS) available to InnoDB background tasks.
* `innodb_io_capacity_max` - The maximum number of IOPS that InnoDB background tasks can perform.
* `innodb_log_file_size` - The size of a single file in bytes in the redo system log.
* `innodb_log_files_in_group` - The number of system log files in a log group.
* `innodb_purge_threads` - The number of background threads allocated for the InnoDB purge operation.
* `innodb_thread_concurrency` - The maximum number of threads permitted inside of InnoDB.
* `innodb_strict_mode` - The MySQL operation mode.
* `innodb_sync_array_size` - The size of the mutex/lock wait array.
* `max_allowed_packet` - The maximum size of one packet, any generated/intermediate string
  or any parameter sent by the _mysql_stmt_send_long_data()_ C API function.
* `max_connect_errors` - The maximum number of connection errors, at which the server blocks the host from further connections.
* `max_connections` - The maximum permitted number of simultaneous client connections that a host can handle.
* `max_heap_table_size` - The maximum size in bytes, to which user-created `MEMORY` tables are permitted to grow.
* `logging` - The logging settings for the service. The structure of this block is [described below](#logging).
* `monitoring` - The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `options` - Other MySQL parameters.
* `pxc_strict_mode` - PXC mode.
* `table_open_cache` - The number of open tables for all threads.
* `thread_cache_size` - The number of threads that the server caches to establish new network connections.
* `tmp_table_size` - The maximum size of internal in-memory temporary tables in bytes.
* `transaction_isolation` - The transaction isolation level.
* `user` - List of MySQL users with parameters. The structure of this block is [described below](#mysql-user).
* `vendor` - The engine vendor.
* `version` - The installed version.
* `wait_timeout` - The number of seconds the server waits for activity on a noninteractive connection before closing it.

### MySQL database

The `database` block has the following structure:

* `backup_enabled` - Indicates whether backup is enabled for the database.
* `backup_id` - The database backup ID.
* `backup_db_name` - The name of a database from the backup specified in the `backup_id` parameter.
* `charset` - The database charset.
* `collate` - The database collation.
* `id` - The ID of the database.
* `name` - The database name.
* `user` - List of database users with parameters. The structure of this block is [described below](#mysql-database-user).

### MySQL database user

The `user` block has the following structure:

* `id` - The ID of the user.
* `name` - The MySQL user name.
* `options` - List of user options.
* `privileges` - List of user privileges.

### MySQL user

The `user` block has the following structure:

* `global_privileges` - Set of privileges on all databases, including system databases (`*.*`).
* `host` - The hostname or IP address.
* `id` - The ID of the user.
* `name` - The MySQL user name.
* `password` - The MySQL user password.

## PostgreSQL Attribute Reference

~> **Note** The following attributes contain default parameter values or user-defined values used when the service is created.

In addition to the common attributes for all services [described above](#attribute-reference),
the following attributes are exported only for a PostgreSQL service:

* `autovacuum` - Indicates whether the server must run the autovacuum launcher daemon.
* `autovacuum_max_workers` - The maximum number of autovacuum processes (other than the autovacuum launcher)
  that can run simultaneously.
* `autovacuum_vacuum_cost_delay` - The cost delay value in milliseconds used in automatic `VACUUM` operations.
* `autovacuum_vacuum_cost_limit` - The cost limit value used in automatic `VACUUM` operations.
* `autovacuum_analyze_scale_factor` - The fraction of the table size to add to `autovacuum_analyze_threshold`
  when deciding whether to trigger `ANALYZE`.
* `autovacuum_vacuum_scale_factor` - The fraction of the table size to add to `autovacuum_vacuum_threshold`
  when deciding whether to trigger `VACUUM`.
* `class` - The service class.
* `database` - List of PostgreSQL databases with parameters. The structure of this block is [described below](#postgresql-database).
* `effective_cache_size` - The planner’s assumption about the effective size of the disk cache in bytes,
  that is available to a single query.
* `effective_io_concurrency` - The number of concurrent disk I/O operations.
* `logging` - The logging settings for the service. The structure of this block is [described below](#logging).
* `maintenance_work_mem` - The maximum amount of memory in bytes, used by maintenance operations,
  such as `VACUUM`, `CREATE INDEX`, and `ALTER TABLE ADD FOREIGN KEY`.
* `max_connections` - The maximum number of simultaneous connections to the database server.
* `max_wal_size` - The maximum size in bytes, that WAL can reach at automatic checkpoints.
* `max_parallel_maintenance_workers` - The maximum number of parallel workers that a single utility command can start.
* `max_parallel_workers` - The maximum number of workers that the system can support for parallel operations.
* `max_parallel_workers_per_gather` - The maximum number of workers that a single _Gather_ node can start.
* `max_worker_processes` - The maximum number of background processes that the system can support.
* `min_wal_size` - The minimum size in bytes to shrink the WAL to. As long as WAL disk usage stays below this setting,
  old WAL files are always recycled for future use at a checkpoint, rather than removed.
* `monitoring` - The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `options` - Other PostgreSQL parameters.
* `replication_mode` - The replication mode in the _Patroni_ cluster.
* `shared_buffers` - The amount of memory in bytes, the database server uses for shared memory buffers.
* `user` - List of PostgreSQL users with parameters. The structure of this block is [described below](#postgresql-user).
* `version` - The installed version.
* `wal_buffers` - The amount of shared memory used for WAL data not yet written to a volume.
* `wal_keep_segments` - The minimum number of log files segments that must be kept in the _pg_xlog_ directory,
  in case a standby server needs to fetch them for streaming replication.
* `work_mem` - The base maximum amount of memory in bytes to be used by a query operation (such as a sort or hash table)
  before writing to temporary disk files.

### PostgreSQL database

The `database` block has the following structure:

* `backup_enabled` - Indicates whether backup is enabled for the database.
* `backup_id` - The database backup ID.
* `backup_db_name` - The name of a database from the backup specified in the `backup_id` parameter.
* `encoding` - The database encoding.
* `extensions` - List of extensions for the database.
* `id` - The ID of the database.
* `locale` - The database locale.
* `name` - The database name.
* `owner` - The name of the user who is the database owner.
* `user` - List of database users with parameters. The structure of this block is [described below](#postgresql-database-user).

### PostgreSQL database user

The `user` block has the following structure:

* `id` - The ID of the user.
* `name` - The PostgreSQL user name.

### PostgreSQL user

The `user` block has the following structure:

* `id` - The ID of the user.
* `name` - The PostgreSQL user name.
* `password` - The PostgreSQL user password when it is returned by the API. This attribute is sensitive.
* `roles` - Set of predefined PostgreSQL roles assigned to this user across the service.

## Kafka Attribute Reference

~> **Note** The following attributes contain default parameter values or user-defined values used when the service is created.

In addition to the common attributes for all services [described above](#attribute-reference),
the following attributes are exported only for a Kafka service:

* `class` - The service class.
* `logging` - The logging settings for the service. The structure of this block is [described below](#logging).
* `monitoring` - The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `options` - Other Kafka parameters.
* `version` - The installed version.

## Prometheus Attribute Reference

~> **Note** Prometheus is a single-node service, so `high_availability` is `false`.

In addition to the common attributes for all services [described above](#attribute-reference),
the following attributes are exported only for a Prometheus service:

* `class` - The service class. Its value is `monitoring`.
* `remote_write_receiver` - Indicates whether the Prometheus service accepts metrics through the Remote Write protocol.
  This attribute is available for Prometheus environment versions `paas_v4_0` and later.

## RabbitMQ Attribute Reference

~> **Note** The following attributes contain default parameter values or user-defined values used when the service is created.

In addition to the common attributes for all services [described above](#attribute-reference),
the following attributes are exported only for a RabbitMQ service:

* `class` - The service class.
* `logging` - The logging settings for the service. The structure of this block is [described below](#logging).
* `monitoring` - The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `options` - Other RabbitMQ parameters.
* `password` - The RabbitMQ admin password.
* `version` - The installed version.

## Redis Attribute Reference

~> **Note** The following attributes contain default parameter values or user-defined values used when the service is created.

In addition to the common attributes for all services [described above](#attribute-reference),
the following attributes are exported only for a Redis service:

* `class` - The service class.
* `cluster_type` - The clustering option.
* `databases` - The number of databases.
* `logging` - The logging settings for the service. The structure of this block is [described below](#logging).
* `maxmemory_policy` - The memory management mode.
* `monitoring` - The monitoring settings for the service. The structure of this block is [described below](#monitoring).
* `options` - Other Redis parameters.
* `password` - The Redis user password.
* `persistence_aof` - Indicates whether the AOF storage mode is enabled.
* `persistence_rdb` - Indicates whether the RDB storage mode is enabled.
* `timeout` - The time in seconds during which the connection to an inactive client is retained.
* `tcp_backlog` - The size of a connection queue.
* `tcp_keepalive` - The time in seconds during which the service sends ACKs to detect dead peers (unreachable clients).
* `version` - The installed version.

## Common Service Attribute Reference

### logging

The `logging` block has the following structure:

* `log_to` - The ID of the logging service.
* `logging_tags` - List of tags that are assigned to the log records of the service.

### monitoring

The `monitoring` block has the following structure:

* `monitor_by` - The ID of the monitoring service.
* `monitoring_labels` - Map containing labels that are assigned to the metrics of the service.
