import numpy as np
import json

object_classes = json.load(open("models/object_classes.json"))
colours = np.random.random_integers(0, 255, size=(len(object_classes), 3))

with open('models/object_colours.json', 'w') as outfile:
    json.dump(colours.tolist(), outfile)
