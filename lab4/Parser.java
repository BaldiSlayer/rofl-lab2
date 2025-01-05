import java.util.LinkedList;
import java.util.ListIterator;

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
	
	Regex resultRegex;
	Group nowRegex;
	
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
		
		resultRegex = new Regex();
		parseS(resultRegex);
		
		if(lexemList.peek() != null){
			throw new ParserException(String.format("В ходе разбора строка была разобрана не полностью. Оставшиеся лексемы %s",
				lexemList));
		}
		return resultRegex;
	}
	
	// <S> ::= <A> <Atail>
	void parseS(Group nowRegex) throws ParserException{
		//System.out.println("S");
		//System.out.println(lexemList);
		
		nowRegex.addAlternative();
		
		parseAlt(nowRegex);
		parseAltTail(nowRegex);
	}
	// <Atail> ::= '|' <A> <Atail> | \epsilon
	void parseAltTail(Group nowRegex) throws ParserException {
		//System.out.println("Atail");
		//System.out.println(lexemList);
		if(peekLexem(Lexer.LexemType.alternative)){
			nowRegex.addAlternative();
			parseLexem(Lexer.LexemType.alternative);
			parseAlt(nowRegex);
			parseAltTail(nowRegex);
		}
	}
	// <A> ::= <rg> <star> <rgtail>
	void parseAlt(Group nowRegex) throws ParserException{
		//System.out.println("A");
		//System.out.println(lexemList);
		Reg r = parseReg();
		parseStar(r);
		nowRegex.addReg(r);
		
		parseRegTail(nowRegex);
	}
	// <rgtail>	::= <rg> <star> <rgtail> | \epsilon
	void parseRegTail(Group nowRegex) throws ParserException{
		//System.out.println("rgtail");
		//System.out.println(lexemList);
		if(peekLexem(Lexer.LexemType.lgroup)
			|| peekLexem(Lexer.LexemType.lneutr) || peekLexem(Lexer.LexemType.reggrab)
			|| peekLexem(Lexer.LexemType.wordgrab) || peekLexem(Lexer.LexemType.letter)){
			Reg r = parseReg();
			parseStar(r);
			nowRegex.addReg(r);
		
			parseRegTail(nowRegex);
		}
	}
	// <rg> ::= '(' <S> ')' | '(:' <S> ')' | '(?[num])' | '\[num]' | [a-z]
	Reg parseReg() throws ParserException{
		//System.out.println("rg");
		//System.out.println(lexemList);
		Reg r;
		if(peekLexem(Lexer.LexemType.lgroup)){
			Lexer.Lexem l = parseLexem(Lexer.LexemType.lgroup);
			Group g = new Group(l.value, l.start);
			parseS(g);
			parseLexem(Lexer.LexemType.rb);
			r = g;
		}else if(peekLexem(Lexer.LexemType.lneutr)){
			Group g = new Group(0, parseLexem(Lexer.LexemType.lneutr).start);
			parseS(g);
			parseLexem(Lexer.LexemType.rb);
			r = g;
		}else if(peekLexem(Lexer.LexemType.reggrab)){
			Lexer.Lexem l = parseLexem(Lexer.LexemType.reggrab);
			r = new Reg(l.value, l.start);
		}else if(peekLexem(Lexer.LexemType.wordgrab)){
			Lexer.Lexem l = parseLexem(Lexer.LexemType.wordgrab);
			r = new Reg(-l.value, l.start);
		}else{
			Lexer.Lexem l = parseLexem(Lexer.LexemType.letter);
			r = new Reg(l.value, l.start);
		}
		return r;
	}
	// <star> ::= '*' | \epsilon
	void parseStar(Reg r) throws ParserException{
		//System.out.println("star");
		//System.out.println(lexemList);
		if(peekLexem(Lexer.LexemType.star)){
			parseLexem(Lexer.LexemType.star);
			//add star
			r.setStar();
		}
	}
	
	boolean peekLexem(Lexer.LexemType t){
		return (lexemList.peek() == null && t == null) || (lexemList.peek() != null && lexemList.peek().type == t); 
	}
	Lexer.Lexem parseLexem(Lexer.LexemType t) throws ParserException{
		if(lexemList.peek() == null)
			throw new ParserException(String.format("Ожидалось %s, найдено конец ввода", t));
		
		Lexer.Lexem l = lexemList.pop();
		if(l.type != t)
			throw new ParserException(String.format("Ожидалось %s, найдено %s", t, l));
		
		return l;
	}
	
}
class ParserException extends GrammarException{
	ParserException(String msg){
		super(msg);
	}
}

class Reg{
	public int startIndex;
	public int value; //'a' - 'z' - alphabetic, 1-9 - reg grab, -9 - -1 - wordgrab
	public boolean star = false;
		
	Reg(int v, int i){
		value = v;
		startIndex = i;
	}
	/*
	Reg(char c){
		value = (int)c;
	}*/
	
	public void setStar(){star = true;}
	
	@Override
	public String toString(){
		String r;
		if(value < 0){
			r = "\\"+(-value);
		}else if(value >= 'a'){
			r = String.valueOf((char)value);
		}else{
			r = "(?"+value+")";
		}
		if(star){
			r += "*";
		}
		return r;
	}
}
class Group extends Reg{
	//value = 0 if neutral group, value = 1-9 grab group
	public LinkedList<LinkedList<Reg>> alternatives; 
	//public boolean canBeWordGrab = false;
	
	Group(int v, int i){
		super(v, i);
		alternatives = new LinkedList<LinkedList<Reg>>();
	}
	public void addAlternative(){
		alternatives.add(new LinkedList<Reg>());
	}
	public LinkedList<Reg> getLastAlternative(){
		return alternatives.getLast();
	}
	public Reg getLastReg(){
		return alternatives.getLast().getLast();
	}
	public void addReg(Reg r){
		alternatives.getLast().add(r);
	}
	
	String alternativesToString(){
		String r = "";
		for(Reg j : alternatives.peek()){
			r += j.toString();
		}
		ListIterator<LinkedList<Reg>> l = alternatives.listIterator(1);
		while(l.hasNext()){
			LinkedList<Reg> i = l.next();
			r += "|";
			for(Reg j : i){
				r += j.toString();
			}
		}
		return r;
	}
	
	@Override
	public String toString(){
		String r;
		if(value == 0){
			r = "(?:";
		}else{
			r = "(";
		}
		r += alternativesToString() + ")";
		if(star)
			r += "*";
		return r;
	}
}

class Regex extends Group{
	
	Regex(){
		super(0,0);
	}
	
	@Override
	public void setStar(){return;}
	
	@Override
	public String toString(){
		return alternativesToString();
	}
}