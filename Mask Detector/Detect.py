# Here we will detect the mask with the video stream
# Importing Pacakges 

from tensorflow.keras.applications.mobilenet_v2 import preprocess_input
from tensorflow.keras.preprocessing.image import img_to_array
from tensorflow.keras.models import load_model
from imutils.video import VideoStream
import numpy as np
import imutils
import time
import cv2
import os

def detect(frame, faceNet, maskNet):
	#Grabbing thedimentions of frame and constucting a blob
	(h,w) = frame.shape[:2]
	blob = cv2.dnn.blobFromImage(frame, 1.0, (300,300),
		                         (104.0, 177.0, 123.0))

	#Now we will pass the blob through network and perform face detection
	faceNet.setInput(blob)
	detections = faceNet.forward()
	print(detections.shape)

	# Now we will initialize list of faces, list of locations and their predictions
	faces = []
	locs = []
	preds = []


	# Now we will loop over the detections
	for i in range(0, detections.shape[2]):
		# Here confidence means probablity, We will extract the probabity from the detections.
		confidence = detections[0, 0, i, 2]

		# Now we will filter out the weak detections and and ensure that probablity is greater than the threshold

		if confidence > 0.5:
			# Here we will compute the (x,y) co-ordinates
			#Hence make the box around the face
			box = detections[0, 0, i, 3:7] * np.array([w, h, w, h]) 
			(startx, starty, endx, endy) = box.astype('int')

			# To make sure the bounding box or Object localization fall in dimension

			(startx, starty) = (max(0,startx), max(0,starty))
			(endx, endy) = (min(w-1, endx), min(h-1, endy))


			# Now we will extract the face and convert to RGB from BGR
			# Also Resizing to 224 by 224 then process it

			face = frame[starty:endy, startx:endx]
			face = cv2.cvtColor(face, cv2.COLOR_BGR2RGB)
			face = cv2.resize(face, (224,224))
			face = img_to_array(face)
                        #face = np.expand_dims(face,axis=0)
			face = preprocess_input(face)

			# Append to the list

			faces.append(face)
			locs.append((startx, starty, endx, endy))


	# Now in order to make predictions if atleast 1 face is detected
	# We will use lines of code below
	if len(faces) > 0:
		faces = np.array(faces, dtype='float32')
		preds = maskNet.predict(faces, batch_size=32)

	#Returning a 2-tuple of faces and locations
	return(locs, preds)


# Loading face detector model from disk
prototxt = r"C:\Users\Adin\Desktop\Mask_Detect\deploy.prototx"
weightsres = r"C:\Users\Adin\Desktop\Mask_Detect\res10_300x300_ssd_iter_140000.caffemodel"
faceNet = cv2.dnn.readNet(prototxt, weightsres)


# Loading the Mask Detector Model from disk

maskNet = load_model("maskDetector.model")

print(" Starting Video Stream")

vs = VideoStream(src=0).start() # src is basically to give index to which camera


# Now we will loop over frames from video stream

while True:
	#Grabbing the frame and resizing to have max width
	frame = vs.read()
	frame = imutils.resize(frame, width=400)

	# Detect Faces and determine if mask is there or not
	(locs, preds) = detect(frame, faceNet, maskNet)

	# loop And detect face and mark box

	for (box, pred) in zip(locs, preds):
		#Unpack the boundaries
		(startx, starty, endx, endy) = box
		(mask, withoutMask) = pred

		# Determine the label
		label = "Mask" if mask > withoutMask else "No Mask"
		color = (0,255,0) if label == "Mask" else (0,0,255)

		# Include Probablity
		label = "{}: {:.2f}%".format(label, max(mask, withoutMask)*100)

		# Displaying 

		cv2.putText(frame, label, (startx, starty-10),
			 cv2.FONT_HERSHEY_SIMPLEX, 0.45, color, 2)
		cv2.rectangle(frame, (startx, starty), (endx, endy), color, 2)

	# Show the output
	
	cv2.imshow("Frame", frame)
	key = cv2.waitKey(1) & 0xFF

	# To quit we will pass q

	if key == ord("q"):
		break

# Clean Up
cv2.destroyAllWindows()
vs.stop()



