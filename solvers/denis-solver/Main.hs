module Main where

import Printer
import Solver
import Models
{-
Главный исполняемый файл солвера. Он компилится с помощью команды `ghc solver.hs`. После он должен вызываться с помощью файла executer.py
Этот файл организует интерфейс, позволяющий итеративно проводить вычисление классов эквивалентности
-}

checkAutomat :: Automat -> IO(Bool, String)
checkAutomat a = do
    putStrLn $ generateStringOfTable a
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
    auto <- generateAutomat ""
    loop auto

check :: IO()
check = do
    auto <- generateAutomat ""
    putStrLn $ generateStringOfTable auto

