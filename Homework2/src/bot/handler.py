from json import loads
from pathlib import Path
from os import getenv
import piexif
from images import add_value_to_metadata, get_face_without_name, get_face_by_tg_file_unique_id, get_originals_by_name
from telegram import send_message, send_photo
from yandex_cloud import get_image, save_image

PHOTOS_MOUNT_POINT = getenv("PHOTOS_MOUNT_POINT")
FACES_MOUNT_POINT = getenv("FACES_MOUNT_POINT")
API_GATEWAY_URL = getenv("API_GATEWAY_URL")

def handle(event, context):
    update = loads(event["body"])
    message = update.get("message")

    if message:
        _message(message)

    return { "statusCode": 200 }

def _message(message):
    if (text := message.get("text")) and text == "/start":
        pass

    elif text := message.get("text") and text == "/getface":
        object_key = get_face_without_name()
        if not object_key:
            send_message("Нет лиц с незаданным именем", message)
            return

        file_unique_id = send_photo(f"{API_GATEWAY_URL}?face={object_key}", message)

        image_path = Path("/function/storage", FACES_MOUNT_POINT, object_key)
        image = get_image(image_path)
        metadata = add_value_to_metadata(image, "0th", piexif.ImageIFD.DocumentName, file_unique_id)
        save_image(image, image_path, metadata)


    elif (text := message.get("text")) and (reply_message := message.get("reply_to_message", {})):
        file_unique_id = reply_message.get("photo", [{}])[-1].get("file_unique_id")
        if not file_unique_id:
            return

        object_key = get_face_by_tg_file_unique_id(file_unique_id)

        image_path = Path("/function/storage", FACES_MOUNT_POINT, object_key)
        image = get_image(image_path)
        metadata = add_value_to_metadata(image, "Exif", piexif.ExifIFD.UserComment, text)
        save_image(image, image_path, metadata)

    elif (text := message.get("text")) and text.startswith("/find"):
        name = text[6:]

        originals = get_originals_by_name(name)
        if not originals:
            send_message(f"Фотографии с {name} не найдены", message)
            return

        for original in originals:
            image_path = Path("/function/storage", PHOTOS_MOUNT_POINT, original)
            send_photo(image_path, message, local_file=True)

    else:
        send_message("Ошибка", message)