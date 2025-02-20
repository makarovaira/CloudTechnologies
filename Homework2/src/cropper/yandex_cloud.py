from requests import put
from time import sleep
from piexif import dump
from os import getenv
from PIL import Image
from requests_aws4auth import AWS4Auth

FACES_MOUNT_POINT = getenv("FACES_MOUNT_POINT")
ACCESS_KEY_ID = getenv("ACCESS_KEY_ID")
SECRET_ACCESS_KEY = getenv("SECRET_ACCESS_KEY")

def get_image(image_path):
    with Image.open(image_path) as image:
        image.load()
    return image

def save_image(image, image_path, exif):
    image.save(image_path, exif=dump(exif))
    sleep(1)

    put(
        f"https://storage.yandexcloud.net/{FACES_MOUNT_POINT}/{image_path.name}",
        headers = {
            "Content-Type": "image/jpeg",
            "X-Amz-Copy-Source": f"/{FACES_MOUNT_POINT}/{image_path.name}",
            "X-Amz-Metadata-Directive": "REPLACE",
        },
        auth = AWS4Auth(
            ACCESS_KEY_ID,
            SECRET_ACCESS_KEY,
            "ru-central1",
            "s3",
        )
    )