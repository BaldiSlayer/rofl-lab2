class FrameGrammar{
	FrameGrammar(Regex r) throws GrammarException{
		Checker c = new Checker();
		c.check(r);
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
		FrameGrammar fg = new FrameGrammar();
		
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

class Grammar{}