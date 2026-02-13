# terraform/environments/dev/main.tf
# SONA IoT Elderly Care Platform - Dev Environment
# Multi-tenant elderly care platform on DigitalOcean.

# ---------------------------------------------------------------------------
# Terraform & Provider Configuration
# ---------------------------------------------------------------------------
terraform {
  required_version = ">= 1.0"

  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

provider "digitalocean" {
  token = var.do_token
}

# ---------------------------------------------------------------------------
# Local values shared across module calls
# ---------------------------------------------------------------------------
locals {
  project_name = var.project_name
  environment  = var.environment
  region       = var.region
}

# ---------------------------------------------------------------------------
# Module: Networking
# Creates the VPC and firewall rules used by all other modules.
# ---------------------------------------------------------------------------
module "networking" {
  source = "../../modules/networking"

  project_name = local.project_name
  environment  = local.environment
  region       = local.region
}

# ---------------------------------------------------------------------------
# Module: Database
# Managed PostgreSQL cluster for multi-tenant data storage.
# ---------------------------------------------------------------------------
module "database" {
  source = "../../modules/database"

  project_name  = local.project_name
  environment   = local.environment
  region        = local.region
  vpc_id        = module.networking.vpc_id
  db_size       = var.db_size
  db_node_count = var.db_node_count
}

# ---------------------------------------------------------------------------
# Module: Compute
# DigitalOcean Kubernetes Service (DOKS) cluster for running platform
# workloads including the IoT ingestion pipeline and care dashboards.
# ---------------------------------------------------------------------------
module "compute" {
  source = "../../modules/compute"

  project_name  = local.project_name
  environment   = local.environment
  region        = local.region
  vpc_id        = module.networking.vpc_id
  k8s_node_size = var.k8s_node_size
  k8s_min_nodes = var.k8s_min_nodes
  k8s_max_nodes = var.k8s_max_nodes
}

# ---------------------------------------------------------------------------
# Module: Storage
# DigitalOcean Spaces bucket for IoT sensor data, tenant assets, and backups.
# ---------------------------------------------------------------------------
module "storage" {
  source = "../../modules/storage"

  project_name  = local.project_name
  environment   = local.environment
  spaces_region = var.spaces_region
}

# ---------------------------------------------------------------------------
# Module: Cache
# Managed Redis cluster for session management, real-time alerts, and
# pub/sub messaging between IoT device handlers.
# ---------------------------------------------------------------------------
module "cache" {
  source = "../../modules/cache"

  project_name = local.project_name
  environment  = local.environment
  region       = local.region
  vpc_id       = module.networking.vpc_id
}

# ---------------------------------------------------------------------------
# Module: Monitoring
# Uptime checks and alerting for platform health observability.
# ---------------------------------------------------------------------------
module "monitoring" {
  source = "../../modules/monitoring"

  project_name        = local.project_name
  environment         = local.environment
  k8s_cluster_id      = module.compute.k8s_cluster_id
  database_cluster_id = module.database.database_cluster_id
}
