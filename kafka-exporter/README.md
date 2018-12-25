

## Kafka exporter

- [https://github.com/danielqsj/kafka_exporter](https://github.com/danielqsj/kafka_exporter)

```
# local run
$ docker run --network=host -ti --rm -p 9308:9308 danielqsj/kafka-exporter --kafka.server=localhost:9092 --kafka.version=1.0.1

# for multiple brokers
$ docker run --network=host -ti --rm -p 9308:9308 danielqsj/kafka-exporter --kafka.server=server-1:9092 --kafka.server=server-2:9192 --kafka.version=1.0.1
```


## Prometheus-to-sd

- https://github.com/GoogleCloudPlatform/k8s-stackdriver/tree/master/prometheus-to-sd
