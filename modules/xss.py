import sys, os,uuid
def XSStrike(target, threads, program):
    outfile = f"programs/{program}/{uuid.uuid4()}"
    os.system(f"python3 XSStrike/xsstrike.py -u {target} -t {threads} --log-file {outfile}_xss.txt")
