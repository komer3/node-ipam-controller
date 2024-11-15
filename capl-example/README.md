# CAPL (Cluster API Provider Linode) Example

This directory contains example manifests and scripts for deploying the node-ipam-controller in a Cluster API Provider Linode (CAPL) environment.

## Contents

- `cluster-manifest.yaml`: A comprehensive manifest that includes all the necessary components for a CAPL cluster:
  - Cilium Network Policies for cluster traffic management
  - Cluster API resources for Linode infrastructure
  - Node configurations and firewall rules
  - ClusterCIRD for IPAM controller
  - and more
  
- `apply-manifest.sh`: A helper script to apply the cluster manifest with proper environment variable substitution.

## Prerequisites

Before using these examples, ensure you have:

1. A Linode account and API token
2. kubectl installed and configured
3. A functioning CAPL management cluster with the latest release installed ([Getting Started Guide](https://linode.github.io/cluster-api-provider-linode/topics/getting-started.html))
4. The following environment variables set:
   - `LINODE_TOKEN`: Your Linode API token
   - `LINODE_REGION`: The target Linode region for deployment
   - `LINODE_SSH_KEY`: Your Linode SSH public key

## Usage

1. Set the required environment variables:
   ```bash
   export LINODE_TOKEN="your-linode-token"
   export LINODE_REGION="your-preferred-region"
   export LINODE_SSH_KEY="your-ssh-public-key" 
   ```

2. Run the apply script:
   ```bash
   ./apply-manifest.sh
   ```

The script will validate that all required environment variables are set and then apply the cluster manifest with the appropriate substitutions.

## Related Documentation

- [Node IPAM Controller Documentation](../README.md)
- [Cluster API Provider Linode (CAPL) Documentation](https://linode.github.io/cluster-api-provider-linode)
