import os, glob
import rembg

def process (input_path: str, output_path: str, session: rembg.bg.BaseSession) -> bool:
    print(f"making file at {input_path} to {output_path}")
    if not os.path.exists(input_path):
        print(f"ERROR: Process failed due to input file does not exist [{input_path}]")
        return False
        
    output_dir_path = os.path.dirname(output_path)
    if not os.path.isdir(output_dir_path):
        print(f"LOG: Process creates output directory [{output_dir_path}].")
        os.makedirs(output_dir_path)

    try:
        with open(input_path, 'rb') as i:
            input = i.read()

        output = rembg.remove(input, session=session)

        with open(output_path, 'wb') as o:
            o.write(output)

    except Exception as error:
        print(f"ERROR: Process failed due to exception [{error}]")
        return False

    return True

# Getting files
in_dir  = f"{os.path.dirname(__file__)}/test/in"
out_dir = f"{os.path.dirname(__file__)}/test/out"

if __name__ == '__main__':
    print("LOG: Started running, looking for file")
    
    in_files  =      [os.path.splitext(os.path.relpath(p, start=in_dir))  for p in glob.glob(f"{in_dir}/**/*.*",  recursive=True)]
    out_files = dict([os.path.splitext(os.path.relpath(p, start=out_dir)) for p in glob.glob(f"{out_dir}/**/*.*", recursive=True)])

    files_to_process = [(stem, ext) for (stem, ext) in in_files if stem not in out_files]

    print(f"LOG: Files to process [{len(files_to_process)}]")

    session = rembg.new_session()
    failed = 0
    for (stem, ext) in files_to_process:
        in_file  = os.path.join(in_dir,  f"{stem}{ext}")
        out_file = os.path.join(out_dir, f"{stem}.png")
        if not process(in_file, out_file, session):
            failed += 1
    
    print(f"LOG: Files succeeded  [{len(files_to_process) - failed}]")
    print(f"LOG: Files failed     [{failed}]")
    
    print('LOG: Stopping Background Process')
