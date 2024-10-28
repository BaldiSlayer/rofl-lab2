module FileItteract where

import System.IO

readFromFile :: String -> IO (Int, Int)
readFromFile filename = do
    fileHandle <- openFile filename ReadMode
    contents <- hGetLine fileHandle
    arg1 <- return $ read contents
    contents <- hGetLine fileHandle
    arg2 <- return $ read contents
    return (arg1, arg2)
