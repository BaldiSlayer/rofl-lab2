module Solver where

import Models
import Data.Map (fromList)

{-Генерирует список всех префиксов, начиная с пустого префикса-}
generatePrefixes :: [Char] -> [[Char]]
generatePrefixes [] = []
generatePrefixes arr = let
                            pref p [] = [p]
                            pref p (x:xs) = p : (pref (p++[x]) xs)
                       in pref [] arr 

{-Генерирует список всех суффиксов. Расположены в обратном порядке-}
generateSuffixes :: String -> [String]
generateSuffixes [] = [""]
generateSuffixes (x:xs) = (x:xs) : generateSuffixes xs


expandStringList :: [String] -> [String]
expandStringList [] = []
expandStringList (x:xs) = x : (x++"N") : (x++"S") : (x++"W") : (x++"E") : (expandStringList xs)
{- (x++"N") : ((x++"S") : ((x++"W") : ((x++"E") : (expandStringList xs))))-}


{-Делает запрос к MAT(или пользователь) на получение данных о вхождение в язык-}
isInLanguage :: String -> IO Bool
isInLanguage str = do
    putStrLn str
    res <- getLine
    return ((head res) == '1')


{-Создает строку таблицы классов эквивалентности. Принимает префикс и список суффиксов-}
listisInLanguage :: [String] -> String -> IO [Bool]
listisInLanguage [] str  = do
                            return []
listisInLanguage (x:xs) str = do 
                            res1 <- isInLanguage $ str++x
                            res2 <- listisInLanguage xs str
                            return (res1 : res2)

{-Собирает строки таблицы классов эквивалентности в один список-}
generateMaplistFromList :: [String] -> [String] -> IO [([Bool], String)]
generateMaplistFromList [] sl = do 
    return []
generateMaplistFromList (x:xs) sl = do 
    s <- listisInLanguage sl x 
    stail <- generateMaplistFromList xs sl
    return ((s, x):stail)

{-Создает таблицу классов эквивалентности по строке-}
generateAutomat :: String -> IO Automat
generateAutomat str = do
                mapa <- generateMaplistFromList prefList suffList
                return $ newAutomat (mapFromList mapa) suffList
                    where 
                        prefList = expandStringList $ generatePrefixes str
                        suffList = ["", "N", "S", "W", "E"]






{-Функция добавляет строку к существующей таблице классов эквивалентности-}
{-addStrToAutomat :: Automat -> String -> Automat
addStrToAutomt automat str = 
-}

