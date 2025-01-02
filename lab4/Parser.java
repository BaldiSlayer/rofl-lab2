/*
grammar - LL(1)
<S>			::= <A> <Atail>
<Atail>		::= '|' <A> <Atail> | \epsilon
<A>			::= <rg> <star> <rgtail>
<rgtail>	::= <rg> <star> <rgtail> | \epsilon
<rg>		::= '(' <S> ')'
				| '(:' <S> ')'
				| '(?[num])'
				| '\[num]'
				| [a-z]
<star>		::= '*' | \epsilon
*/

class Parser{
	String text;
	LinkedList<Lexem> lexemList;
	
	public static void main(String[] args){
		String test="((?: a | bc* | c)a(?1))\1*  |\td";
		
		Parser r = new Parser();
		try{
			System.out.println(r.parse(test));
		}catch(LexerException e){
			System.err.println(e);
		}
	}
	
	Regex parse(String msg) throws GrammarException{
		this.text = msg;
		
		Lexer l = new Lexer();
		lexemList = l.lex(msg);
		
		parseS();
		
	}
	
	void parseS(){}
	void parseAltTail(){}
	void parseAlt(){}
	void parseRegTail(){}
	void parseReg(){}
	void parseStar(){}
	
}
class ParserException extends GrammarException{
	ParserException(String msg){
		super("parse: "+msg);
	}
}

class Regex{
	
}