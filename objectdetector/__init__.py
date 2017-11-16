from flask import Flask, jsonify, render_template, make_response
from cfenv import AppEnv

env = AppEnv()
kafkaService = env.get_service(name='raw-images-topic')
if kafkaService is None:
    kafkaServers   = "localhost:9092"
    kafkaTopicName = "opencv-kafka-demo-status"

else:
    kafkaServers   = kafkaService.credentials.get("hostname")
    kafkaTopicName = kafkaService.credentials.get("topicName")

print("  hostname",  kafkaServers)
print("  topicName", kafkaTopicName)

from kafka import KafkaConsumer
consumer = KafkaConsumer(kafkaTopicName, bootstrap_servers=kafkaServers)

metrics = consumer.metrics()
print(metrics)

for msg in consumer:
    print(msg)
