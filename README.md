
## Kubernetes HPA Autoscaling With Kafka Metrics

Repo contains :

+ go-kafka  # contains kafka golang client
+ kafka-exporter    # kafka-exporter manifest details
+ consumer-kafka-client-deployment.yaml # consumer deployment manifest
+ producer-kafka-client-deployment.yaml # producer deployment manifest
+ kafka-custom-metrics-hpa.yaml     # hpa manifest configured with external metrics
