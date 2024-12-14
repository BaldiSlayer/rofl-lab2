module Printer where

import Models 

{-
        | epsilon 
------------------
epsilon |   0     
S       |   0     
N       |   0     
W       |   0   
E       |   1    
-}

{-Получает максимальную длину строки в списке. 
Может можно было сделать с помощью встроенных функций, но я написал сам-}
maxLenString :: [String] -> Int
maxLenString [] = 0
maxLenString (x:xs) = max (length x) (maxLenString xs)


{-Генерирует первую строку вывода-}
firstLine :: Int -> [String] -> String
firstLine n names = (printNames names) ++ "\n"
    where
        printNames [] = ""
        printNames (x:xs) = " "++(checkEpsilon x)++ (printNames xs)

{-Генерирует сепаратор таблицы-}
separator :: Int -> String
separator n = replicate (n+10) '-'

checkEpsilon :: String -> String
checkEpsilon str | str == "" = "e"
                 | otherwise = str

{-Выводит саму таблицу-}
printLines :: [([Bool], String)] -> String
printLines a = printLinesList a
    where
        printLinesList [] = ""
        printLinesList ((key, str):xs) = (checkEpsilon str) ++ " " ++ ({-boolToNumStr (last key)-}toString key) ++ "\n" ++ (printLinesList xs)
        boolToNumStr x | x = "1"
                    | otherwise = "0"
        toString [] = ""
        toString (x:xs) = (boolToNumStr x) ++ " "++ (toString xs)

{-Собирает все выводы в единую строчку-}
generateStringOfTable :: Automat -> String
generateStringOfTable a = let
       table = prefixesAndColumns a
       names = suffixes a
       columnlen = maxLenString $ "epsilon":(elems table) 
    in (firstLine columnlen names) ++ (printLines table)
