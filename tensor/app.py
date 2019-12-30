import numpy as np
import tensorflow as tf
from keras.models import load_model
from flask import Flask, request
from flask_cors import CORS

app = Flask(__name__)
CORS(app)
model = load_model('cnn.h5')
graph = tf.Graph()

@app.route('/predict',methods=['POST'])
def predict():
    req = [int(i) for i in request.form['image'].split(',')]
    img = np.array(req, dtype=np.float32).reshape(1,28,28,1)
    result = np.argmax(model.predict(img))
    return str(result)

@app.route('/')
def hello_world():
    with open('index.html', "rb") as f:
        s = f.read()
        return s
    

if __name__ == '__main__':
    app.run()
