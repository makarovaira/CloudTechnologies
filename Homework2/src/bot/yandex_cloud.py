from piexif import dump
from PIL import Image

def get_image(image_path):
    with Image.open(image_path) as image:
        image.load()
    return image


def save_image(image, image_path, exif=()):
    image.save(image_path, exif=dump(exif))
