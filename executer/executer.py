#!/usr/bin/python3

import sys
import subprocess as sp

def runSolve(MATpath, solverpath):
    MAT = sp.Popen(MATpath, stdout=sp.PIPE, stdin=sp.PIPE)
    solver = sp.Popen(solverpath, stdout=sp.PIPE, stdin=sp.PIPE)
    solvermessage, _ = solver.communicate("start".encode())    
    while solver.poll() == None:
        print(solvermessage.decode())
        MATmessage, _ = MAT.communicate(solvermessage)
        print(MATmessage.decode())
        solvermessage, _ = solver.communicate(MATmessage)
        

if __name__ == "__main__":
    if len(sys.argv) != 3: #executer.pu + MAT path + solver path
        raise Exception("need 2 arguments: path to MAT, path to Solver")
    MAT, solver = sys.argv[1:]
    #optional: check if it MAT and if it solver    
    runSolve(MAT, solver)
