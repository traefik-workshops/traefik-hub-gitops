# Traefik Hub Kubecon NA 2023

## Deploying the stack locally

Login to the Traefik Hub UI, add a new agent and copy your token.

Now, open a terminal and run these commands:


```bash
kubectl apply -k clusters/kubecon/flux-system
```

Edit the file [apps/overlays/local/token.yaml](./apps/overlays/local/token.yaml) with you base64 Hub token

```bash
kubectl apply -k apps/overlays/local
```
