from flask import Flask, jsonify, render_template, make_response, request
from cfenv import AppEnv

env = AppEnv()
statusService = env.get_service(name='status-topic')
if statusService is None:
    statusKafka = "localhost:9092"
    statusTopic  = "opencv-kafka-demo-status"
else:
    statusKafka  = statusService.credentials.get("hostname")
    statusTopic  = statusService.credentials.get("topicName")

imagesService = env.get_service(name='raw-images-topic')
if imagesService is None:
    imagesKafka = "localhost:9092"
    imagesTopic = "opencv-kafka-demo-raw-images"
else:
    imagesKafka  = imagesService.credentials.get("hostname")
    imagesTopic  = imagesService.credentials.get("topicName")

import json
from kafka import KafkaProducer
statusProducer = KafkaProducer(value_serializer=lambda v: json.dumps(v).encode('utf-8'), bootstrap_servers=statusKafka)
statusProducer.send(statusTopic, {"status":"starting", "client": "images-from-postapi"})

imagesService = env.get_service(name='raw-images-topic')
if imagesService is None:
    imagesKafka = "localhost:9092"
    imagesTopic = "opencv-kafka-demo-raw-images"
else:
    imagesKafka  = imagesService.credentials.get("hostname")
    imagesTopic  = imagesService.credentials.get("topicName")

imagesProducer = KafkaProducer(bootstrap_servers=imagesKafka)

app = Flask(__name__)

@app.route('/', methods = ['GET'])
def get_status():
    statusProducer.send(statusTopic, {"status":"get_status", "client": "images-from-postapi"})

    return jsonify({"statusKafka": statusKafka, "statusTopic": statusTopic, "imagesKafka": imagesKafka, "imagesTopic": imagesTopic})

@app.route('/image', methods = ['POST', 'PUT'])
def receive_image():
    image = request.data
    imagesProducer.send(imagesTopic, image)
    return '{}', 200
