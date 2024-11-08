module Solver where

import Models
import Checker
import Printer


--Собирает строки таблицы классов эквивалентности в один список
generateMaplistFromList :: [String] -> [String] -> IO [([Bool], String)]
generateMaplistFromList [] sl = do 
    return []
generateMaplistFromList (x:xs) sl = do 
    s <- listisInLanguage sl x 
    stail <- generateMaplistFromList xs sl
    return ((s, x) `insert` stail)

-- Генерирует "обратный" путь
revertPath :: String -> String
revertPath "" = ""
revertPath str = map rev $ reverse str where
    rev 'N' = 'S'
    rev 'S' = 'N'
    rev 'E' = 'W'
    rev 'W' = 'E'
    rev x = x



{-Функция добавляет суффиксы строки к существующей таблице классов эквивалентности-}
addSufsToAutomat :: Automat -> [String] -> IO Automat
addSufsToAutomat automat strlist = let 
    suff = strlist `unqConcat` (suffixes automat)
    in do
        mapa <- generateMaplistFromList (elems (prefixesAndColumns automat)) suff
        return $ (Automat mapa suff (mazeSize automat))

-- Функция добавляет префиксы строки с сущетсвующей таблице классов эквивалентности
addPrefsToAutomat :: Automat -> [String] -> IO Automat
addPrefsToAutomat automat strlist = let
    suff = suffixes automat  
    in do
        expandedTable <- generateMaplistFromList strlist suff
        mapa <- return (insertList (prefixesAndColumns automat) expandedTable)
        return $ (Automat mapa suff (mazeSize automat))

-- Функция добавляет и префиксы и суффиксы к существующей таблице классов эквивалентности
addStringToAutomat :: Automat -> String -> IO Automat
addStringToAutomat automat str = let
    suffs = (generatePrefixes $ revertPath str) `unqConcat` ((expandList $ generateSuffixes str) `unqConcat` (suffixes automat))
    prefs = (expandList $ generatePrefixes str) `unqConcat` (elems (prefixesAndColumns automat))
    in do
        mapa  <- generateMaplistFromList prefs suffs
        --mapa1 <- filterOut a1 mapa
        return $ (Automat mapa suffs (mazeSize automat))

-----------------------------------------------------------------------------

-- функция, прохождения в соответсвии с путем
goAlong :: String -> (Int, Int)
goAlong str = goAlong' str (0,0) where
    goAlong' [] (a,b) = (a,b)
    goAlong' ('N':xs) (a,b) = goAlong' xs (a,b-1)
    goAlong' ('S':xs) (a,b) = goAlong' xs (a, b+1)
    goAlong' ('W':xs) (a,b) = goAlong' xs (a-1, b)
    goAlong' ('E':xs) (a,b) = goAlong' xs (a+1, b)
    goAlong' (x:xs) (a,b) = (0,0)

-- генерирует префиксы обхода по границе лабиринта
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
    line <- listisInLanguage (suffixes a) path -- получаем строку таблицы классов эквивалентности для выхода из лабиринта
    minPathMB <- return $ lookupMap (prefixesAndColumns a) line -- получаем минимизированный вариант(гарантировано будет)
    minPath <- return $ maybe path (\x-> x) minPathMB
    (x,y) <- return $ goAlong minPath -- получает координаты выхода из лабиринта
    pathList <- return $ map ((++) minPath) $ generatePathes (x,y) (mazeSize a) -- генерирует шаги по границе лабиринта
    pathRev <- return $ map revertPath pathList -- генерирует обратные суффиксы
    a1 <- addSufsToAutomat a pathRev -- добавляет суффиксы и префиксы к автомату
    a2 <- addPrefsToAutomat a1 pathList
    return a2
    where
        
