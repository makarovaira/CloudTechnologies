from json import loads
from piexif import ImageIFD
from os import getenv
from yandex_cloud import get_image, save_image
from pathlib import Path
from uuid import uuid4

PHOTOS_MOUNT_POINT = getenv("PHOTOS_MOUNT_POINT")
FACES_MOUNT_POINT = getenv("FACES_MOUNT_POINT")

def handle(event, context):
    message = loads(event["messages"][0]["details"]["message"]["body"])

    object_key = message["object_key"]
    image = get_image(Path("/function/storage", PHOTOS_MOUNT_POINT, object_key))

    x, y, w, h = message["rectangle"]
    face = image.crop((x, y, x + w, y + h))

    metadata = { "0th": { ImageIFD.ImageDescription: object_key.encode("utf-8") }}
    save_image(face, Path("/function/storage", FACES_MOUNT_POINT, f"{uuid4()}.jpg"), exif=metadata)

    return { "statusCode": 200 }
