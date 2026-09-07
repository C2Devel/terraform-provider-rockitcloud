module "ru-msk" {
  source   = "./region"
  region   = "ru-msk"
  vpc_cidr = var.vpc_cidr_blocks["ru-msk"]

  access_key = "${local.dict_creds["C2_PROJECT"]}:${local.dict_creds["BASE_ACCESS_KEY"]}"
  secret_key = local.dict_creds["EC2_SECRET_KEY"]
}

module "ru-spb" {
  source   = "./region"
  region   = "ru-spb"
  vpc_cidr = var.vpc_cidr_blocks["ru-spb"]

  access_key = "${local.dict_creds["C2_PROJECT"]}:${local.dict_creds["BASE_ACCESS_KEY"]}"
  secret_key = local.dict_creds["EC2_SECRET_KEY"]
}
