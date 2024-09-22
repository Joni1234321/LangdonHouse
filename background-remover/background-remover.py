import rembg

def process (input_path: str, output_path: str, session = rembg.new_session()):
    with open(input_path, 'rb') as i:
        with open(output_path, 'wb') as o:
            input = i.read()
            output = rembg.remove(input, session=session)
            o.write(output)