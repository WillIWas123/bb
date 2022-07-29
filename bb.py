import json, sys, os,re,uuid
from modules.pm import ParamMiner
from modules.cd import ContentDiscovery
from modules.xss import XSStrike
from modules.tplmap import TPLMap
from modules.scope import getScope, checkScope,getDomains
from modules.subfuz import subfuz

def main(program=None,threads=None,recursiveness=None):
    if len(sys.argv) >=4:
        program = sys.argv[1]
        threads = sys.argv[2]
        recursiveness = sys.argv[3]
    include, exclude = getScope(program)
    domains=getDomains(program)
    subdomains = subfuz(domains, program)
    for i in subdomains:
        inScope = checkScope(include, exclude,i)
        if not inScope:
            continue
        cd = ContentDiscovery(f"https://{i}/", threads, recursiveness,program)
        for j in cd:
            pm = ParamMiner(j, threads,program)
            for k in pm:
                XSStrike(k, threads,program)
                TPLMap(k,threads,program)

if __name__ == "__main__":
    if len(sys.argv) < 4:
        print(f"Usage: {sys.argv[0]} [program] [threads] [recursiveness]")
        sys.exit(1)
    main()
