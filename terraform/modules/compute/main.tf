# terraform/modules/compute/main.tf
# SONA IoT Elderly Care Platform - Compute Module
# DigitalOcean Kubernetes Service (DOKS) cluster for running platform workloads.

resource "digitalocean_kubernetes_cluster" "sona" {
  name    = "${var.project_name}-k8s-${var.environment}"
  region  = var.region
  version = var.k8s_version

  vpc_uuid = var.vpc_id

  # Enable HA control plane in production
  ha = var.ha_control_plane

  # Automatic upgrades for minor/patch versions
  auto_upgrade  = var.auto_upgrade
  surge_upgrade = true

  maintenance_policy {
    day        = var.maintenance_window_day
    start_time = var.maintenance_window_hour
  }

  # -------------------------------------------------------------------------
  # Default Node Pool
  # General-purpose nodes for SONA platform services (API, web, workers).
  # -------------------------------------------------------------------------
  node_pool {
    name       = "${var.project_name}-default-pool"
    size       = var.default_node_size
    node_count = var.auto_scale_enabled ? null : var.default_node_count
    auto_scale = var.auto_scale_enabled
    min_nodes  = var.auto_scale_enabled ? var.auto_scale_min_nodes : null
    max_nodes  = var.auto_scale_enabled ? var.auto_scale_max_nodes : null

    tags   = var.tags
    labels = merge(var.node_labels, {
      "sona.io/pool" = "default"
      "sona.io/env"  = var.environment
    })
  }

  tags = var.tags
}

# -----------------------------------------------------------------------------
# Additional Node Pool: IoT Workers (optional)
# Dedicated nodes for IoT message processing and robot telemetry ingestion.
# -----------------------------------------------------------------------------
resource "digitalocean_kubernetes_node_pool" "iot_workers" {
  count = var.enable_iot_worker_pool ? 1 : 0

  cluster_id = digitalocean_kubernetes_cluster.sona.id
  name       = "${var.project_name}-iot-pool"
  size       = var.iot_worker_node_size
  auto_scale = true
  min_nodes  = var.iot_worker_min_nodes
  max_nodes  = var.iot_worker_max_nodes

  tags = var.tags
  labels = {
    "sona.io/pool"     = "iot-workers"
    "sona.io/env"      = var.environment
    "sona.io/workload" = "iot-telemetry"
  }

  taint {
    key    = "sona.io/workload"
    value  = "iot-telemetry"
    effect = "NoSchedule"
  }
}
