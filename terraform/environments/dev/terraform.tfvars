# terraform/environments/dev/terraform.tfvars
# SONA IoT Elderly Care Platform - Dev Environment Variable Values
#
# NOTE: do_token should NEVER be committed to version control.
# Set it via an environment variable instead:
#
#   export TF_VAR_do_token="<your-digitalocean-api-token>"

region       = "nyc1"
project_name = "sona"
environment  = "dev"

# Kubernetes cluster sizing (smaller for dev to reduce cost)
k8s_node_size = "s-2vcpu-2gb"
k8s_min_nodes = 1
k8s_max_nodes = 2

# Database cluster sizing
db_size       = "db-s-1vcpu-1gb"
db_node_count = 1
