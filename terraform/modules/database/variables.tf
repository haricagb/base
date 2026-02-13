# terraform/modules/database/variables.tf
# SONA IoT Elderly Care Platform - Database Module Variables

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
  description = "DigitalOcean region for the database cluster"
  type        = string
  default     = "nyc3"
}

variable "vpc_id" {
  description = "VPC UUID to place the database cluster in"
  type        = string
}

variable "vpc_cidr" {
  description = "VPC CIDR range for database firewall rules"
  type        = string
}

variable "pg_version" {
  description = "PostgreSQL engine version"
  type        = string
  default     = "16"
}

variable "pg_size" {
  description = "Database droplet size slug"
  type        = string
  default     = "db-s-1vcpu-1gb"
}

variable "pg_node_count" {
  description = "Number of nodes in the database cluster"
  type        = number
  default     = 1

  validation {
    condition     = var.pg_node_count >= 1 && var.pg_node_count <= 3
    error_message = "Node count must be between 1 and 3."
  }
}

variable "database_name" {
  description = "Name of the application database"
  type        = string
  default     = "robotics_mgmt"
}

variable "database_user" {
  description = "Name of the application database user"
  type        = string
  default     = "app_user"
}

variable "k8s_cluster_id" {
  description = "Kubernetes cluster ID to allow database access (leave empty to skip)"
  type        = string
  default     = ""
}

variable "trusted_sources" {
  description = "Additional trusted sources for database firewall rules"
  type = list(object({
    type  = string
    value = string
  }))
  default = []
}

variable "connection_pool_size" {
  description = "Number of connections in the PgBouncer connection pool"
  type        = number
  default     = 20
}

variable "maintenance_window_day" {
  description = "Day of week for maintenance window (monday, tuesday, etc.)"
  type        = string
  default     = "sunday"
}

variable "maintenance_window_hour" {
  description = "Hour (UTC) for maintenance window start (00-23)"
  type        = string
  default     = "04:00"
}

variable "tags" {
  description = "Tags to apply to the database cluster"
  type        = list(string)
  default     = ["sona", "database", "postgresql"]
}
