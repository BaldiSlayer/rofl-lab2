module Models where

import Data.Map (Map, elems, fromList, toList)

mapToList :: Map k a -> [(k,a)]
mapToList mapa = Data.Map.toList mapa

mapFromList ::(Ord k) => [(k,a)] -> Map k a
mapFromList xs = Data.Map.fromList xs

newAutomat :: Map[Bool] String -> [String] -> Automat
newAutomat m l = Automat m l

data Automat = Automat{
    prefixesAndColumns :: Map [Bool] String,
    suffixes :: [String]
}
