## Instal Prometheus Helm
```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install monitoring prometheus-community/kube-prometheus-stack
```

## Port Forward
```bash
kubectl port-forward -n chat pod/chat-db-1 9187:9187
```

## Inspect on
`http://localhost:9187/metrics`