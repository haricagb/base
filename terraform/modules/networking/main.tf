# terraform/modules/networking/main.tf
# SONA IoT Elderly Care Platform - Networking Module
# Creates VPC and firewall rules for the platform infrastructure.

resource "digitalocean_vpc" "sona" {
  name        = "${var.project_name}-vpc-${var.environment}"
  region      = var.region
  description = "VPC for SONA ${var.environment} environment"
  ip_range    = var.vpc_ip_range
}

# -----------------------------------------------------------------------------
# Firewall: Web Traffic (HTTP/HTTPS)
# Applied to resources tagged with the web tag.
# -----------------------------------------------------------------------------
resource "digitalocean_firewall" "web" {
  name = "${var.project_name}-fw-web-${var.environment}"

  tags = [var.web_tag]

  # Inbound: Allow HTTP
  inbound_rule {
    protocol         = "tcp"
    port_range       = "80"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  # Inbound: Allow HTTPS
  inbound_rule {
    protocol         = "tcp"
    port_range       = "443"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  # Inbound: Allow SSH (restricted to trusted CIDRs)
  inbound_rule {
    protocol         = "tcp"
    port_range       = "22"
    source_addresses = var.ssh_allowed_cidrs
  }

  # Outbound: Allow all traffic
  outbound_rule {
    protocol              = "tcp"
    port_range            = "1-65535"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }

  outbound_rule {
    protocol              = "udp"
    port_range            = "1-65535"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }

  outbound_rule {
    protocol              = "icmp"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }
}

# -----------------------------------------------------------------------------
# Firewall: Database Access (PostgreSQL 5432 - VPC only)
# Restricts database port access to resources within the VPC.
# -----------------------------------------------------------------------------
resource "digitalocean_firewall" "database" {
  name = "${var.project_name}-fw-db-${var.environment}"

  tags = [var.database_tag]

  # Inbound: PostgreSQL from VPC CIDR only
  inbound_rule {
    protocol         = "tcp"
    port_range       = "5432"
    source_addresses = [digitalocean_vpc.sona.ip_range]
  }

  # Inbound: Redis from VPC CIDR only
  inbound_rule {
    protocol         = "tcp"
    port_range       = "6379"
    source_addresses = [digitalocean_vpc.sona.ip_range]
  }

  # Outbound: Allow all traffic
  outbound_rule {
    protocol              = "tcp"
    port_range            = "1-65535"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }

  outbound_rule {
    protocol              = "udp"
    port_range            = "1-65535"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }

  outbound_rule {
    protocol              = "icmp"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }
}

# -----------------------------------------------------------------------------
# Firewall: Kubernetes internal communication
# Allows node-to-node traffic within the VPC for DOKS.
# -----------------------------------------------------------------------------
resource "digitalocean_firewall" "k8s_internal" {
  name = "${var.project_name}-fw-k8s-${var.environment}"

  tags = [var.k8s_tag]

  # Inbound: All TCP from VPC (node-to-node, pod-to-pod)
  inbound_rule {
    protocol         = "tcp"
    port_range       = "1-65535"
    source_addresses = [digitalocean_vpc.sona.ip_range]
  }

  # Inbound: All UDP from VPC (DNS, overlay networking)
  inbound_rule {
    protocol         = "udp"
    port_range       = "1-65535"
    source_addresses = [digitalocean_vpc.sona.ip_range]
  }

  # Inbound: NodePort range from load balancer
  inbound_rule {
    protocol         = "tcp"
    port_range       = "30000-32767"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  # Outbound: Allow all traffic
  outbound_rule {
    protocol              = "tcp"
    port_range            = "1-65535"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }

  outbound_rule {
    protocol              = "udp"
    port_range            = "1-65535"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }

  outbound_rule {
    protocol              = "icmp"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }
}
