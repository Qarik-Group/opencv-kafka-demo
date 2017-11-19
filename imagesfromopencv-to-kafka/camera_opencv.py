import cv2
import imutils
from camera_base import BaseCamera


class Camera(BaseCamera):
    video_source = 0

    @staticmethod
    def set_video_source(source):
        Camera.video_source = source

    @staticmethod
    def frames():
        camera = cv2.VideoCapture(Camera.video_source)
        if not camera.isOpened():
            raise RuntimeError('Could not start camera.')

        while True:
            # read current frame
            _, image = camera.read()
            image = imutils.resize(image, width=400)

            # (h, w) = image.shape[:2]
            # print(h, w)

            # encode as a jpeg image and return it
            yield cv2.imencode('.jpg', image)[1].tobytes()
