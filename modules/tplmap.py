import sys, os,uuid
def TPLMap(target, threads, program):
    outfile = f"programs/{program}/{uuid.uuid4()}"
    os.system(f"python3 tplmap/tplmap.py -u {target}|tee {outfile}_tplmap.txt")
