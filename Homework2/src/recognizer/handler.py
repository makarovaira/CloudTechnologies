import cv2
from pathlib import Path
from os import getenv
from json import dumps
from boto3 import client

MOUNT_POINT = getenv("MOUNT_POINT")
QUEUE_URL = getenv("QUEUE_URL")
ACCESS_KEY_ID = getenv("ACCESS_KEY_ID")
SECRET_ACCESS_KEY = getenv("SECRET_ACCESS_KEY")

def handle(event, context):
    object_key = event["messages"][0]["details"]["object_id"]

    for face in _recognize_faces(Path("/function/storage", MOUNT_POINT, object_key)):
        client(
            service_name="sqs",
            endpoint_url="https://message-queue.api.cloud.yandex.net",
            region_name="ru-central1",
            aws_access_key_id=ACCESS_KEY_ID,
            aws_secret_access_key=SECRET_ACCESS_KEY,
        ).send_message(
            QueueUrl=QUEUE_URL,
            MessageBody=dumps({"object_key": object_key, "rectangle": face})
        )

    return { "statusCode": 200 }

def _recognize_faces(image_path):
    image = cv2.cvtColor(cv2.imread(image_path), cv2.COLOR_BGR2GRAY)
    face_cascade = cv2.CascadeClassifier(f"{cv2.data.haarcascades}haarcascade_frontalface_default.xml")

    faces = face_cascade.detectMultiScale(image, 1.2, 5)
    return map(lambda face: list(map(int, face)), faces)