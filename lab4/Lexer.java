/*
[rg] ::= [rg][rg] | [rg]|[rg]
 | ([rg]) | (?: [rg] ) | [rg]*
 | (?[num]) | [a−z]
 [num] ::= [1−9]
 [rg] ::= \[num]


lexems
letter		// [a-z]
lgroup		// (
lneutr		// (?:
rb 			// )
reggrab		// (? [1-9] )
wordgrab	// \[1-9]
star		// *
alternative // |
*/


import java.util.LinkedList;


class Lexer{
	int parseIndex = 0;
	String inputString;
	
	int groupCounter = 1; // считает номер очередной группы захвата. При выходе за число 9 - бросаем ошибку
	
	public static void main(String[] args){
		String test="((?: a | bc* | c)a(?1))\\1*  |\td";
		
		Lexer r = new Lexer();
		try{
			for(Lexem i: r.lex(test)){
				System.out.println(i);
			}
		}catch(LexerException e){
			System.err.println(e);
		}
	}
	
	public LinkedList<Lexem> lex(String msg) throws LexerException{
		inputString = msg;
		LinkedList<Lexem> arr = new LinkedList<Lexem>();
		//arr.add(null); // aka EOF
		
		for(parseIndex = 0; parseIndex<msg.length();){
			char nowChar = msg.charAt(parseIndex);
			//System.out.println(nowChar + " "+nowChar)
			if(isSpace(nowChar)){ // space \t 
				//nothing
				parseIndex++;
			}else if(Character.isAlphabetic(nowChar)){ // [a-z]
				arr.add(new Lexem(parseIndex, parseIndex++, LexemType.letter));
			}else if(nowChar == ')'){
				arr.add(new Lexem(parseIndex, parseIndex++, LexemType.rb));
			}else if(nowChar == '|'){
				arr.add(new Lexem(parseIndex, parseIndex++, LexemType.alternative));
			}else if(nowChar == '*'){
				arr.add(new Lexem(parseIndex, parseIndex++, LexemType.star));
			}else if(nowChar == '\\'){ // \[num]
				arr.add(parseWordGrab());
			}else if(nowChar == '('){// ( | (?: | (? [num] )
				// (? [num]) накладывает больше всего условий - парсим в начале
				int start = parseIndex;
				try{
					arr.add(parseRegGrab());
					continue;
				}catch(LexerException e){}
				// (?: после пытаемся распарсить это
				parseIndex = start;
				try{
					arr.add(parseFreeBracket());
					continue;
				}catch(LexerException e){}
				// после парсим только (
				parseIndex = start;
				arr.add(parseGroup());
			}else{
				throw new LexerException(String.format("На позиции %d, найден не поддерживаемый символ '%c'",
					parseIndex+1, nowChar));
			}
		}
		return arr;
	}
	
	private Lexem parseWordGrab() throws LexerException{	// \[1-9]
		parseIndex++;
		
		Lexem r = new Lexem(parseIndex-1, parseIndex, LexemType.wordgrab);
		r.value = parseNum();
		return r;
	}
	private Lexem parseRegGrab() throws LexerException{	// (? [1-9] )
		parseIndex++;
		parseChar("(?[1-9])", '?');
		Lexem r = new Lexem(parseIndex-1, parseIndex, LexemType.reggrab);
		r.value = parseNum();
		
		parseChar("(?[1-9])", ')');
		return r;
	}
	private Lexem parseFreeBracket() throws LexerException{	// (?:
		parseIndex++;
		parseChar("free bracket", '?');
		parseChar("free bracket", ':');
		return new Lexem(parseIndex-3, parseIndex-1, LexemType.lneutr);
	}
	private Lexem parseGroup() throws LexerException{		// ( при входе в функцию известно, что на позиции стоит (
		if(groupCounter > 9 && groupCounter < 1)
			throw new LexerException("Превышено максимальное количество групп захвата: от 1 до 9");
		Lexem r = new Lexem(parseIndex, parseIndex++, LexemType.lgroup);
		r.value = groupCounter++;
		return r;
	}
	
	private boolean isSpace(char a){
		return (a == ' ') || (a == '\t');
	}
	private void tryLength(String parse, String expect) throws LexerException{
		if(parseIndex == inputString.length())
			throw new LexerException(String.format("На позиции %d при обработке %s ожидалась %s, найден конец строки",
				parseIndex+1, parse, expect));
	}
	private void tryLength(String parse, char expect) throws LexerException{
		if(parseIndex == inputString.length())
			throw new LexerException(String.format("На позиции %d при обработке %s ожидалась '%c', найден конец строки",
				parseIndex+1, parse, expect));
	}
	
	private void parseChar(String parse, char c) throws LexerException{
		tryLength(parse, c);
		
		char nowChar = inputString.charAt(parseIndex);
		if(nowChar != c){
			throw new LexerException(String.format("На позиции %d при обработке %s ожидалась цифра '%c', найден '%c'",
				parseIndex+1, parse, c, nowChar));
		}
		parseIndex++;
	}
	private int parseNum() throws LexerException{
		tryLength("\\[1-9]", "цифра");
		
		char nowChar = inputString.charAt(parseIndex);
		if(nowChar >= '1' && nowChar <= '9'){
			int r = (int)(nowChar - '0');
			parseIndex++;
			return r;
		}else{
			throw new LexerException(String.format("На позиции %d при обработке \\[1-9] ожидалась цифра от 1 до 9, найден '%c'",
				parseIndex+1, nowChar));
		}
	}
	
	enum LexemType{
		letter("[a-z]"),		// [a-z]
		lgroup("("),			// (
		lneutr("(?:"),			// (?:
		rb(")"),				// )
		reggrab("REGgrab"),		// (? [1-9] )
		wordgrab("WORDgrab"),	// \[1-9]
		star("STAR"),			// *
		alternative("ALT"); 	// |
		
		
		String msg;
		LexemType(String s){msg=s;}
		public String toString(){return msg;}
	}
	public class Lexem{
		public int start;
		public int end;
		public LexemType type;
		public Integer value = null; // специальная константа для упрощения дальнейшего парсинга
		
		Lexem(int s, int e, LexemType t){
			start=s;end=e;type=t;
			
			if(type == Lexer.LexemType.letter){
				value = (int)inputString.charAt(s);
			}
		}
		public String toString(){return String.format("%d:%d %s", start, end, type);}
	}
}


class LexerException extends GrammarException{
	LexerException(String msg){
		super(msg);
	}
}