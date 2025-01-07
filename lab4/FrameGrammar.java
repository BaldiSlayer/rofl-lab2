import java.util.ArrayList;
import java.util.LinkedList;

class FrameGrammar{
	int freeIndex = 0;
	int starIndex = 0;
	ArrayList<NonTerm> grammar;
	
	
	FrameGrammar(Regex r) throws GrammarException{
		Checker c = new Checker();
		c.check(r); // проверяем регекс на правильные ссылки
		
		
		grammar = new ArrayList<NonTerm>();
		
		groupToNonTerm(r, -1);
		
		/*for(Group g: c.groups){
			groups.add(groupToNonTerm(c.groups, g, groups.size()));
		}*/
		
		for(NonTerm nt : grammar){
			System.out.println(nt);
		}
	}
	NonTerm groupToNonTerm(Group g, int index){
		String name;
		if(index == -1){
			name = "S";
		}else if(index == 0){ // neutral group
			name = "P"+(freeIndex++);
		}else{ // capture group
			name = "T"+index;
		}
		
		NonTerm r = new NonTerm(name);
		grammar.add(r);
		
		for(LinkedList<Reg> alt : g.alternatives){
			r.addRule();
			
			for(Reg i : alt){
				if(i instanceof Group){
					if(i.star){ // добавить нетерминал для звездочки, вставить терминал для звездочки
						NonTerm star = new NonTerm("S"+(starIndex++));
						grammar.add(star);
						r.addToLastRule(star);
						
						star.addToLastRule(groupToNonTerm((Group)i, i.value));
						
						star.addToLastRule(star);
						
						star.addRule();
						star.addToLastRule(new Term("."));
					}else{
						// создать нетерминал рекурсивно
						r.addToLastRule(groupToNonTerm((Group)i, i.value));
					}
				}else{
					if(i.star){ // добавить нетерминал для звездочки, 
						NonTerm star = new NonTerm("S"+(starIndex++));
						grammar.add(star);
						r.addToLastRule(star);
						
						if(i.value >= 'a' && i.value <= 'z'){
							star.addToLastRule(new Term(String.valueOf((char)i.value)));
						}else{
							int groupIndex = i.value;
							if(groupIndex < 0) groupIndex = -groupIndex;
							
							star.addToLastRule(new NonTerm("T"+groupIndex));
						}
						star.addToLastRule(star);
						
						star.addRule();
						star.addToLastRule(new Term("."));
						
					}else{ // добавить терминал, если буква. Иначе добавить нетерминал соответсвующей группы
						if(i.value >= 'a' && i.value <= 'z'){
							r.addToLastRule(new Term(String.valueOf((char)i.value)));
						}else{
							int groupIndex = i.value;
							if(groupIndex < 0) groupIndex = -groupIndex;
							
							r.addToLastRule(new NonTerm("T"+groupIndex));
						}
					}
				}
			}
		}
		
		return r;
	}
	
	public static void main(String[] args){
		LinkedList<String> tests = new LinkedList<String>();
		/*
		tests.add("\\1(a|b)");				// no
		tests.add("(?1)(a*|b)");				// ok
		tests.add("(a(?1)b | c)");			// ok
		tests.add("(a | (bb*)) (?2)");		// ok
		tests.add("(a(bb) | b (cc)) \\2"); 	// no
		tests.add("( (?: a (?2)* | (bb) )* (?1))"); // ok
		tests.add("(a (?1))");				// ok
		tests.add("(a)(b(c))(d(e(f)(g)))");	// ok
		tests.add("(a |(?:(b))) \\2"); 		// no
		tests.add("(a) \\1 | b");			// ok
		tests.add("(a)\\3");				// no
		tests.add("(aa");					// no*/
		tests.add("(aaba | c)* \\1| d(?1)");
		
		Parser p = new Parser();
		
		int counter=1;
		for(String test : tests){
			System.out.println("unit test "+ (counter++) + ":");
			
			try{
				Regex r = p.parse(test);
				System.out.println(r);
				FrameGrammar fg = new FrameGrammar(r);
				
			}catch(GrammarException e){System.out.println(e);}
			
			System.out.println("---------------------");
		}
	}
}