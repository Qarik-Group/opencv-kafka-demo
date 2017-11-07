#!/bin/bash

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
cd $ROOT

aws s3 sync s3://opencv-buildpack/demos/cnn-example/ models/
