module Models where

-- import Data.Map as Map (Map, elems, fromList, toList)
-- Почему-то перестала работать библиотека
----------------------------------------------------

keys :: (Ord k) => [(k, a)] -> [k]
keys = map fst

elems :: (Ord k) => [(k,a)] -> [a]
elems = map snd

insert :: (Ord k) => (k,a) -> [(k,a)] -> [(k,a)]
insert x [] = [x]
insert (k, e) ((d,f):xs) | k > d  = (d,f) : (insert (k,e) xs)
                                | k == d = ((d,f):xs)
                                | otherwise = (k,e): (d,f):xs

insertList :: (Ord k) => [(k,a)] -> [(k,a)] -> [(k,a)]
insertList [] xs = xs
insertList ((x,y):ys) xs = ys `insertList` ((x,y) `insert` xs)

unqPairList :: (Ord k) => [(k, a)] -> [(k, a)]
unqPairList arr = arr `insertList` [] 

----------------------------------------------

data Automat = Automat{
    prefixesAndColumns :: [([Bool], String)],
    suffixes :: [String]
}

newAutomat :: [([Bool], String)] -> [String] -> Automat
newAutomat m l = Automat m l

--------------------------------------------------------------------------------

-- Конкатенация списков с обеспечением уникальности элементов(второй список должен быть гарантировано отсортирован)
unqConcat :: (Ord a) => [a] -> [a] -> [a]
unqConcat [] xs = xs
unqConcat (y:ys) xs = ys `unqConcat` (y `unqAppend` xs)  

unqAppend :: (Ord a) =>  a -> [a] -> [a]
unqAppend x [] = [x]
unqAppend x (y:ys) | x > y = y : (unqAppend x ys)
                   | x == y = y : ys
                   | otherwise = x : y : ys  

---------------------------------------------------------------------------------
{-Генерирует список всех префиксов, начиная с пустого префикса-}
generatePrefixes :: String -> [String]
generatePrefixes "" = [""]
generatePrefixes arr = let
                            pref p [] = [p]
                            pref p (x:xs) = p `unqAppend` (pref (p++[x]) xs)
                       in pref [] arr 

{-Генерирует список всех суффиксов, начиная с пустого. Результат отсортирован по длине-}
generateSuffixes :: String -> [String]
generateSuffixes str = reverse (genSuf str) where
    genSuf [] = [""]
    genSuf (x:xs) = (x:xs) `unqAppend` genSuf xs

-- Расширяет список строк
expandList :: [String] -> [String]
expandList [] = []
expandList (x:xs) = [x, (x++"N"), (x++"S"), (x++"W"), (x++"E")] `unqConcat` (expandList xs)

