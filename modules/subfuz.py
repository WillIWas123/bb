import os, uuid
def subfuz(domains, program):
    outfile = f"programs/{program}/{uuid.uuid4()}"
    output=[]
    for i in domains:
        if i == "":
            continue
        if i[0] == "*":
            domain = i[2:]
            output.append(domain)
            os.system(f"python3 /opt/subfuz/subfuz.py -d {domain} -all -csv {outfile}")
        else:
            output.append(i)
        with open(outfile, "r") as f:
            data=f.read().split("\n")
        for j in data:
            subdomain=j.split(",")[0]
            output.append(subdomain)
    return output
