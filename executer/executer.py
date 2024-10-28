#!/usr/bin/python3

import sys
import subprocess as sp

def runSolve(MATpath, solverpath):
    MAT = sp.Popen(MATpath, stdout=sp.PIPE, stdin=sp.PIPE)
    solver = sp.Popen(solverpath, stdout=sp.PIPE, stdin=sp.PIPE)
    solvermessage, _ = solver.communicate()    
    while solver.poll() == None:
        print(solvermessage.decode())
        MATmessage, _ = MAT.communicate(solvermessage)
        print(MATmessage.decode())
        solvermessage, _ = solver.communicate(MATmessage)
        

if __name__ == "__main__":
    MAT = "./MAT"
    solver = "./Solver"    
    if len(sys.argv) == 3: #executer.pu + MAT path + solver path
        MAT, solver = sys.argv[1:]
    
    print("used mat: {0}\nsolver: {1}\n".format(MAT, solver))
    #optional: check if it MAT and if it solver    
    runSolve(MAT, solver)
