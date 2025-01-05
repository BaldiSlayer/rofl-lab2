import java.util.LinkedList;

class Checker{
	LinkedList<Group> groups;
	
	
	void check(Regex r) throws GrammarException{
		// проверить есть ли захваты слов на группы под звездочкой и в альтернативе
		groups = new LinkedList<Group>();
		
		recursiveFound(r); // список всех групп захвата, нужно для ускорения поиска
		recursiveCheck(r);
	}
	
	void recursiveFound(Group r){
		for(LinkedList<Reg> alt : r.alternatives){
			for(Reg i : alt){
				if(i instanceof Group){//если группа
					Group groupa = (Group)i;
					if(groupa.value > 0){
						if(!groupa.star && (r.alternatives.size() == 1)){
							groupa.canBeWordGrab = true;
						}
						
						groups.add(groupa);
					}
					recursiveFound(groupa);
				}
			}
		}
	}
	void recursiveCheck(Group r) throws GrammarException{
		for(LinkedList<Reg> alt : r.alternatives){
			for(Reg i : alt){
				if(!(i instanceof Group)){ //если reg,
					if(i.value < 0){
						for(Group g: groups){
							if(g.value == -i.value){
								if(!g.canBeWordGrab)
									throw new GrammarException(String.format("На позиции %d ссылка на группу в альтернативе или под замыканием Клини %d: %s",
										i.startIndex, g.value, g.toString()));
								if(i.startIndex < g.startIndex)
									throw new GrammarException(String.format("На позиции %d ссылка на неинициализированную группу %d: %s",
										i.startIndex, g.value, g.toString()));
								break;
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
		tests.add("(?1)(a|b)");				// yes
		tests.add("(a(?1)b | c)");			// yes
		tests.add("(a | (bb)) (?2)");		// yes
		tests.add("(a(bb) | b (cc)) \\2"); 	// no
		tests.add("( (?: a (?2) | (bb) ) (?1))"); // yes
		tests.add("(a (?1))");
		tests.add("(a)(b(c))(d(e(f)(g)))");
		
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