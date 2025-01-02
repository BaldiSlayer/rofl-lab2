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
		String test="((?: a | bc* | c)a(?1))\1*  |\td";
		
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
		arr.add(null); // aka EOF
		
		for(parseIndex = 0; parseIndex<msg.length();){
			char nowChar = msg.charAt(parseIndex);
			//System.out.println(nowChar + " "+nowChar)
			if(isSpace(nowChar)){ // space \t 
				//nothing
				parseIndex++;
			}else if(Char.isAlphabet(nowChar)){ // [a-z]
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
				// (?: после пытаемся распарсить это
				// после парсим только (
			}else{
				throw new LexerException(String.format("На позиции %d, найден не поддерживаемый символ '%c'",
					parseIndex+1, nowChar));
			}
		}
		return arr;
	}
	
	private Lexem parseWordGrab(){	// \[1-9]
		parseIndex++;
		if(parseIndex == inputString.length())
			throw LexerException(String.format("На позиции %d при обработке \[1-9] ожидалась цифра, найден конец строки", parseIndex+1);
		
		char nowChar = inputString.charAt(parseIndex);
		if(nowChar > '1' && nowChar < '9'){
			Lexem r = new Lexem(parseIndex-1, parseIndex++, LexemType.wordgrab);
			r.value = (int)(nowChar - '0');
			return r;
		}else{
			throw LexerException(String.format("На позиции %d при обработке \[1-9] ожидалась цифра от 0 до 9, найден '%с'",
				parseIndex+1, nowChar);
		}
	}
	private Lexem parseRegGrab(){	// (? [1-9] )
		
	}
	private Lexem parseFreeBracket(){	// (?:
		
	}
	private Lexem parseGroup(){		// ( при входе в функцию известно, что на позиции стоит (
		if(groupCounter > 9 && groupCounter < 1)
			throw new LexerException("Превышено максимальное количество групп захвата: от 1 до 9");
		Lexem r = new Lexem(parseIndex, parseIndex++, LexemType.lgroup);
		r.value = groupCounter++;
		return r;
	}
	
	private boolean isSpace(char a){
		return (a == ' ') || (a == '\t');
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
		
		Lexem(int s, int e, LexemType t){start=s;end=e;type=t;}
		public String toString(){return String.format("%d:%d %s", start, end, type);}
	}
}


class LexerException extends Exception{
	LexerException(String msg){
		super("lexer: "+msg);
	}
}