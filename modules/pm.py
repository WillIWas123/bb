import os,uuid
def ParamMiner(target, threads, program):
    filename = f"programs/{program}/{uuid.uuid4()}"
    os.system(f"./theePM/theePM -w theePM/wordlists/ -url {target} -t {threads}|tee {filename}")
    with open(filename, "r") as f:
        content=f.read()
    output=[]
    lines=content.split("\n")
    for i in lines:
        result=":".join(i.split(":")[0:2])
        if result != "":
            output.append(result)
    return output
