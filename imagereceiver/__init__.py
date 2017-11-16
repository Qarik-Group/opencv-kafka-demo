from flask import Flask, jsonify, render_template, make_response
from cfenv import AppEnv

env = AppEnv()
statusService = env.get_service(name='status')
if statusService is None:
    statusKafka = "localhost:9092"
    statusTopic  = "opencv-kafka-demo-status"

else:
    statusKafka  = statusService.credentials.get("hostname")
    statusTopic  = statusService.credentials.get("topicName")

imagesService = env.get_service(name='images')
if imagesService is None:
    imagesKafka = "localhost:9092"
    imagesTopic = "opencv-kafka-demo-images"

else:
    imagesKafka  = imagesService.credentials.get("hostname")
    imagesTopic  = imagesService.credentials.get("topicName")

from kafka import KafkaProducer
import json
statusProducer = KafkaProducer(value_serializer=lambda v: json.dumps(v).encode('utf-8'), bootstrap_servers=statusKafka)
statusProducer.send(statusTopic, {"status":"starting", "client": "imagereceiver"})

app = Flask(__name__)

@app.route('/', methods = ['GET'])
def get_env():
    statusProducer.send(statusTopic, {"status":"get_env", "client": "imagereceiver"})

    return jsonify({"statusKafka": statusKafka, "statusTopic": statusTopic, "imagesKafka": imagesKafka, "imagesTopic": imagesTopic})
