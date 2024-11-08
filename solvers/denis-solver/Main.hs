module Main where

import Printer
import Solver
import Models
import FileItteract

import System.IO (hFlush, stdout)
{-
Главный исполняемый файл солвера. Он компилится с помощью команды `ghc solver.hs`. После он должен вызываться с помощью файла executer.py
Этот файл организует интерфейс, позволяющий итеративно проводить вычисление классов эквивалентности
-}

sendTable :: String -> IO()
sendTable str = do
    putStrLn "table"
    putStr str
    putStrLn "end" 
    hFlush stdout

checkAutomat :: Automat -> IO(Bool, String)
checkAutomat a = do
    sendTable $ generateStringOfTable a
    str <- getLine
    if str == "TRUE"
        then return (True, "")
        else return (False, str)


loop :: Automat -> IO()
loop a = do
    (res, contr) <- checkAutomat a
    if res 
        then putStrLn "Succes"
        else do
            newauto <- addStringToAutomat a contr
            loop newauto    

main :: IO()
main = do
    (x,y) <- readFromFile "./parameters.txt"
    --putStrLn $ show (x,y)
    sendTable "e\ne 0\n"
    str <- getLine
    auto <- addStringToAutomat (emptyAutomat (x,y)) str
    auto1 <- addBorder auto str 
    loop auto1

