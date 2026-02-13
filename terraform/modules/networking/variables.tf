# terraform/modules/networking/variables.tf
# SONA IoT Elderly Care Platform - Networking Module Variables

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
  default     = "sona"
}

variable "environment" {
  description = "Deployment environment (dev, staging, production)"
  type        = string
  default     = "production"

  validation {
    condition     = contains(["dev", "staging", "production"], var.environment)
    error_message = "Environment must be one of: dev, staging, production."
  }
}

variable "region" {
  description = "DigitalOcean region for the VPC"
  type        = string
  default     = "nyc3"
}

variable "vpc_ip_range" {
  description = "IP range for the VPC in CIDR notation"
  type        = string
  default     = "10.10.0.0/16"
}

variable "ssh_allowed_cidrs" {
  description = "List of CIDR blocks allowed to SSH into resources"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "web_tag" {
  description = "Tag applied to resources that should receive web firewall rules"
  type        = string
  default     = "sona-web"
}

variable "database_tag" {
  description = "Tag applied to resources that should receive database firewall rules"
  type        = string
  default     = "sona-database"
}

variable "k8s_tag" {
  description = "Tag applied to Kubernetes node resources"
  type        = string
  default     = "sona-k8s"
}
