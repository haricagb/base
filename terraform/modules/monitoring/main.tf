# terraform/modules/monitoring/main.tf
# SONA IoT Elderly Care Platform - Monitoring Module
# DigitalOcean monitoring alerts for infrastructure health.

# -----------------------------------------------------------------------------
# High CPU Usage Alert
# Triggers when CPU utilization exceeds the threshold on k8s nodes / droplets.
# -----------------------------------------------------------------------------
resource "digitalocean_monitor_alert" "cpu_high" {
  alerts {
    email = var.alert_email_addresses
    slack {
      channel = var.slack_channel
      url     = var.slack_webhook_url
    }
  }

  window      = var.alert_window
  type        = "v1/insights/droplet/cpu"
  compare     = "GreaterThan"
  value       = var.cpu_threshold_percent
  enabled     = var.alerts_enabled
  entities    = var.droplet_ids
  description = "SONA ${var.environment}: CPU usage exceeds ${var.cpu_threshold_percent}% on monitored nodes"
  tags        = var.alert_tags
}

# -----------------------------------------------------------------------------
# High Memory Usage Alert
# Triggers when memory utilization exceeds the threshold.
# -----------------------------------------------------------------------------
resource "digitalocean_monitor_alert" "memory_high" {
  alerts {
    email = var.alert_email_addresses
    slack {
      channel = var.slack_channel
      url     = var.slack_webhook_url
    }
  }

  window      = var.alert_window
  type        = "v1/insights/droplet/memory_utilization_percent"
  compare     = "GreaterThan"
  value       = var.memory_threshold_percent
  enabled     = var.alerts_enabled
  entities    = var.droplet_ids
  description = "SONA ${var.environment}: Memory usage exceeds ${var.memory_threshold_percent}% on monitored nodes"
  tags        = var.alert_tags
}

# -----------------------------------------------------------------------------
# High Disk Usage Alert
# Triggers when disk utilization exceeds the threshold.
# Critical for preventing data loss on nodes handling IoT telemetry.
# -----------------------------------------------------------------------------
resource "digitalocean_monitor_alert" "disk_high" {
  alerts {
    email = var.alert_email_addresses
    slack {
      channel = var.slack_channel
      url     = var.slack_webhook_url
    }
  }

  window      = var.alert_window
  type        = "v1/insights/droplet/disk_utilization_percent"
  compare     = "GreaterThan"
  value       = var.disk_threshold_percent
  enabled     = var.alerts_enabled
  entities    = var.droplet_ids
  description = "SONA ${var.environment}: Disk usage exceeds ${var.disk_threshold_percent}% on monitored nodes"
  tags        = var.alert_tags
}

# -----------------------------------------------------------------------------
# Database CPU Alert
# Monitors the managed PostgreSQL cluster CPU usage.
# -----------------------------------------------------------------------------
resource "digitalocean_monitor_alert" "db_cpu_high" {
  count = var.enable_database_alerts ? 1 : 0

  alerts {
    email = var.alert_email_addresses
    slack {
      channel = var.slack_channel
      url     = var.slack_webhook_url
    }
  }

  window      = var.alert_window
  type        = "v1/insights/droplet/cpu"
  compare     = "GreaterThan"
  value       = var.db_cpu_threshold_percent
  enabled     = var.alerts_enabled
  entities    = var.database_droplet_ids
  description = "SONA ${var.environment}: Database CPU usage exceeds ${var.db_cpu_threshold_percent}%"
  tags        = var.database_alert_tags
}

# -----------------------------------------------------------------------------
# Database Connections Alert (via uptime check)
# Monitors the database for connectivity issues.
# Uses a DO Uptime Check to verify the database endpoint is reachable.
# -----------------------------------------------------------------------------
resource "digitalocean_uptime_check" "database" {
  count = var.enable_database_alerts && var.database_health_url != "" ? 1 : 0

  name    = "${var.project_name}-db-health-${var.environment}"
  target  = var.database_health_url
  type    = "https"
  regions = var.uptime_check_regions
  enabled = var.alerts_enabled
}

resource "digitalocean_uptime_alert" "database_down" {
  count = var.enable_database_alerts && var.database_health_url != "" ? 1 : 0

  name       = "${var.project_name}-db-down-${var.environment}"
  check_id   = digitalocean_uptime_check.database[0].id
  type       = "down"
  period     = "2m"
  comparison = "less_than"
  threshold  = 1

  notifications {
    email = var.alert_email_addresses
    slack {
      channel = var.slack_channel
      url     = var.slack_webhook_url
    }
  }
}

# -----------------------------------------------------------------------------
# Application Uptime Check
# Monitors the SONA platform API endpoint for availability.
# -----------------------------------------------------------------------------
resource "digitalocean_uptime_check" "app" {
  count = var.app_health_url != "" ? 1 : 0

  name    = "${var.project_name}-app-health-${var.environment}"
  target  = var.app_health_url
  type    = "https"
  regions = var.uptime_check_regions
  enabled = var.alerts_enabled
}

resource "digitalocean_uptime_alert" "app_down" {
  count = var.app_health_url != "" ? 1 : 0

  name       = "${var.project_name}-app-down-${var.environment}"
  check_id   = digitalocean_uptime_check.app[0].id
  type       = "down"
  period     = "2m"
  comparison = "less_than"
  threshold  = 1

  notifications {
    email = var.alert_email_addresses
    slack {
      channel = var.slack_channel
      url     = var.slack_webhook_url
    }
  }
}
