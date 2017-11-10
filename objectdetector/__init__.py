from flask import Flask, jsonify, render_template, make_response
from cfenv import AppEnv

env = AppEnv()
kafkaService = env.get_service(name='raw-images-topic')
if kafkaService is None:
    kafkaServers   = "localhost:9092"
    kafkaTopicName = "raw-images-topic-demo"

else:
    kafkaServers   = kafkaService.credentials.get("hostname")
    kafkaTopicName = kafkaService.credentials.get("topicName")

app = Flask(__name__)

@app.route('/', methods = ['GET'])
def get_kafka_env():
    if kafkaService is None:
        return jsonify({"hostname": kafkaServers, "topicName": kafkaTopicName})
    else:
        return jsonify(kafkaService.credentials)
