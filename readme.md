<br/>

<div align="center" style="margin: 30px;">
<a href="https://hub.traefik.io/">
  <img src="https://doc.traefik.io/traefik-hub/assets/images/logos-traefik-hub-horizontal.svg"   style="width:250px;" align="center" />
</a>
<br />
<br />

<div align="center">
    <a href="https://hub.traefik.io">Log In</a> |
    <a href="https://doc.traefik.io/traefik-hub/">Documentation</a> |
    <a href="https://community.traefik.io/c/traefik-hub/20">Community</a>
</div>
</div>

<br />

<div align="center"><strong>Traefik Hub - GitOps Complete Tutorial</strong>

<br />
<br />
</div>

<div align="center">Welcome to this tutorial !</div>

# Summary

This tutorial details how to deploy locally a Kubernetes cluster, with GitOps and Observability components.

# Prerequisites

If you'd like to follow along with this tutorial on your own machine, you'll need a few things first:

1. [kubectl](https://github.com/kubernetes/kubectl) command-line tool installed and configured to access the cluster
2. [gh](https://cli.github.com/) command-line tool installed and configured with your account
3. [Flux CD](https://fluxcd.io/flux/cmd/) command-line tool installed.
4. A Kubernetes cluster running. In this tutorial, we'll use [kind](https://kind.sigs.k8s.io). You may also use alternatives like [k3d](https://k3d.io/), cloud providers or others.

kind requires some config in order to use an IngressController on localhost:

```yaml
$ cat kind.config
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 30000
    hostPort: 80
    protocol: TCP
  - containerPort: 30001
    hostPort: 443
    protocol: TCP
```

All config files are in this public GitHub repository, so you may be interested to clone it:

```shell
git clone https://github.com/traefik-workshops/traefik-hub-gitops.git
cd traefik-hub-gitops
kind create cluster --config=kind.config
kubectl cluster-info
kubectl wait --for=condition=ready nodes kind-control-plane
```

# Fork the repo and deploy Flux CD

Flux needs to be able to commit on the repository, so this tutorial can only works on a fork that _you_ owns.
You will also need a GitHub personal access token (PAT) with administration permissions.

```shell
gh repo fork --remote
kubectl apply -f clusters/kind/flux-system/gotk-components.yaml
```


# Deploy the stack on this cluster

Login to the Traefik Hub UI, add a new agent and copy your token.

Now, open a terminal and run these commands to create secret needed for Traefik Hub

```shell
export TRAEFIK_HUB_TOKEN=xxx
kubectl create namespace traefik-hub
kubectl create secret generic hub-agent-token --namespace traefik-hub --from-literal=token=${TRAEFIK_HUB_TOKEN}
```

Then, we can configure flux and launch it on the fork.

```shell
export GITHUB_ACCOUNT=xxx
export GITHUB_TOKEN=yyy

# Configure Flux CD for a repository you owned
flux create secret git git-auth  --url="https://github.com/${GITHUB_ACCOUNT}/traefik-hub-gitops" --namespace=flux-system -u git -p "${GITHUB_TOKEN}"
sed -i -e "s/traefik-workshops/${GITHUB_ACCOUNT}/g" clusters/kind/flux-system/gotk-sync.yaml

# Deploy GitRepository and Kustomization
kubectl apply -f clusters/kind/flux-system/gotk-sync.yaml
```

```bash
kubectl apply -k apps/overlays/local
```
