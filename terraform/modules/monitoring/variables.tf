# terraform/modules/monitoring/variables.tf
# SONA IoT Elderly Care Platform - Monitoring Module Variables

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

variable "alerts_enabled" {
  description = "Enable or disable all monitoring alerts"
  type        = bool
  default     = true
}

# ---- Alert Notification Channels ----

variable "alert_email_addresses" {
  description = "Email addresses to receive alert notifications"
  type        = list(string)
  default     = []
}

variable "slack_webhook_url" {
  description = "Slack webhook URL for alert notifications"
  type        = string
  default     = ""
  sensitive   = true
}

variable "slack_channel" {
  description = "Slack channel for alert notifications"
  type        = string
  default     = "#sona-alerts"
}

# ---- Alert Thresholds ----

variable "cpu_threshold_percent" {
  description = "CPU usage percentage threshold to trigger alert"
  type        = number
  default     = 80

  validation {
    condition     = var.cpu_threshold_percent > 0 && var.cpu_threshold_percent <= 100
    error_message = "CPU threshold must be between 1 and 100."
  }
}

variable "memory_threshold_percent" {
  description = "Memory usage percentage threshold to trigger alert"
  type        = number
  default     = 85

  validation {
    condition     = var.memory_threshold_percent > 0 && var.memory_threshold_percent <= 100
    error_message = "Memory threshold must be between 1 and 100."
  }
}

variable "disk_threshold_percent" {
  description = "Disk usage percentage threshold to trigger alert"
  type        = number
  default     = 90

  validation {
    condition     = var.disk_threshold_percent > 0 && var.disk_threshold_percent <= 100
    error_message = "Disk threshold must be between 1 and 100."
  }
}

variable "db_cpu_threshold_percent" {
  description = "Database CPU usage percentage threshold to trigger alert"
  type        = number
  default     = 75
}

# ---- Alert Window ----

variable "alert_window" {
  description = "Time window for alert evaluation (e.g., '5m', '10m', '30m', '1h')"
  type        = string
  default     = "5m"

  validation {
    condition     = contains(["5m", "10m", "30m", "1h"], var.alert_window)
    error_message = "Alert window must be one of: 5m, 10m, 30m, 1h."
  }
}

# ---- Target Resources ----

variable "droplet_ids" {
  description = "List of Droplet IDs to monitor (leave empty to use tags)"
  type        = list(string)
  default     = []
}

variable "alert_tags" {
  description = "Tags used to target resources for compute alerts"
  type        = list(string)
  default     = ["sona", "kubernetes"]
}

# ---- Database Monitoring ----

variable "enable_database_alerts" {
  description = "Enable alerts specific to the database cluster"
  type        = bool
  default     = true
}

variable "database_droplet_ids" {
  description = "List of database Droplet IDs to monitor"
  type        = list(string)
  default     = []
}

variable "database_alert_tags" {
  description = "Tags used to target database resources for alerts"
  type        = list(string)
  default     = ["sona", "database"]
}

variable "database_health_url" {
  description = "URL for database health check endpoint (leave empty to skip uptime check)"
  type        = string
  default     = ""
}

# ---- Application Monitoring ----

variable "app_health_url" {
  description = "URL for application health check endpoint (leave empty to skip uptime check)"
  type        = string
  default     = ""
}

variable "uptime_check_regions" {
  description = "Regions to run uptime checks from"
  type        = list(string)
  default     = ["us_east", "eu_west"]
}
