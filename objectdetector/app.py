import os
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

inImagesService = env.get_service(name='raw-images-topic')
if inImagesService is None:
    inImagesKafka = "localhost:9092"
    inImagesTopic = "opencv-kafka-demo-raw-images"
else:
    inImagesKafka  = inImagesService.credentials.get("hostname")
    inImagesTopic  = inImagesService.credentials.get("topicName")

outImagesService = env.get_service(name='objectdetector-images-topic')
if outImagesService is None:
    outImagesKafka = "localhost:9092"
    outImagesTopic = "opencv-kafka-demo-objectdetector-images"
else:
    outImagesKafka  = outImagesService.credentials.get("hostname")
    outImagesTopic  = outImagesService.credentials.get("topicName")

import json
from kafka import KafkaProducer
statusProducer = KafkaProducer(value_serializer=lambda v: json.dumps(v).encode('utf-8'), bootstrap_servers=statusKafka)
statusProducer.send(statusTopic, {"status":"starting", "client": "objectdetector", "language": "python"})


# initialize the list of class labels MobileNet SSD was trained to
# detect, then generate a set of bounding box colors for each class
object_classes = json.load(open("models/object_classes.json", "r"))
object_colours = json.load(open("models/object_colours.json", "r"))

# load our serialized model from disk
print("[INFO] loading model...")
protext="models/MobileNetSSD_deploy.prototxt.txt"
model="models/MobileNetSSD_deploy.caffemodel"
net = cv2.dnn.readNetFromCaffe(protext, model)

CONFIDENCE=0.2 # minimum probability to filter weak detections

kafka_group_id=inImagesTopic
kafka_client_id='%s-%d' % (kafka_group_id, os.getpid())
print(kafka_group_id, kafka_client_id)

from kafka import KafkaConsumer, KafkaProducer
inImages = KafkaConsumer(inImagesTopic,
                         bootstrap_servers=inImagesKafka,
                         client_id=kafka_client_id,
                         group_id=kafka_group_id)
outImages = KafkaProducer(bootstrap_servers=outImagesKafka)

from tempfile import NamedTemporaryFile
for image in inImages:
    print(image.partition, image.offset)
    with NamedTemporaryFile(suffix=".jpg") as frame_file:
        frame_file.write(image.value)
        frame_file.seek(0)

        frame = cv2.imread(frame_file.name, cv2.IMREAD_COLOR)

    	# grab the frame dimensions and convert it to a blob
        (h, w) = frame.shape[:2]
        blob = cv2.dnn.blobFromImage(frame, 0.007843, (300, 300), 127.5)

    	# pass the blob through the network and obtain the detections and
        # predictions
        net.setInput(blob)
        detections = net.forward()

        # loop over the detections
        for i in np.arange(0, detections.shape[2]):
        	# extract the confidence (i.e., probability) associated with
        	# the prediction
        	confidence = detections[0, 0, i, 2]

        	# filter out weak detections by ensuring the `confidence` is
        	# greater than the minimum confidence
        	if confidence > CONFIDENCE:
        		# extract the index of the class label from the
        		# `detections`, then compute the (x, y)-coordinates of
        		# the bounding box for the object
        		idx = int(detections[0, 0, i, 1])
        		box = detections[0, 0, i, 3:7] * np.array([w, h, w, h])
        		(startX, startY, endX, endY) = box.astype("int")

        		# draw the prediction on the frame
        		label = "{}: {:.2f}%".format(object_classes[idx],
        			confidence * 100)
        		cv2.rectangle(frame, (startX, startY), (endX, endY),
        			object_colours[idx], 2)
        		y = startY - 15 if startY - 15 > 15 else startY + 15
        		cv2.putText(frame, label, (startX, y),
        			cv2.FONT_HERSHEY_SIMPLEX, 0.5, object_colours[idx], 2)

        image = cv2.imencode('.jpg', frame)[1].tobytes()
        outImages.send(outImagesTopic, image)

    inImages.seek_to_end()
