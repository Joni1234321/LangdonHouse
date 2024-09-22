import os
import rembg

def process (input_path: str, output_path: str, session = rembg.new_session()):
    if not os.path.exists(input_path) or not os.path.exists(output_path):
        return False

    with open(input_path, 'rb') as i:
        with open(output_path, 'wb') as o:
            input = i.read()
            output = rembg.remove(input, session=session)
            o.write(output)

    return True