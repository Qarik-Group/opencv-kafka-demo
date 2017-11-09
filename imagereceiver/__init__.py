from flask import Flask, jsonify, render_template, make_response
from cfenv import AppEnv

env = AppEnv()
kafka = env.get_service(name='raw-images-topic')

kafkaServers   = kafka.credentials.get("hostname")
kafkaTopicName = kafka.credentials.get("topicName")
print("  hostname",  kafkaServers)
print("  topicName", kafkaTopicName)

from kafka import KafkaProducer
producer = KafkaProducer(bootstrap_servers=kafkaServers)
producer.send(kafkaTopicName, key="status", value="starting")

app = Flask(__name__)

@app.route('/', methods = ['GET'])
def get_kafka_env():
    return jsonify(kafka.credentials)
