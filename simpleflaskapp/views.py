import datetime
from simpleflaskapp import app
from flask import jsonify

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

