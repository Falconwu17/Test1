#!/bin/bash

docker exec -i kafka_broker bash <<EOF
kafka-topics --bootstrap-server localhost:9092 \
  --create --topic entries \
  --partitions 6 --replication-factor 1 \
  --if-not-exists

kafka-topics --bootstrap-server localhost:9092 --describe --topic entries
EOF
