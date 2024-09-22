import os
import rembg

def process (input_path: str, output_path: str, session = rembg.new_session()):
    if not os.path.exists(input_path):
        print(f"ERROR: input file does not exist [{input_path}]")
        return False
    
    output_dir_path = os.path.dirname(output_path)
    if not os.path.isdir(output_dir_path):
        print(f"ERROR: output directory does not exist [{output_dir_path}] ")
        return False

    with open(input_path, 'rb') as i:
        with open(output_path, 'wb') as o:
            input = i.read()
            output = rembg.remove(input, session=session)
            o.write(output)

    return False