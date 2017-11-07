import datetime
from simpleflaskapp import app
from flask import jsonify, render_template, make_response
import subprocess

@app.route('/', methods = ['GET'])
def hello_world():
    return jsonify({'message': 'Welcome to simpleflaskapp!'})

@app.route('/datetime', methods = ['GET'])
def get_datetime():
    return jsonify({'datetime': str(datetime.datetime.now())})


@app.route('/date', methods = ['GET'])
def get_date():
    return jsonify({'date': str(datetime.datetime.now().date())})

@app.route('/time', methods = ['GET'])
def get_time():
    return jsonify({'time': str(datetime.datetime.now().time())})


@app.route('/aws', methods = ['GET'])
def get_aws_bucket():
    cmd = "aws s3 ls s3://opencv-buildpack/demos/cnn-example/"
    stdout = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE).stdout.read()
    resp = make_response(stdout)
    resp.headers['Content-type'] = 'text/plain; charset=utf-8'
    return resp
