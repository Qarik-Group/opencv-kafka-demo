

On Raspberry Pi:

```
CAMERA=
POST_ENDPOINT=http://drnic-laptop.local/image
DEVICE_ID=drnic-pi
docker rm -f $DEVICE_ID; \
docker run -d \
  --name $DEVICE_ID \
  --device=/dev/vchiq --device=/dev/vcsm \
  -e POST_ENDPOINT=$POST_ENDPOINT \
  -e DEVICE_ID=$DEVICE_ID \
  -e CAMERA=$CAMERA \
  starkandwayne/imagesfromopencv-to-postapi:armv7 \
  /app/start.sh; \
docker logs $DEVICE_ID -f
```
