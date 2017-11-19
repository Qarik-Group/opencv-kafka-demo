#!/bin/bash

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
cd $ROOT

: ${AWS_ACCESS_KEY_ID:?required}
: ${AWS_SECRET_ACCESS_KEY:?required}

aws s3 sync s3://opencv-buildpack/demos/cnn-example/ models/
