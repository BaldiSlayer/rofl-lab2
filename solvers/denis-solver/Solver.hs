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

filterOut :: Automat -> [([Bool], String)] -> IO [([Bool], String)]
filterOut a [] = return []
filterOut a (x:xs) = do
    r1 <- isOut a x
    r2 <- filterOut a xs
    if r1
        then return r2
        else return (x:r2)
    where
    isOut auto (bs, str) = do
        if str == ""
            then return False
            else do
                (_, res) <- listisInLanguageCheck auto ["", "S", "W", "N", "E"] str
                return (res  == [True, True, True, True, True]) 

-- Мемоизированная версия функции
generateMaplistFromListCheck :: Automat -> [String] -> [String] -> IO (Automat, [([Bool], String)])
generateMaplistFromListCheck a [] sl = do
    return (a,[])
generateMaplistFromListCheck a (x:xs) sl = do
    (a1, s) <- listisInLanguageCheck a sl x
    (a2, stail) <- generateMaplistFromListCheck a1 xs sl
    return (a2, (s,x) `insert` stail)

{-Создает таблицу классов эквивалентности по строке-}
generateAutomat ::(Int, Int) -> String -> IO Automat
generateAutomat size str = do
                (a, mapa) <- generateMaplistFromListCheck (emptyAutomat size) prefList suffList
                mapa1 <- filterOut a mapa
                return $ (Automat (unqPairList mapa1) suffList (knownResults a) size)
                    where 
                        prefList = expandList $ generatePrefixes str
                        suffList = ["","E","N","S","W"]

{-Функция добавляет суффиксы строки к существующей таблице классов эквивалентности-}
addSufToAutomat :: Automat -> String -> IO Automat
addSufToAutomat automat str = let 
    suff = (expandList $ generateSuffixes str) `unqConcat` (suffixes automat)
    in do
        (a1, table) <- generateMaplistFromListCheck automat (elems (prefixesAndColumns automat)) suff
        return $ (Automat table suff (knownResults a1) (mazeSize a1))

-- Функция добавляет префиксы строки с сущетсвующей таблице классов эквивалентности
addPrefToAutomat :: Automat -> String -> IO Automat
addPrefToAutomat automat str = let
    suff = suffixes automat  
    in do
        (a1, expandedTable) <- generateMaplistFromListCheck automat (expandList $ generatePrefixes str) suff
        table <- return (insertList (prefixesAndColumns automat) expandedTable)
        return $ (Automat table suff (knownResults a1) (mazeSize a1))

-- Функция добавляет и префиксы и суффиксы к существующей таблице классов эквивалентности
addStringToAutomat :: Automat -> String -> IO Automat
addStringToAutomat automat str = let
    suffs = (expandList $ generateSuffixes str) `unqConcat` (suffixes automat)
    prefs = (expandList $ generatePrefixes str) `unqConcat` (elems (prefixesAndColumns automat))
    in do
        (a1, mapa)  <- generateMaplistFromListCheck automat prefs suffs
        return $ (Automat mapa suffs (knownResults a1) (mazeSize a1))

