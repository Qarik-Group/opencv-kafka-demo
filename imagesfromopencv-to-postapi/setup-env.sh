#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

set -eu

apt-get update && apt-get install -y --no-install-recommends \
        build-essential \
        curl wget \
        libjpeg-dev \
        libpng12-dev \
        pkg-config \
        rsync \
        software-properties-common \
        unzip \
        python3 python3-dev python3-numpy \
        cmake git pkg-config \
        libjpeg-dev \
        libtiff5-dev \
        libjasper-dev \
        libpng12-dev \
        libavcodec-dev \
        libavformat-dev \
        libswscale-dev \
        libv4l-dev \
        libxvidcore-dev \
        libx264-dev \
        libgtk2.0-dev \
        libatlas-base-dev \
        gfortran

curl -O https://bootstrap.pypa.io/get-pip.py && \
    python3 get-pip.py && \
    rm get-pip.py

OPENCV_VERSION=3.3.1
cd /tmp
if [[ ! -f opencv-$OPENCV_VERSION.tar.gz ]]; then
  curl -L https://github.com/opencv/opencv/archive/$OPENCV_VERSION.tar.gz -o opencv-$OPENCV_VERSION.tar.gz
fi
rm -rf opencv-*/
tar xvzf opencv-$OPENCV_VERSION.tar.gz
cd opencv-*
    mkdir release && cd release
    cmake -D CMAKE_BUILD_TYPE=RELEASE \
          -D CMAKE_INSTALL_PREFIX=/usr/local \
          -D BUILD_PYTHON_SUPPORT=ON \
          -D BUILD_NEW_PYTHON_SUPPORT=ON \
          -D WITH_XINE=ON \
          -D WITH_TBB=ON ..
    make -j4
    make install

cd $DIR
rm -rf /tmp/opencv*/

pip3 install -r requirements.txt
