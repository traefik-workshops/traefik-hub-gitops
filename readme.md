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

This tutorial details how to deploy locally a Kubernetes cluster, with GitOps and Observability components.

# Deploy a Kubernetes cluster

TODO

# Deploy the stack on this cluster

Login to the Traefik Hub UI, add a new agent and copy your token.

Now, open a terminal and run these commands:


```bash
kubectl apply -k clusters/kubecon/flux-system
```

Edit the file [apps/overlays/local/token.yaml](./apps/overlays/local/token.yaml) with you base64 Hub token

```bash
kubectl apply -k apps/overlays/local
```
