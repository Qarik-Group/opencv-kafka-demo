import os
import re

import cv2
import numpy as np
import imutils

from cfenv import AppEnv
env = AppEnv()
statusService = env.get_service(name='status-topic')
if statusService is None:
    statusKafka = "localhost:9092"
    statusTopic  = "opencv-kafka-demo-status"
else:
    statusKafka  = statusService.credentials.get("hostname")
    statusTopic  = statusService.credentials.get("topicName")

inImagesService = env.get_service(name=re.compile('raw'))
if inImagesService is None:
    if not os.environ.get('DEVICE_ID'):
        print("Must provide $DEVICE_ID when not running on Cloud Foundry")
        exit(1)

    inImagesKafka = "localhost:9092"
    inImagesTopic = "opencv-kafka-demo-raw-" + os.environ['DEVICE_ID']
else:
    inImagesKafka  = inImagesService.credentials.get("hostname")
    inImagesTopic  = inImagesService.credentials.get("topicName")
    print("Found inbound Cloud Foundry service", inImagesService.name)

outImagesService = env.get_service(name=re.compile('edgedetector'))
if outImagesService is None:
    outImagesKafka = "localhost:9092"
    outImagesTopic = "opencv-kafka-demo-edgedetector-" + os.environ['DEVICE_ID']
else:
    outImagesKafka  = outImagesService.credentials.get("hostname")
    outImagesTopic  = outImagesService.credentials.get("topicName")
    print("Found outbound Cloud Foundry service", outImagesService.name)

import json
from kafka import KafkaProducer
statusProducer = KafkaProducer(value_serializer=lambda v: json.dumps(v).encode('utf-8'), bootstrap_servers=statusKafka)
statusProducer.send(statusTopic, {"status":"starting", "client": "edgedetector", "language": "python"})


from kafka import KafkaConsumer, KafkaProducer
inImages = KafkaConsumer(inImagesTopic,
                         bootstrap_servers=inImagesKafka,
                         group_id=inImagesTopic)
outImages = KafkaProducer(bootstrap_servers=outImagesKafka)

from tempfile import NamedTemporaryFile
for image in inImages:
    print(image.partition, image.offset)
    with NamedTemporaryFile(suffix=".jpg") as frame_file:
        frame_file.write(image.value)
        frame_file.seek(0)

        frame = cv2.imread(frame_file.name, 0)
        edges = cv2.Canny(frame,100,200)

        image = cv2.imencode('.jpg', edges)[1].tobytes()
        outImages.send(outImagesTopic, image)

    inImages.seek_to_end()
