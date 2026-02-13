# terraform/modules/cache/variables.tf
# SONA IoT Elderly Care Platform - Cache Module Variables

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
  description = "DigitalOcean region for the Redis cluster"
  type        = string
  default     = "nyc3"
}

variable "vpc_id" {
  description = "VPC UUID to place the Redis cluster in"
  type        = string
}

variable "vpc_cidr" {
  description = "VPC CIDR range for Redis firewall rules"
  type        = string
}

variable "redis_version" {
  description = "Redis engine version"
  type        = string
  default     = "7"
}

variable "redis_size" {
  description = "Redis droplet size slug"
  type        = string
  default     = "db-s-1vcpu-1gb"
}

variable "redis_node_count" {
  description = "Number of nodes in the Redis cluster"
  type        = number
  default     = 1

  validation {
    condition     = var.redis_node_count >= 1 && var.redis_node_count <= 3
    error_message = "Node count must be between 1 and 3."
  }
}

variable "redis_maxmemory_policy" {
  description = "Redis eviction policy when max memory is reached"
  type        = string
  default     = "allkeys-lru"

  validation {
    condition = contains([
      "noeviction",
      "allkeys-lru",
      "allkeys-random",
      "volatile-lru",
      "volatile-random",
      "volatile-ttl",
    ], var.redis_maxmemory_policy)
    error_message = "Invalid Redis maxmemory policy."
  }
}

variable "redis_notify_keyspace_events" {
  description = "Redis keyspace notification events (empty string to disable)"
  type        = string
  default     = ""
}

variable "redis_timeout" {
  description = "Redis client idle timeout in seconds (0 to disable)"
  type        = number
  default     = 0
}

variable "k8s_cluster_id" {
  description = "Kubernetes cluster ID to allow Redis access (leave empty to skip)"
  type        = string
  default     = ""
}

variable "trusted_sources" {
  description = "Additional trusted sources for Redis firewall rules"
  type = list(object({
    type  = string
    value = string
  }))
  default = []
}

variable "maintenance_window_day" {
  description = "Day of week for maintenance window (monday, tuesday, etc.)"
  type        = string
  default     = "sunday"
}

variable "maintenance_window_hour" {
  description = "Hour (UTC) for maintenance window start (00-23)"
  type        = string
  default     = "05:00"
}

variable "tags" {
  description = "Tags to apply to the Redis cluster"
  type        = list(string)
  default     = ["sona", "cache", "redis"]
}
