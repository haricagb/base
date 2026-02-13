# terraform/modules/cache/main.tf
# SONA IoT Elderly Care Platform - Cache Module
# Managed Redis cluster for session management, caching, and pub/sub messaging.

resource "digitalocean_database_cluster" "redis" {
  name                 = "${var.project_name}-redis-${var.environment}"
  engine               = "redis"
  version              = var.redis_version
  size                 = var.redis_size
  region               = var.region
  node_count           = var.redis_node_count
  private_network_uuid = var.vpc_id

  maintenance_window {
    day  = var.maintenance_window_day
    hour = var.maintenance_window_hour
  }

  tags = var.tags
}

# -----------------------------------------------------------------------------
# Redis Eviction Policy Configuration
# Configures the eviction policy for the Redis cluster via database config.
# -----------------------------------------------------------------------------
resource "digitalocean_database_redis_config" "sona" {
  cluster_id            = digitalocean_database_cluster.redis.id
  maxmemory_policy      = var.redis_maxmemory_policy
  notify_keyspace_events = var.redis_notify_keyspace_events
  timeout               = var.redis_timeout
}

# -----------------------------------------------------------------------------
# Database Firewall
# Restrict Redis access to VPC and Kubernetes resources only.
# -----------------------------------------------------------------------------
resource "digitalocean_database_firewall" "redis" {
  cluster_id = digitalocean_database_cluster.redis.id

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

  # Allow additional trusted sources
  dynamic "rule" {
    for_each = var.trusted_sources
    content {
      type  = rule.value.type
      value = rule.value.value
    }
  }
}
