## Confluent-Kafka-Go

Automated docker build image for:

+ Alpine
+ Golang
+ Librdkafka
+ Confluent-Kafka-Go
+ consumer.go   # client consumes configured topic
+ producer.go   # client produce sample data to configured topic

## Build docker image
```
docker build -t kafka-client:1.1 .
```

## Golang Kafka Consumer
```
docker run --network=host -ti --rm kafka-client:1.1 ./consumer localhost:9092 custom-topic golang-cg1 100 2000 plaintext
```

## Golang Kafka Producer
```
docker run --network=host -ti --rm kafka-client:1.1 ./producer localhost:9092 custom-topic 100 10000 none 100 3000 plaintext
```

#### Repo contains
- consumer.go     # to consume topic and display reading stats
- producer.go   # to produce same packet repetitively and display stats
- sample-json.txt   # sample data packet which will be produced by producer.go



#### Producer 500,000 records (Producer Client)

> Usage: ./producer <broker:port> <topic> <msgBurst> <total> <none/gzip> <lingerMs> <waitMs> <ssl/plaintext/sasl_plaintext/sasl_ssl>
 Additional Params needed : <ssl.ca.location> <ssl.keystore.location> <ssl.keystore.password>

- broker:port : all Kafka Broker details
- topic : Topic where you want to producer messages
- msgBurst : To display stats after producing these many messages
- total : Total message to produce
- none/gzip : Compression codec
- lingerMs : Time to burst producer messages
- waitMs : Wait time when internal queue is full.
- ssl/plaintext : Choose protocol to communicate.
- ssl.ca.location : position of CA certs, which kafka servers trust
- ssl.keystore.location : key-pairs signed with CA certs
- ssl.keystore.password : Password to read key-pairs


#### Consume message indefinitely (Consumer Client)

> Usage: ./consumer <broker:port> <topic> <groupID> <groupNumber> <waitMs> <ssl/plaintext>
 Additional Params needed : <ssl.ca.location> <ssl.keystore.location> <ssl.keystore.password>

- broker:port : all Kafka Broker details
- topic : Topic where you want to producer messages
- groupID : Consumer Group Name
- groupNumber : To display stats after reading these many messages
- waitMs: Wait after reading some batch of consumers, keep it 0 for no wait.
- ssl/plaintext : Choose protocol to communicate. If using 'ssl', configure ssl parameters in code. Refer Kafka Documentation to generate certificate.
- ssl.ca.location : position of CA certs, which kafka servers trust
- ssl.keystore.location : key-pairs signed with CA certs
- ssl.keystore.password : Password to read key-pairs
