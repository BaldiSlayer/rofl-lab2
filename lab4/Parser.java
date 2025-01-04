import java.util.LinkedList;

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


First
S - '(', '(:', '(?[num])', '\[num]', [a-z]
A - '(', '(:', '(?[num])', '\[num]', [a-z]
Atail - '|', epsilon
rg - '(', '(:', '(?[num])', '\[num]', [a-z]
rgtail - '(', '(:', '(?[num])', '\[num]', [a-z]
star - '*', epsilon
*/

class Parser{
	String text;
	LinkedList<Lexer.Lexem> lexemList;
	
	public static void main(String[] args){
		String test="((?: a | bc* | c)a(?1))\\1*  |\td";
		
		Parser r = new Parser();
		try{
			System.out.println(r.parse(test));
		}catch(GrammarException e){
			System.err.println(e);
		}
	}
	
	Regex parse(String msg) throws GrammarException{
		this.text = msg;
		
		Lexer l = new Lexer();
		lexemList = l.lex(msg);
		//System.out.println(lexemList);
		parseS();
		if(lexemList.peek() != null){
			throw new ParserException(String.format("В ходе разбора строка была разобрана не полностью. Оставшиеся лексемы %s",
				lexemList));
		}
		return null;
	}
	
	// <S> ::= <A> <Atail>
	void parseS() throws ParserException{
		//System.out.println("S");
		//System.out.println(lexemList);
		parseAlt();
		parseAltTail();
	}
	// <Atail> ::= '|' <A> <Atail> | \epsilon
	void parseAltTail() throws ParserException {
		//System.out.println("Atail");
		//System.out.println(lexemList);
		if(peekLexem(Lexer.LexemType.alternative)){
			parseLexem(Lexer.LexemType.alternative);
			parseAlt();
			parseAltTail();
		}
	}
	// <A> ::= <rg> <star> <rgtail>
	void parseAlt() throws ParserException{
		//System.out.println("A");
		//System.out.println(lexemList);
		parseReg();
		parseStar();
		parseRegTail();
	}
	// <rgtail>	::= <rg> <star> <rgtail> | \epsilon
	void parseRegTail() throws ParserException{
		//System.out.println("rgtail");
		//System.out.println(lexemList);
		if(peekLexem(Lexer.LexemType.lgroup)
			|| peekLexem(Lexer.LexemType.lneutr) || peekLexem(Lexer.LexemType.reggrab)
			|| peekLexem(Lexer.LexemType.wordgrab) || peekLexem(Lexer.LexemType.letter)){
			parseReg();
			parseStar();
			parseRegTail();
		}
	}
	// <rg> ::= '(' <S> ')' | '(:' <S> ')' | '(?[num])' | '\[num]' | [a-z]
	void parseReg() throws ParserException{
		//System.out.println("rg");
		//System.out.println(lexemList);
		if(peekLexem(Lexer.LexemType.lgroup)){
			parseLexem(Lexer.LexemType.lgroup);
			parseS();
			parseLexem(Lexer.LexemType.rb);
		}else if(peekLexem(Lexer.LexemType.lneutr)){
			parseLexem(Lexer.LexemType.lneutr);
			parseS();
			parseLexem(Lexer.LexemType.rb);
		}else if(peekLexem(Lexer.LexemType.reggrab)){
			parseLexem(Lexer.LexemType.reggrab);
		}else if(peekLexem(Lexer.LexemType.wordgrab)){
			parseLexem(Lexer.LexemType.wordgrab);
		}else{
			parseLexem(Lexer.LexemType.letter);
		}
	}
	// <star> ::= '*' | \epsilon
	void parseStar() throws ParserException{
		//System.out.println("star");
		//System.out.println(lexemList);
		if(peekLexem(Lexer.LexemType.star)){
			parseLexem(Lexer.LexemType.star);
			//add star
		}
	}
	
	boolean peekLexem(Lexer.LexemType t){
		return (lexemList.peek() == null && t == null) || (lexemList.peek() != null && lexemList.peek().type == t); 
	}
	Lexer.Lexem parseLexem(Lexer.LexemType t) throws ParserException{
		Lexer.Lexem l = lexemList.pop();
		if(l == null || l.type != t)
			throw new ParserException(String.format("Ожидалось %s, найдено %s", t, l));
		
		return l;
	}
	
}
class ParserException extends GrammarException{
	ParserException(String msg){
		super("parse: "+msg);
	}
}

class Regex{
	class Reg{
		int value; //'a' - 'z' - alphabetic, 1-9 - reg grab, -9 - -1 - wordgrab
		boolean star = false;
		
		public void setStar(){star = true;}
	}
	class Group extends Reg{
		//value = 0 if neutral group, value = 1-9 grab group
		LinkedList<LinkedList<Reg>> alternatives; 
	}
}