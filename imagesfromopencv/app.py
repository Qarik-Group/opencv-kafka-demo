#!/usr/bin/env python
from importlib import import_module
import os
from flask import Flask, render_template, Response

from cfenv import AppEnv

env = AppEnv()
statusService = env.get_service(name='status-topic')
if statusService is None:
    statusKafka = "localhost:9092"
    statusTopic  = "opencv-kafka-demo-status"
else:
    statusKafka  = statusService.credentials.get("hostname")
    statusTopic  = statusService.credentials.get("topicName")

# FIXME: this app will be run on RaspberryPis that do not have access to Kafka
# This is temporary whilst I play
import json
from kafka import KafkaProducer
statusProducer = KafkaProducer(value_serializer=lambda v: json.dumps(v).encode('utf-8'), bootstrap_servers=statusKafka)
statusProducer.send(statusTopic, {"status":"starting", "client": "imagesfromopencv"})

imagesService = env.get_service(name='raw-images-topic')
if imagesService is None:
    imagesKafka = "localhost:9092"
    imagesTopic = "opencv-kafka-demo-images"
else:
    imagesKafka  = imagesService.credentials.get("hostname")
    imagesTopic  = imagesService.credentials.get("topicName")


imagesProducer = KafkaProducer(bootstrap_servers=imagesKafka)


# import camera driver
if os.environ.get('CAMERA'):
    Camera = import_module('camera_' + os.environ['CAMERA']).Camera
else:
    from camera import Camera

camerastream = Camera()
while True:
    frame = camerastream.get_frame()
    # print("frame")
    imagesProducer.send(imagesTopic, frame)
    # statusProducer.send(statusTopic, {"status":"image-sent", "client": "imagesfromopencv", "camera": camerastream.__class__.__name__})
