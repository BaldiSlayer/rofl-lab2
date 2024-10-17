module Main where

import Printer
import Solver
{-
Главный исполняемый файл солвера. Он компилится с помощью команды `ghc solver.hs`. После он должен вызываться с помощью файла executer.py
Этот файл организует интерфейс, позволяющий итеративно проводить вычисление классов эквивалентности
-}

{-loop :: Automat -> IO ()
loop automat = do
    ans <- getLine
    if ans == "TRUE"
        then return
        else else_branch
    where
        else_branch = do
            newautomat <- addStrToAutomat automat ans
            putStrLn $ generateStringOfTable newautomat
            loop newautomat
-}


main :: IO()
main = do
    str <- getLine
    automat <- generateAutomat str
    putStrLn $ generateStringOfTable automat
