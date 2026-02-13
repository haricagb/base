# terraform/environments/dev/variables.tf
# SONA IoT Elderly Care Platform - Dev Environment Variables

variable "do_token" {
  description = "DigitalOcean API token for provider authentication"
  type        = string
  sensitive   = true
}

variable "region" {
  description = "DigitalOcean region for infrastructure deployment"
  type        = string
  default     = "nyc1"
}

variable "project_name" {
  description = "Project name used as a prefix for all resource names"
  type        = string
  default     = "sona"
}

variable "environment" {
  description = "Deployment environment identifier (dev, staging, production)"
  type        = string
  default     = "dev"
}

variable "k8s_node_size" {
  description = "Droplet size slug for Kubernetes worker nodes"
  type        = string
  default     = "s-2vcpu-4gb"
}

variable "k8s_min_nodes" {
  description = "Minimum number of nodes in the Kubernetes auto-scale pool"
  type        = number
  default     = 1
}

variable "k8s_max_nodes" {
  description = "Maximum number of nodes in the Kubernetes auto-scale pool"
  type        = number
  default     = 3
}

variable "db_size" {
  description = "Droplet size slug for the managed PostgreSQL database cluster"
  type        = string
  default     = "db-s-1vcpu-1gb"
}

variable "db_node_count" {
  description = "Number of nodes in the managed PostgreSQL database cluster"
  type        = number
  default     = 1
}

variable "spaces_region" {
  description = "DigitalOcean region for the Spaces object storage bucket"
  type        = string
  default     = "nyc3"
}
