from camera_base import BaseCamera
from cfenv import AppEnv

env = AppEnv()
imagesService = env.get_service(name='objectdetector-images-topic')
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
