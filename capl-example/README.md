# CAPL Example

This example demonstrates how to deploy a Kubernetes cluster using Cluster API Provider Linode (CAPL) with custom networking configuration using Cilium CNI and the node-ipam-controller.

## Prerequisites

- Linode API Token with read/write access
- `kubectl` configured with access to a management cluster
- `clusterctl` installed
- Helm 3.x installed

## Configuration Overview

### Network Architecture

The cluster uses the following networking configuration:

- Pod CIDR: `10.0.0.0/8`
- Native routing mode with BGP enabled
- Host firewall enabled with policy audit mode
- Direct node-to-node routing

### Components

1. **Cilium CNI**
   - Version: 1.15.4
   - Configuration:
     - Native routing mode
     - BGP control plane enabled
     - Host firewall enabled
     - Policy audit mode
     - Direct node routing enabled
     - IPv4 masquerade enabled

2. **Node IPAM Controller**
   - Version: v0.0.1
   - Manages pod CIDR allocation
   - Runs on control plane node
   - Uses host network and kubeconfig

3. **Network Policies**
   - Default cluster-wide policy allowing:
     - Intra-cluster traffic
     - Traffic from `10.0.0.0/8`
     - API server access (port 6443)
     - Kubelet access (port 10250)

## Deployment

1. Set environment variables:
```bash
export LINODE_TOKEN="your-token-here"
export LINODE_REGION="us-east"  # or your preferred region
```

2. Deploy the cluster:
```bash
./apply-manifest.sh
```

## Verification

1. Check node status:
```bash
kubectl get nodes -o wide
```

2. Verify Cilium status:
```bash
kubectl get pods -n kube-system -l k8s-app=cilium
```

3. Check IPAM controller:
```bash
kubectl get pods -l app.kubernetes.io/name=node-ipam-controller
```

## Troubleshooting

Common issues and solutions:

1. **Node-to-Node Connectivity**
   - Ensure Cilium policies allow port 10250 for kubelet access
   - Check BGP peering status
   - Verify direct node routes are configured

2. **Pod CIDR Allocation**
   - Verify ClusterCIDR resource is created
   - Check node-ipam-controller logs
   - Ensure IPAM mode is set to "kubernetes" in Cilium

3. **Network Policy**
   - Check Cilium policy status
   - Verify policy audit mode if needed
   - Ensure correct CIDR ranges in policies

## Notes

- The node-ipam-controller must run on the control plane node for proper API server access
- Host firewall rules are managed by Cilium
- BGP is enabled for native routing between nodes
- Policy audit mode allows observing policy violations without enforcement
