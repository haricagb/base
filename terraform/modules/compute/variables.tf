# terraform/modules/compute/variables.tf
# SONA IoT Elderly Care Platform - Compute Module Variables

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
  default     = "sona"
}

variable "environment" {
  description = "Deployment environment (dev, staging, production)"
  type        = string
  default     = "production"
}

variable "region" {
  description = "DigitalOcean region for the Kubernetes cluster"
  type        = string
  default     = "nyc3"
}

variable "vpc_id" {
  description = "VPC UUID to place the Kubernetes cluster in"
  type        = string
}

variable "k8s_version" {
  description = "Kubernetes version slug (e.g., '1.31.1-do.3' or 'latest')"
  type        = string
  default     = "latest"
}

variable "ha_control_plane" {
  description = "Enable high-availability control plane (recommended for production)"
  type        = bool
  default     = false
}

variable "auto_upgrade" {
  description = "Enable automatic minor version upgrades"
  type        = bool
  default     = true
}

# ---- Default Node Pool ----

variable "default_node_size" {
  description = "Droplet size for default node pool"
  type        = string
  default     = "s-2vcpu-4gb"
}

variable "default_node_count" {
  description = "Number of nodes in the default pool (used when auto-scale is disabled)"
  type        = number
  default     = 2
}

variable "auto_scale_enabled" {
  description = "Enable auto-scaling for the default node pool"
  type        = bool
  default     = true
}

variable "auto_scale_min_nodes" {
  description = "Minimum number of nodes when auto-scaling is enabled"
  type        = number
  default     = 1

  validation {
    condition     = var.auto_scale_min_nodes >= 1
    error_message = "Minimum node count must be at least 1."
  }
}

variable "auto_scale_max_nodes" {
  description = "Maximum number of nodes when auto-scaling is enabled"
  type        = number
  default     = 5

  validation {
    condition     = var.auto_scale_max_nodes >= 1
    error_message = "Maximum node count must be at least 1."
  }
}

# ---- IoT Worker Node Pool ----

variable "enable_iot_worker_pool" {
  description = "Create a dedicated node pool for IoT telemetry processing"
  type        = bool
  default     = false
}

variable "iot_worker_node_size" {
  description = "Droplet size for IoT worker node pool"
  type        = string
  default     = "s-2vcpu-4gb"
}

variable "iot_worker_min_nodes" {
  description = "Minimum nodes in the IoT worker pool"
  type        = number
  default     = 1
}

variable "iot_worker_max_nodes" {
  description = "Maximum nodes in the IoT worker pool"
  type        = number
  default     = 3
}

# ---- Metadata ----

variable "tags" {
  description = "Tags to apply to the Kubernetes cluster and node pools"
  type        = list(string)
  default     = ["sona", "kubernetes"]
}

variable "node_labels" {
  description = "Additional Kubernetes labels to apply to default pool nodes"
  type        = map(string)
  default     = {}
}

variable "maintenance_window_day" {
  description = "Day of week for maintenance window"
  type        = string
  default     = "sunday"
}

variable "maintenance_window_hour" {
  description = "Hour (UTC) for maintenance window start (HH:MM format)"
  type        = string
  default     = "04:00"
}
