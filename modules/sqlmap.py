import sys, os,uuid
def SQLMap(target, threads, program):
    outfile = f"programs/{program}/{uuid.uuid4()}"
    os.system(f"python3 sqlmap/sqlmap.py --batch -u {target}|tee {outfile}_sqlmap.txt")
