module Checker where
import Models

{-Делает запрос к MAT(или пользователь) на получение данных о вхождение в язык-}
isInLanguage :: String -> IO Bool
isInLanguage str = do
            putStrLn str
            res <- getLine
            return ((length res) /= 0 && (head res) == '1') -- для простоты тестирования, "не 1" - 0

{-Создает строку таблицы классов эквивалентности. Принимает список суффиксов и префикс
listisInLanguage :: [String] -> String -> IO [Bool]
listisInLanguage [] str  = do
                            return []
listisInLanguage (x:xs) str = do 
                            res1 <- isInLanguage $ str++x
                            res2 <- listisInLanguage xs str
                            return (res1 : res2)
-}

isInLanguageCheck :: Automat -> String -> IO (Automat, Bool)
isInLanguageCheck auto str = let 
    check = lookupMap (knownResults auto)
    isIn = check str
    isJust Nothing = False
    isJust (Just x) = True
    in do 
        if (isJust isIn)
            then return $ maybe (auto, False) (\x->(auto, x)) isIn
            else do
                res <- isInLanguage str
                newcheck <- return $ insert (str, res) (knownResults auto)
                return ((newAutomat (prefixesAndColumns auto) (suffixes auto) newcheck ), res) 

listisInLanguageCheck :: Automat -> [String] ->  String -> IO (Automat, [Bool])
listisInLanguageCheck auto [] _ = do return (auto, [])
listisInLanguageCheck auto (x:xs) str = do
    (newauto1, res1) <- isInLanguageCheck auto $ str++x
    (newauto2, res2) <- listisInLanguageCheck newauto1 xs str
    return (newauto2, (res1:res2))
