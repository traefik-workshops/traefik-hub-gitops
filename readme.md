# Traefik Hub Kubecon NA 2023

## Deploying Traefik Hub

Login to the Traefik Hub UI, add a new agent and copy your token.

Now, open a terminal and run these commands:

## Install monitoring

### API management metrics with OTel

Install Prometheus with OTel ingestion support + install Grafana + add dashboards:

```shell
kubectl apply -f tools/monitoring-api-otel/monitoring.yaml
kubectl create configmap grafana-hub-dashboards \
  --from-file=tools/monitoring-api-otel/dashboard-hub.json --from-file=tools/monitoring-api-otel/hub-users.json \
  --from-file=tools/monitoring-api-otel/hub-apis.json \
  --from-file=tools/monitoring-api-otel/hub-main.json \
  -o yaml --dry-run=client -n monitoring | kubectl apply -f -
```

(Re)install the Hub agent with OTel metrics enabled and Hub exposed as an LB:

```shell
# Complete the install steps (token, etc) from above or the Hub UI, then run helm with the new parameters
# Export your token
export HUB_TOKEN=a03940a1-3540-4726-81aa-bd926ce1c5f1

# Hub prod env
helm upgrade --install --namespace traefik-hub traefik-hub traefik/traefik-hub \
  --set additionalArguments='{--hub.metrics.opentelemetry.insecure,--hub.metrics.opentelemetry.address=prometheus.monitoring:9090,--hub.metrics.opentelemetry.path=/api/v1/otlp/v1/metrics}' \
  --set service.type=LoadBalancer
```

### Ingress controller mode

```shell
# Add the Helm repository
helm repo add traefik https://traefik.github.io/charts
helm repo update

# Install the Ingress Controller
kubectl create namespace traefik-hub
kubectl create secret generic hub-agent-token --namespace traefik-hub --from-literal=token=$HUB_TOKEN
helm upgrade --install --namespace traefik-hub traefik-hub traefik/traefik-hub \
  --set additionalArguments='{--hub.metrics.opentelemetry.insecure,--hub.metrics.opentelemetry.address=prometheus.monitoring:9090,--hub.metrics.opentelemetry.path=/api/v1/otlp/v1/metrics}' \
  --set service.type=LoadBalancer
```

## Custom domain

> :warning: Traefik Labs Internal!

When using an environment from [env-on-demand](https://github.com/traefik/env-on-demand/issues "Link to GitHub repo for creating on-demand envs") you can control the DNS records to enable custom domains.

After updating the records, edit the custom domain section in `api-gateway.yaml` and `api-portal.yaml` and apply them (or reapply if did it previously).

```shell
# Add your env name and DNS token
export ENV_NAME=your-env-on-demand-cluster-name
export DNS_TOKEN=your-dns-token-from-the-env-on-demand-bot-reply

# Choose one of the ingresses
export INGRESS_CONTROLLER=hub
export INGRESS_CONTROLLER=nginx
export INGRESS_CONTROLLER=emissary
```

Run these commands for the domain setup:

```shell
case $INGRESS_CONTROLLER in
  hub)
    export SVC_NAME=traefik
    export SVC_NS=traefik-hub
  ;;
  nginx)
    export SVC_NAME=ingress-nginx-controller
    export SVC_NS=ingress-nginx
  ;;
  emissary)
    export SVC_NAME=emissary-emissary-ingress
    export SVC_NS=emissary
  ;;
esac

for subdomain in portal gateway prometheus grafana ; do
  echo ""
  echo "Creating CNAME record for ${subdomain} ..."
  curl --location --request POST "https://cf.infra.traefiklabs.tech/dns/env-on-demand?s=${subdomain}" \
  --header "X-TraefikLabs-User: ${ENV_NAME}" \
  --header "Content-Type: application/json" \
  --header "Authorization: Bearer ${DNS_TOKEN}" \
  --data-raw "{\"Value\": \"$(kubectl get service $SVC_NAME -n $SVC_NS --no-headers | awk {'print $4'})\"}"
  echo ""
  echo "Creating CNAME record for ${subdomain} ACME challenge ..."
  curl --location --request POST "https://cf.infra.traefiklabs.tech/dns/env-on-demand?s=_acme-challenge.${subdomain}" \
  --header "X-TraefikLabs-User: ${ENV_NAME}" \
  --header "Content-Type: application/json" \
  --header "Authorization: Bearer ${DNS_TOKEN}" \
  --data-raw "{\"Value\": \"acme-challenge.traefikhub.io\"}"
done
```


## Deploy APIs and Portal

Apply all resources in one step:

```shell
kubectl apply \
  -f apis/namespace.yaml \
  -f apis/customers/ \
  -f apis/employee/ \
  -f apis/flight/ \
  -f apis/ticket/ \
  -f apis/external/ \
  -f apis/api-collections.yaml \
  -f api-access.yaml \
  -f api-rate-limit.yaml \
  -f api-gateway.yaml \
  -f api-portal.yaml \
  -f custom-ui.yaml
```

## Demo playbook

If you need to do a demo on Hub, here is a [demo playbook](./demo-playbook.md) for you.

