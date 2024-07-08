import sys
from PIL import Image
import os

def process_image(image_path):
    try:
        image = Image.open(image_path)
        image = image.convert("L")  # Convert image to grayscale
        processed_image_path = f"processed_{os.path.basename(image_path)}"
        image.save(processed_image_path)
        return processed_image_path
    except Exception as e:
        return str(e)

if __name__ == "__main__":
    image_path = sys.argv[1]
    result = process_image(image_path)
    print(result)
