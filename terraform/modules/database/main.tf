# terraform/modules/database/main.tf
# SONA IoT Elderly Care Platform - Database Module
# Managed PostgreSQL cluster for multi-tenant robotics management data.

resource "digitalocean_database_cluster" "postgres" {
  name                 = "${var.project_name}-pg-${var.environment}"
  engine               = "pg"
  version              = var.pg_version
  size                 = var.pg_size
  region               = var.region
  node_count           = var.pg_node_count
  private_network_uuid = var.vpc_id

  maintenance_window {
    day  = var.maintenance_window_day
    hour = var.maintenance_window_hour
  }

  tags = var.tags
}

# -----------------------------------------------------------------------------
# Application Database
# Primary database for the SONA robotics management platform.
# -----------------------------------------------------------------------------
resource "digitalocean_database_db" "app_database" {
  cluster_id = digitalocean_database_cluster.postgres.id
  name       = var.database_name
}

# -----------------------------------------------------------------------------
# Application User
# Dedicated user for application-level database access.
# -----------------------------------------------------------------------------
resource "digitalocean_database_user" "app_user" {
  cluster_id = digitalocean_database_cluster.postgres.id
  name       = var.database_user
}

# -----------------------------------------------------------------------------
# Database Firewall
# Restrict access to the database cluster to VPC and Kubernetes resources only.
# -----------------------------------------------------------------------------
resource "digitalocean_database_firewall" "postgres" {
  cluster_id = digitalocean_database_cluster.postgres.id

  # Allow access from the VPC
  rule {
    type  = "ip_addr"
    value = var.vpc_cidr
  }

  # Allow access from the Kubernetes cluster (if provided)
  dynamic "rule" {
    for_each = var.k8s_cluster_id != "" ? [var.k8s_cluster_id] : []
    content {
      type  = "k8s"
      value = rule.value
    }
  }

  # Allow additional trusted sources (e.g., CI/CD, admin IPs)
  dynamic "rule" {
    for_each = var.trusted_sources
    content {
      type  = rule.value.type
      value = rule.value.value
    }
  }
}

# -----------------------------------------------------------------------------
# Connection Pool
# PgBouncer connection pooling for efficient multi-tenant connections.
# -----------------------------------------------------------------------------
resource "digitalocean_database_connection_pool" "app_pool" {
  cluster_id = digitalocean_database_cluster.postgres.id
  name       = "${var.project_name}-pool"
  mode       = "transaction"
  size       = var.connection_pool_size
  db_name    = digitalocean_database_db.app_database.name
  user       = digitalocean_database_user.app_user.name
}
