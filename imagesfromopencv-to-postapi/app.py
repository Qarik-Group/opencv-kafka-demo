#!/usr/bin/env python
from importlib import import_module
import os
import requests
from time import sleep

if not os.environ.get('POST_ENDPOINT'):
    print("Requires $POST_ENDPOINT")
    exit(1)

postEndpoint = os.environ['POST_ENDPOINT']

# import camera driver
if os.environ.get('CAMERA'):
    print("Using camera", os.environ['CAMERA'])
    Camera = import_module('camera_' + os.environ['CAMERA']).Camera
else:
    print("Using demo camera")
    from camera import Camera

camerastream = Camera()
while True:
    frame = camerastream.get_frame()
    try:
        r = requests.post(postEndpoint, data=frame)
        if r.status_code != 200:
            print(r.status_code)
            sleep(1)
    except requests.exceptions.RequestException as err:
        print("Connection error: {0}".format(err))
        sleep(1)
