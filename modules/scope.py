import json, re
def getDomains(program):
    with open(f"programs/{program}/domains.txt", "r") as f:
        return f.read().split("\n")

def getScope(program):
    with open(f"programs/{program}/config.json", "r") as f:
        data=json.loads(f.read())
        scope = data["target"]["scope"]
        include=scope["include"]
        exclude=scope["exclude"]
        return include, exclude

def checkScope(include, exclude, result):
    exists=False
    for i in include:
        if i["enabled"]:
            inScope = re.match(i["host"], result)
            if inScope:
                exists=True
                break
    if not exists:
        return False
    for i in exclude:
        if i["enabled"]:
            exScope = re.match(i["host"], result)
            if exScope:
                return False
    return True
