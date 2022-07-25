import os,uuid
def ContentDiscovery(target, threads, recursive):
    filename = f"programs/{program}/{uuid.uuid4()}"
    os.system(f"./theeCD/theeCD -w theeCD/wordlists/ -url {target} -t {threads} -r {recursive}|tee {filename}")
    with open(filename, "r") as f:
        content=f.read()
    output=[]
    lines=content.split("\n")
    for i in lines:
        result=":".join(i.split(":")[0:2])
        if result != "":
            output.append(result)
    return output
