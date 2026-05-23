terraform {
  required_version = ">= 1.7.0"
}

module "redis" {
  source = "./modules/redis"
}

module "kafka" {
  source = "./modules/kafka"
}

module "k8s" {
  source = "./modules/k8s"
}
