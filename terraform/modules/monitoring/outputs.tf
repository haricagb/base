# terraform/modules/monitoring/outputs.tf
# SONA IoT Elderly Care Platform - Monitoring Module Outputs

output "cpu_alert_id" {
  description = "ID of the CPU usage monitoring alert"
  value       = digitalocean_monitor_alert.cpu_high.id
}

output "memory_alert_id" {
  description = "ID of the memory usage monitoring alert"
  value       = digitalocean_monitor_alert.memory_high.id
}

output "disk_alert_id" {
  description = "ID of the disk usage monitoring alert"
  value       = digitalocean_monitor_alert.disk_high.id
}

output "db_cpu_alert_id" {
  description = "ID of the database CPU monitoring alert (if enabled)"
  value       = var.enable_database_alerts ? digitalocean_monitor_alert.db_cpu_high[0].id : null
}

output "database_uptime_check_id" {
  description = "ID of the database uptime check (if configured)"
  value       = var.enable_database_alerts && var.database_health_url != "" ? digitalocean_uptime_check.database[0].id : null
}

output "app_uptime_check_id" {
  description = "ID of the application uptime check (if configured)"
  value       = var.app_health_url != "" ? digitalocean_uptime_check.app[0].id : null
}

output "alert_summary" {
  description = "Summary of configured monitoring alerts"
  value = {
    cpu_threshold    = "${var.cpu_threshold_percent}%"
    memory_threshold = "${var.memory_threshold_percent}%"
    disk_threshold   = "${var.disk_threshold_percent}%"
    db_cpu_threshold = var.enable_database_alerts ? "${var.db_cpu_threshold_percent}%" : "disabled"
    alert_window     = var.alert_window
    notifications    = {
      email_count = length(var.alert_email_addresses)
      slack       = var.slack_webhook_url != "" ? var.slack_channel : "not configured"
    }
  }
}
