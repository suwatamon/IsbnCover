import numpy as np
import tensorflow as tf
from keras.models import load_model

model = load_model('tensor/cnn.h5')

req = [int(i) for i in input().split(' ')]
img = np.array(req, dtype=np.float32).reshape(1,28,28,1)
result = np.argmax(model.predict(img))

print(str(result))
