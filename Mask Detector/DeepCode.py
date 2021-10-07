
#Data Preprocessing

#--------------------- Importing the Libraries----------------------------- 
from tensorflow.keras.preprocessing.image import ImageDataGenerator
from tensorflow.keras.applications import MobileNetV2
from tensorflow.keras.layers import AveragePooling2D
from tensorflow.keras.layers import Dropout
from tensorflow.keras.layers import Flatten
from tensorflow.keras.layers import Dense
from tensorflow.keras.layers import Input
from tensorflow.keras.models import Model
from tensorflow.keras.optimizers import Adam
from tensorflow.keras.applications.mobilenet_v2 import preprocess_input
from tensorflow.keras.preprocessing.image import img_to_array
from tensorflow.keras.preprocessing.image import load_img
from tensorflow.keras.utils import to_categorical
from sklearn.preprocessing import LabelBinarizer
from sklearn.model_selection import train_test_split
from sklearn.metrics import classification_report
from imutils import paths
import matplotlib.pyplot as plt
import numpy as np
import os

# Now we will create the directory in which our data is stored by the name Directory
DIC = r"/content/drive/MyDrive/Mask_Detect/data"
# Now we will create the categories as there are two categories 
# With Mask and the Without Mask
CAT = ["with_mask","without_mask"]

print("Loading the images...")

# We will now have two lists the data (i.e Images) and the labels (i.e Categories) later on we will append
data = []
labels = []
# We are here looping through the categories in directory which are with and without mask and joining to form the path
# Then loop through all the path and join path and corresponding image.
# Then we load the image of the particular path at size (224,224)
for category in CAT:
	path = os.path.join(DIC,category) # with_mask or without_mask will get join 
	for img in os.listdir(path):
		img_path = os.path.join(path,img) # each image in with_mas or without_mask will get join
		image = load_img(img_path, target_size=(224,224))
		image = img_to_array(image) #converts to array (As deeplearning only works with array)
		image = preprocess_input(image) # If you use MobileNet you need to use this 
		# Here we append
		data.append(image)
		labels.append(category)

# Converting data and labels to numpy arrays

data = np.array(data,dtype='float32')
labels = np.array(labels)

# Splitting the images to train and test (Splits)

(trainx, testx, trainy, testy) = train_test_split(data, labels, test_size=0.20, stratify=labels, random_state=42)

# Data Augmentation (Create more data in the memory)

aug = ImageDataGenerator(
	       rotation_range=20,
	       zoom_range=0.15,
	       width_shift_range=0.2,
	       height_shift_range=0.2,
	       shear_range=0.15,
	       horizontal_flip=True,
	       fill_mode='nearest' )

Lr = 1e-4
epochs = 20
batch = 32

# There are two type of models here I have made the top model and bottom model
# Top Model is basically the top layers made by me
# Bottom Model are the bottom layes of MobileNetv2
# Let's start the make of model now using Functional API

bottomModel = MobileNetV2(weights="imagenet", include_top=False, input_tensor=Input(shape=(224,224,3))) 

topModel = bottomModel.output
topModel = AveragePooling2D(pool_size=(7,7))(topModel)
topModel = Flatten()(topModel)
topModel = Dense(128, activation='relu')(topModel)
topModel = Dropout(0.5)(topModel)# To Avoid Overfitting
topModel = Dense(2, activation='softmax')(topModel)

# Connecting both top and bottom model

model = Model(inputs=bottomModel.input, outputs=topModel)


# Now we will loop over all the layers in bottom model and freeze so they won't be updated

for layer in bottomModel.layers:
	layer.trainable = False
# Compile the model

print("Compiling...")

opt = Adam(lr=Lr, decay=Lr/epochs)
model.compile(loss='binary_crossentropy', optimizer=opt, metrics=['accuracy'])

# train the head
print("Training the Head...")

history = model.fit(aug.flow(trainx, trainy, batch_size=batch), 
	steps_per_epoch=len(trainx)//batch,
    validation_data=(testx, testy),
    validation_steps=len(testx)//batch,
    epochs=epochs)


#prediction on testing set

print("Evaluating Network...")

predIdxs = model.predict(testx, batch_size=batch)

# For image intest set we need to find the index of the label so inorder to do that we are using argmax
# This will give us labels with largest predicted probablity
predIdxs = np.argmax(predIdxs, axis=1)

# Saving the model

print("Saving...")
model.save('maskDetector', save_format='h5')


#Plotting the training loss and Accuracy
N = epochs
plt.style.use('ggplot')
plt.figure()
plt.plot(np.arange(0,N), history.history['loss'], label='train_loss')
plt.plot(np.arange(0,N), history.history['val_loss'], label='val_loss')
plt.plot(np.arange(0,N), history.history['accuracy'], label='train_acc')
plt.plot(np.arange(0,N), history.history['val_accuracy'], label='val_acc')
plt.title('Training Loass and Accuracy')
plt.xlabel('Epochs')
plt.ylabel('Loss-Accuracy')
plt.legend(loc='lower left')
plt.savefig('plot.png')



# ------------------------------------ONE-HOT-Encoding--------------------------------

# We have labels as with mask or without mask so it's not useful for us so in order to convert we will use the one hot encoding
lb = LabelBinarizer()
labels = lb.fit_transform(labels)
labels = to_categorical(labels)
