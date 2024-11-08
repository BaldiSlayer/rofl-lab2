module Checker where
import Models
import System.IO (hFlush, stdout)

{-Делает запрос к MAT(или пользователь) на получение данных о вхождение в язык-}
isInLanguage :: String -> IO Bool
isInLanguage str = do
            putStrLn "isin"
            putStrLn str
            hFlush stdout
            res <- getLine
            return ((length res) /= 0 && (head res) == '1') -- для простоты тестирования, "не 1" - 0

--Создает строку таблицы классов эквивалентности. Принимает список суффиксов и префикс
listisInLanguage :: [String] -> String -> IO [Bool]
listisInLanguage [] str  = do
                            return []
listisInLanguage (x:xs) str = do 
                            res1 <- isInLanguage $ str++x
                            res2 <- listisInLanguage xs str
                            return (res1 : res2)

