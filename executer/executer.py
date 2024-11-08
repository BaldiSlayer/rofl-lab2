#!/usr/bin/python3

import sys, os
import subprocess as sp

def readLen(pipe):
    return os.fstat(pipe.fileno()).st_size

def runSolve(MATpath, solverpath):
    solver = sp.Popen(solverpath, stdout=sp.PIPE, stdin=sp.PIPE, stderr=sys.stdout, text=True, bufsize=1)
    MAT = sp.Popen(MATpath, stdout=solver.stdin, stdin=solver.stdout, stderr=sys.stdout,text=True, bufsize=1)

    try:
        solver.wait()
    except:
        for line in solver.stderr:
            print("Solver-err: "+line, end="")
        for line in MAT.stderr:
            print("MAT-err: "+line, end="")
    finally:
        #print("MAT-err:" + MAT.stderr.read())
        #print("solver-err:"+solver.stderr.read())
        MAT.kill()
        solver.kill()
        

if __name__ == "__main__":
    MAT = "./MAT"
    solver = "./Solver"    
    if len(sys.argv) == 3: #executer.pu + MAT path + solver path
        MAT, solver = list(map(lambda x: x.split(" "), sys.argv[1:]))
    
    print("used mat: {0}\nsolver: {1}".format(MAT, solver))
    #optional: check if it MAT and if it solver    
    runSolve(MAT, solver)
