module Models where

----------------------------------------------------

keys :: (Ord k) => [(k, a)] -> [k]
keys = map fst

elems :: (Ord k) => [(k,a)] -> [a]
elems = map snd

lookupMap :: (Ord k) => [(k,a)] -> k -> Maybe a
lookupMap [] k = Nothing
lookupMap ((x,y):xs) k | k == x = Just y
                    | k < x = Nothing
                    | otherwise = lookupMap xs k

insert :: (Ord k, Ord a) => (k,a) -> [(k,a)] -> [(k,a)]
insert x [] = [x]
insert (x, y) ((d,f):xs) | x > d  = (d,f) : (insert (x,y) xs)
                         | x == d = ((d,(min y f)):xs)
                         | otherwise = (x,y):(d,f):xs

insertList :: (Ord k, Ord a) => [(k,a)] -> [(k,a)] -> [(k,a)]
insertList [] xs = xs
insertList ((x,y):ys) xs = ys `insertList` ((x,y) `insert` xs)

unqPairList :: (Ord k, Ord a) => [(k, a)] -> [(k, a)]
unqPairList arr = arr `insertList` [] 

----------------------------------------------

data Automat = Automat{
    prefixesAndColumns :: [([Bool], String)],
    suffixes :: [String],
    mazeSize :: (Int, Int)
}

emptyAutomat :: (Int, Int) -> Automat
emptyAutomat (x,y) = Automat [([False], "")] [""] (x,y)


--------------------------------------------------------------------------------

-- Конкатенация списков с обеспечением уникальности элементов(второй список должен быть гарантировано отсортирован)
unqConcat :: [String] -> [String] -> [String]
unqConcat [] xs = xs
unqConcat (y:ys) xs = y `unqAppend` (ys `unqConcat` xs)  

unqCompare :: String -> String -> Bool
unqCompare x y | (length x) > (length y) = True
               | (length x) == (length y) = x > y
               | otherwise = False

unqAppend :: String -> [String] -> [String]
unqAppend x [] = [x]
unqAppend x (y:ys) | x `unqCompare` y = y : (unqAppend x ys)
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

{-Генерирует список всех суффиксов, начиная с пустого. Результат отсортирован -}
generateSuffixes :: String -> [String]
generateSuffixes str = reverse (genSuf str) where
    genSuf [] = [""]
    genSuf (x:xs) = (x:xs) `unqAppend` genSuf xs

-- Расширяет список строк
expandList :: [String] -> [String]
expandList [] = []
expandList (x:xs) = [ (x++"N"), (x++"S"), (x++"W"), (x++"E"), x] `unqConcat` (expandList xs)

