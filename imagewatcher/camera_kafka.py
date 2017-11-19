from camera_base import BaseCamera
from cfenv import AppEnv
import os

lookupServiceName='objectdetector-images-topic'
if os.environ.get('IMAGE_TOPIC_SERVICE_NAME'):
    lookupServiceName = os.environ['IMAGE_TOPIC_SERVICE_NAME']

env = AppEnv()
imagesService = env.get_service(name=lookupServiceName)
if imagesService is None:
    imagesKafka = "localhost:9092"
    imagesTopic = "opencv-kafka-demo-objectdetector-images"
    if os.environ.get('IMAGE_TOPIC'):
        imagesTopic = os.environ['IMAGE_TOPIC']
else:
    imagesKafka  = imagesService.credentials.get("hostname")
    imagesTopic  = imagesService.credentials.get("topicName")

print("Loading images from Kafka topic", imagesTopic)

from kafka import KafkaConsumer
imagesConsumer = KafkaConsumer(imagesTopic, bootstrap_servers=imagesKafka)

class Camera(BaseCamera):
    @staticmethod
    def frames():
        for image in imagesConsumer:
            yield image.value
