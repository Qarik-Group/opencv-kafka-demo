from camera_base import BaseCamera
from cfenv import AppEnv
import os

cf_kafka_topic='objectdetector-images-topic'
if os.environ.get('IMAGE_TOPIC_SERVICE_NAME'):
    cf_kafka_topic = os.environ['IMAGE_TOPIC_SERVICE_NAME']

env = AppEnv()
imagesService = env.get_service(name=cf_kafka_topic)
if imagesService is None:
    imagesKafka = "localhost:9092"
    imagesTopic = "opencv-kafka-demo-objectdetector-images"
else:
    imagesKafka  = imagesService.credentials.get("hostname")
    imagesTopic  = imagesService.credentials.get("topicName")

from kafka import KafkaConsumer
imagesConsumer = KafkaConsumer(imagesTopic, bootstrap_servers=imagesKafka)

class Camera(BaseCamera):
    @staticmethod
    def frames():
        for image in imagesConsumer:
            yield image.value
