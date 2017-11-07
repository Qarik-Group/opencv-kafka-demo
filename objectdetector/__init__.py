from flask import Flask, jsonify, render_template, make_response
from cfenv import AppEnv

env = AppEnv()
kafka = env.get_service(name='raw-images-topic')

app = Flask(__name__)

@app.route('/', methods = ['GET'])
def get_kafka_env():
    return jsonify(kafka.credentials)
