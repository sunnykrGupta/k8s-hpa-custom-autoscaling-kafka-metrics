# Confluent-Kafka-Go

Automated docker build image for:

+ Alpine
+ Golang
+ Librdkafka
+ Confluent-Kafka-Go

# CONSUMER Docker run
docker run --network=host -ti --rm kafka-client:1.1 ./consumer localhost:9092 custom-topic golang-cg1 100 2000 plaintext

# PRODUCER Docker run
docker run --network=host -ti --rm kafka-client:1.1 ./producer localhost:9092 custom-topic 100 10000 none 100 3000 plaintext

