import java.util.LinkedList;

class Checker{
	public LinkedList<Group> groups;
	
	
	void check(Regex r) throws GrammarException{
		// проверить есть ли захваты слов на группы под звездочкой и в альтернативе
		groups = new LinkedList<Group>();
		
		//recursiveFound(r, true); // список всех групп захвата, нужно для ускорения поиска
		recursiveFindGroup(r);
		recursiveCheck(r);
	}
	void recursiveFindGroup(Group r){
		
		for(LinkedList<Reg> alt : r.alternatives){
			for(Reg i : alt){
				if(i instanceof Group){//если группа
					Group groupa = (Group)i;
					if(groupa.value > 0){
						groups.add(groupa);
					}
					recursiveFindGroup(groupa);
				}
			}
		}
	}
	
	Group recursiveFound(Group r, Reg wordgrab){
		if(r.alternatives.size() > 1)
			return null;
		
		for(LinkedList<Reg> alt : r.alternatives){
			for(Reg i : alt){
				if(i instanceof Group){//если группа
					Group groupa = (Group)i;
					if(groupa.value == -wordgrab.value && !groupa.star && groupa.startIndex<wordgrab.startIndex){
						return groupa;
					}
					if(!groupa.star){
						groupa = recursiveFound(groupa, wordgrab);
						if(groupa != null)
							return groupa;
					}
				}
			}
		}
		return null;
	}
	void recursiveCheck(Group r) throws GrammarException{
		for(LinkedList<Reg> alt : r.alternatives){
			for(Reg i : alt){
				if(!(i instanceof Group)){ //если reg,
					if(i.value < 0){
						/*for(Group g: groups){
							if(g.value == -i.value){
								if(!g.canBeWordGrab)
									throw new GrammarException(String.format("На позиции %d ссылка на группу в альтернативе или под замыканием Клини %d: %s",
										i.startIndex, g.value, g.toString()));
								if(i.startIndex < g.startIndex)
									throw new GrammarException(String.format("На позиции %d ссылка на неинициализированную группу %d: %s",
										i.startIndex, g.value, g.toString()));
								break;
							}
						}*/
						Group g = null;
						for(Reg j : alt){
							if(j instanceof Group){
								if(j.value == -i.value && !j.star && j.startIndex < i.startIndex){
									g = (Group)j;
									break;
								}
								recursiveFound((Group)j, i);
							}
						}
						if(g == null){
							for(Group j: groups){
								if(j.value == -i.value){
									g = j;
									break;
								}
							}
							if(g == null){
								throw new GrammarException(String.format("На позиции %d ссылка на неинициализированную группу %d",
									i.startIndex, -i.value));
							}else{
								throw new GrammarException(String.format("На позиции %d ссылка на неинициализированную группу %d: %s",
									i.startIndex, -i.value, g.toString()));
							}
						}
						
					}
				}else{
					recursiveCheck((Group)i);
				}
			}
		}
	}
	
	public static void main(String[] args){
		LinkedList<String> tests = new LinkedList<String>();
		
		tests.add("\\1(a|b)");				// no
		tests.add("(?1)(a|b)");				// ok
		tests.add("(a(?1)b | c)");			// ok
		tests.add("(a | (bb)) (?2)");		// ok
		tests.add("(a(bb) | b (cc)) \\2"); 	// no
		tests.add("( (?: a (?2) | (bb) ) (?1))"); // ok
		tests.add("(a (?1))");				// ok
		tests.add("(a)(b(c))(d(e(f)(g)))");	// ok
		tests.add("(a |(?:(b))) \\2"); 		// no
		tests.add("(a) \\1 | b");			// ok
		tests.add("(a)\\3");				// no
		tests.add("(aa");					// no
		
		Parser p = new Parser();
		Checker c = new Checker();
		
		int counter=1;
		for(String test : tests){
			System.out.println("unit test "+ (counter++) + ":");
			
			try{
				Regex r = p.parse(test);
				c.check(r);
				System.out.println(r);
			}catch(GrammarException e){System.out.println(e);}
			
			System.out.println("---------------------");
		}
	}
}