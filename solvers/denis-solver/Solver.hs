module Solver where

import Models
import Checker
import Printer
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

revertPath :: String -> String
revertPath "" = ""
revertPath str = map rev $ reverse str where
    rev 'N' = 'S'
    rev 'S' = 'N'
    rev 'E' = 'W'
    rev 'W' = 'E'
    rev x = x

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
addSufsToAutomat :: Automat -> [String] -> IO Automat
addSufsToAutomat automat strlist = let 
    suff = strlist `unqConcat` (suffixes automat)
    in do
        (a1, mapa) <- generateMaplistFromListCheck automat (elems (prefixesAndColumns automat)) suff
        --mapa1 <- filterOut a1 mapa
        return $ (Automat mapa suff (knownResults a1) (mazeSize a1))

-- Функция добавляет префиксы строки с сущетсвующей таблице классов эквивалентности
addPrefsToAutomat :: Automat -> [String] -> IO Automat
addPrefsToAutomat automat strlist = let
    suff = suffixes automat  
    in do
        (a1, expandedTable) <- generateMaplistFromListCheck automat strlist suff
        mapa <- return (insertList (prefixesAndColumns automat) expandedTable)
        --mapa1 <- filterOut a1 mapa
        return $ (Automat mapa suff (knownResults a1) (mazeSize a1))

-- Функция добавляет и префиксы и суффиксы к существующей таблице классов эквивалентности
addStringToAutomat :: Automat -> String -> IO Automat
addStringToAutomat automat str = let
    suffs = (generatePrefixes $ revertPath str) `unqConcat` ((expandList $ generateSuffixes str) `unqConcat` (suffixes automat))
    prefs = (expandList $ generatePrefixes str) `unqConcat` (elems (prefixesAndColumns automat))
    in do
        (a1, mapa)  <- generateMaplistFromListCheck automat prefs suffs
        --mapa1 <- filterOut a1 mapa
        return $ (Automat mapa suffs (knownResults a1) (mazeSize a1))

-----------------------------------------------------------------------------
goAlong :: String -> (Int, Int)
goAlong str = goAlong' str (0,0) where
    goAlong' [] (a,b) = (a,b)
    goAlong' ('N':xs) (a,b) = goAlong' xs (a,b-1)
    goAlong' ('S':xs) (a,b) = goAlong' xs (a, b+1)
    goAlong' ('W':xs) (a,b) = goAlong' xs (a-1, b)
    goAlong' ('E':xs) (a,b) = goAlong' xs (a+1, b)
    goAlong' (x:xs) (a,b) = (0,0)

generatePathes :: (Int, Int) -> (Int, Int) -> [String]
generatePathes (x,y) (w, h) | x == -1 = generatePathes' (x,y) (w,h) (x,y+1) "S"
                            | y == -1 = generatePathes' (x,y) (w,h) (x-1,y) "W"
                            | x == w = generatePathes' (x,y) (w,h) (x,y-1) "N"
                            | y == h = generatePathes' (x,y) (w,h) (x+1,y) "E"
                            | otherwise = []
    where
        generatePathes' (x0, y0) (w,h) (x,y) p | x == x0 && y == y0 = [""]
                                               | x == -1 && y == h = p : (generatePathes' (x0,y0) (w,h) (x+1, y) (p++"E"))
                                               | y == -1 && x == -1 = p : (generatePathes' (x0,y0) (w,h) (x, y+1) (p++"S"))
                                               | y == h && x == w = p : (generatePathes' (x0,y0) (w,h) (x, y-1) (p++"N"))                             
                                               | x == w && y == -1 = p : (generatePathes' (x0,y0) (w,h) (x-1, y) (p++"W"))
                                               | x == w = p : (generatePathes' (x0,y0) (w,h) (x, y-1) (p++"N"))
                                               | y == -1 = p : (generatePathes' (x0,y0) (w,h) (x-1, y) (p++"W"))
                                               | x == -1 = p : (generatePathes' (x0,y0) (w,h) (x, y+1) (p++"S"))
                                               | y == h = p : (generatePathes' (x0,y0) (w,h) (x+1, y) (p++"E"))
                                               | otherwise = [""]

-- добавляет префиксы до граничных клеток и суффиксы
addBorder :: Automat -> String -> IO Automat
addBorder a path = do
    (_, line) <- listisInLanguageCheck a (suffixes a) path -- получаем строку таблицы классов эквивалентности для выхода из лабиринта
    minPathMB <- return $ lookupMap (prefixesAndColumns a) line -- получаем минимизированный вариант(гарантировано будет)
    minPath <- return $ maybe path (\x-> x) minPathMB
    (x,y) <- return $ goAlong minPath -- получает координаты выхода из лабиринта
    --putStrLn $ show (x,y)
    pathList <- return $ map ((++) minPath) $ generatePathes (x,y) (mazeSize a)
    pathRev <- return $ map revertPath pathList
    --putStrLn $ show pathList
    --putStrLn $ show pathRev
    a1 <- addSufsToAutomat a pathRev
    a2 <- addPrefsToAutomat a1 pathList
    --putStrLn $ generateStringOfTable a
    --putStrLn $ generateStringOfTable a1
    --putStrLn $ generateStringOfTable a2
    return a2
    where
        
