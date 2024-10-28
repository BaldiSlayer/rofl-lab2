module Solver where

import Models
import Checker
--import Data.Map (Map, findWithDefault, fromList, member)


{-Собирает строки таблицы классов эквивалентности в один список
generateMaplistFromList :: [String] -> [String] -> IO [([Bool], String)]
generateMaplistFromList [] sl = do 
    return []
generateMaplistFromList (x:xs) sl = do 
    s <- listisInLanguage sl x 
    stail <- generateMaplistFromList xs sl
    return ((s, x):stail)
-}

-- Мемоизированная версия функции
generateMaplistFromListCheck :: Automat -> [String] -> [String] -> IO (Automat, [([Bool], String)])
generateMaplistFromListCheck a [] sl = do
    return (a,[])
generateMaplistFromListCheck a (x:xs) sl = do
    (a1, s) <- listisInLanguageCheck a sl x
    (a2, stail) <- generateMaplistFromListCheck a1 xs sl
    return (a2, (s,x) `insert` stail)

{-Создает таблицу классов эквивалентности по строке-}
generateAutomat :: String -> IO Automat
generateAutomat str = do
                (a, mapa) <- generateMaplistFromListCheck emptyAutomat prefList suffList
                return $ newAutomat (unqPairList mapa) suffList (knownResults a)
                    where 
                        prefList = expandList $ generatePrefixes str
                        suffList = ["", "N", "S", "W", "E"]

{-Функция добавляет суффиксы строки к существующей таблице классов эквивалентности-}
addSufToAutomat :: Automat -> String -> IO Automat
addSufToAutomat automat str = let 
    suff = (expandList $ generateSuffixes str) `unqConcat` (suffixes automat)
    in do
        (a1, table) <- generateMaplistFromListCheck automat (elems (prefixesAndColumns automat)) suff
        return $ newAutomat table suff (knownResults a1)

-- Функция добавляет префиксы строки с сущетсвующей таблице классов эквивалентности
addPrefToAutomat :: Automat -> String -> IO Automat
addPrefToAutomat automat str = let
    suff = suffixes automat  
    in do
        (a1, expandedTable) <- generateMaplistFromListCheck automat (expandList $ generatePrefixes str) suff
        table <- return (insertList (prefixesAndColumns automat) expandedTable)
        return $ newAutomat table suff (knownResults a1)

-- Функция добавляет и префиксы и суффиксы к существующей таблице классов эквивалентности
addStringToAutomat :: Automat -> String -> IO Automat
addStringToAutomat automat str = let
    suffs = (expandList $ generateSuffixes str) `unqConcat` (suffixes automat)
    prefs = (expandList $ generatePrefixes str) `unqConcat` (elems (prefixesAndColumns automat))
    in do
        (a1, mapa)  <- generateMaplistFromListCheck automat prefs suffs
        return $ newAutomat mapa suffs (knownResults a1)

