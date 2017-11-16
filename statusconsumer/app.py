from flask import Flask, jsonify, render_template, make_response
from cfenv import AppEnv

env = AppEnv()
statusService = env.get_service(name='status-topic')
if statusService is None:
    statusKafka = "localhost:9092"
    statusTopic  = "opencv-kafka-demo-status"
else:
    statusKafka  = statusService.credentials.get("hostname")
    statusTopic  = statusService.credentials.get("topicName")

from kafka import KafkaConsumer
statusConsumer = KafkaConsumer(statusTopic, bootstrap_servers=statusKafka)

metrics = statusConsumer.metrics()
print(metrics)

for msg in statusConsumer:
    print(msg)