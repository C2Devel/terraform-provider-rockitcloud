output "vpc_ids" {
  value = {
    "ru-msk" = module.ru-msk.vpc_id
    "ru-spb" = module.ru-spb.vpc_id
  }
}

output "subnet_ids" {
  value = {
    "ru-msk" = module.ru-msk.subnet_ids
    "ru-spb" = module.ru-spb.subnet_ids
  }
}

output "security_group_ids" {
  value = {
    "ru-msk" = module.ru-msk.security_group_id
    "ru-spb" = module.ru-spb.security_group_id
  }
}
