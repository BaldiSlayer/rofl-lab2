#!/usr/bin/python3

import sys
import subprocess as sp



def runSolve(MATpath, solverpath):
    MAT = sp.Popen(MATpath, stdout=sp.PIPE, stdin=sp.PIPE, stderr=sp.PIPE,text=True, bufsize=1)
    solver = sp.Popen(solverpath, stdout=sp.PIPE, stdin=sp.PIPE, stderr=sp.PIPE, text=True, bufsize=1)
    try:
        solver.stdin.write("isin\nS\n")
        solver.stdin.flush()
        s = solver.stdout.readline().strip()
        print(s)
    except:
        for line in solver.stderr:
            print("Solver-err: "+line, end="")
        for line in MAT.stderr:
            print("MAT-err: "+line, end="")
    finally:
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
