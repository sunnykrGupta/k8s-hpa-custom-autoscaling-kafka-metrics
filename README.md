
## Kubernetes HPA Autoscaling With Kafka Metrics

Blog post: https://medium.com/google-cloud/kubernetes-hpa-autoscaling-with-kafka-metrics-88a671497f07

---

### Repo contains:

```
go-kafka  # contains kafka golang client
kafka-exporter    # kafka-exporter manifest details
consumer-kafka-client-deployment.yaml # consumer deployment manifest
producer-kafka-client-deployment.yaml # producer deployment manifest
kafka-custom-metrics-hpa.yaml     # hpa manifest configured with external metrics
```
