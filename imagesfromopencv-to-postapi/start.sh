#!/bin/bash

: ${POST_ENDPOINT:?required}
: ${DEVICE_ID:?required}

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

python3 app.py
