from pathlib import Path
import piexif
from os import getenv
from PIL import Image

FACES_MOUNT_POINT = getenv("FACES_MOUNT_POINT")

def add_value_to_metadata(image, idf, key, value):
    if not (exif := image.info.get("exif")):
        exif = piexif.dump(())

    exif = piexif.load(exif)
    exif[idf][key] = value.encode("utf-8")

    return exif


def get_face_without_name():
    for image_path in Path("/function/storage", FACES_MOUNT_POINT).iterdir():
        with Image.open(image_path) as image:
            image.load()

        if not _get_value_from_metadata(image, "Exif", piexif.ExifIFD.UserComment):
            return image_path.name


def get_face_by_tg_file_unique_id(file_unique_id):
    for image_path in Path("/function/storage", FACES_MOUNT_POINT).iterdir():
        with Image.open(image_path) as image:
            image.load()

        if _get_value_from_metadata(image, "0th", piexif.ImageIFD.DocumentName) == file_unique_id:
            return image_path.name


def get_originals_by_name(name):
    originals = []

    for image_path in Path("/function/storage", FACES_MOUNT_POINT).iterdir():
        with Image.open(image_path) as image:
            image.load()

        if _get_value_from_metadata(image, "Exif", piexif.ExifIFD.UserComment) == name:
            originals.append(_get_value_from_metadata(image, "0th", piexif.ImageIFD.ImageDescription))

    return originals

def _get_value_from_metadata(image, idf, key):
    if not (exif := image.info.get("exif")):
        return None

    exif = piexif.load(exif)
    value = exif[idf].get(key)

    if not value:
        return None

    return value.decode("utf-8")