import numpy as np
from keras.models import load_model
import sys

model = load_model('script/cnn.h5')

img_input = input()
while img_input:
    req = [int(i) for i in img_input.split(' ')]
    img = np.array(req, dtype=np.float32).reshape(1,28,28,1)
    result = np.argmax(model.predict(img))
    print(str(result))
    img_input = input()
