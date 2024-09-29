import os
import argparse
import rembg



def process(input_file: str, output_file: str, session: rembg.bg.BaseSession) -> bool:
    if not os.path.exists(input_file):
        print(f"ERROR: Process failed due to input file does not exist [{input_file}]")
        return False

    if os.path.basename(output_file) == '':
        print(f"ERROR: Process failed due to output path is not a file [{output_file}]")
        return False

    output_dir_path = os.path.dirname(output_file)
    if not os.path.isdir(output_dir_path):
        print(f"LOG: Process creates output directory [{output_dir_path}].")
        os.makedirs(output_dir_path)

    try:
        with open(input_file, 'rb') as i:
            input = i.read()

        output = rembg.remove(input, session=session)

        with open(output_file, 'wb') as o:
            o.write(output)

    except Exception as error:
        print(f"ERROR: Process failed due to exception [{error}]")
        return False

    return True

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Removes background image from files')
    parser.add_argument("i", help="path to input file")
    parser.add_argument("o", help="path to output file, this will create a new directory if not already existing")
    args = parser.parse_args()
    
    success = process(args.i, args.o, rembg.new_session())
    
    if not success:
        print(f"LOG: Failed conversion")
    else:
        print(f"LOG: Succesfull conversion")