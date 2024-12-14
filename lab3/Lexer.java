import java.util.LinkedList;


class Lexer{
	int parseIndex = 0;
	String inputString;
	
	int startOfLine=0;
	int numberOfLine=0;
	
	public static void main(String[] args){
		String test="A -> aA2\n A2 -> bb[aa47]\n\n\n [aa47]-> c";
		
		Lexer r = new Lexer();
		try{
			for(Lexem i: r.lex(test)){
				System.out.println(String.format("%d:%d %s", i.start, i.end, i.type));
			}
		}catch(LexerException e){
			System.err.println(e);
		}
	}
	
	public LinkedList<Lexem> lex(String msg) throws LexerException{
		inputString = msg;
		LinkedList<Lexem> arr = new LinkedList<Lexem>();
		
		for(parseIndex = 0; parseIndex<msg.length();){
			char nowChar = msg.charAt(parseIndex);
			//System.out.println(nowChar + " "+nowChar);
			
			if(nowChar == '-'){ 						// ->
				arr.add(parseTO());
			}else if(Character.isLowerCase(nowChar)){	// T
				arr.add(parseT());
			}else if(Character.isUpperCase(nowChar) || nowChar == '['){	// NT
				arr.add(parseNT());
			}else if(nowChar == ' ' || nowChar == '\t'){// blank
				parseBLANK(); // blank lexems dont add to result
			}else if(nowChar == '\n'){					// EOL
				arr.add(parseEOL());
			}else{
				throw new LexerException(String.format("В строке %d, позиции %d, найден не поддерживаемый символ '%c'",
					numberOfLine+1, parseIndex-startOfLine+1, nowChar));
			}
		}
		return arr;
	}
	Lexem parseTO() throws LexerException{
		int start=parseIndex;
		parseChar('-', "->");
		checkNotEOF("->");
		parseChar('>', "->");
		return new Lexem(start, parseIndex-1, LexemType.to);
	}
	Lexem parseT() throws LexerException{
		char nowChar = inputString.charAt(parseIndex);
		if(!Character.isLowerCase(nowChar)){
			throw new LexerException(String.format("В строке %d, позиции %d, при обработке %s ожидалось '%с', найден '%c'",
				numberOfLine+1, parseIndex-startOfLine+1, "term", "a-z", nowChar));
		}
		parseIndex++;
		return new Lexem(parseIndex-1, parseIndex-1, LexemType.T);
	}
	Lexem parseNT() throws LexerException{
		int start=parseIndex;
		char nowChar = inputString.charAt(parseIndex);
		if(Character.isUpperCase(nowChar)){
			parseIndex++;
			if(parseIndex == inputString.length() || !Character.isDigit(inputString.charAt(parseIndex)))
				return new Lexem(start, parseIndex-1, LexemType.NT);
			parseIndex++;
		}else if(nowChar == '['){
			for(; parseIndex < inputString.length() && inputString.charAt(parseIndex) != ']'; parseIndex++);
			parseIndex++;
		}
		return new Lexem(start, parseIndex-1, LexemType.NT);
	}
	Lexem parseBLANK() throws LexerException{
		int start=parseIndex;
		for(;parseIndex < inputString.length() && (inputString.charAt(parseIndex) == ' ' || inputString.charAt(parseIndex) == '\t'); parseIndex++);
		
		return new Lexem(start, parseIndex-1, LexemType.blank);
	}
	Lexem parseEOL() throws LexerException{
		int start=parseIndex;
		parseChar('\n', "eol");
		for(;parseIndex < inputString.length() && inputString.charAt(parseIndex) == '\n'; parseIndex++);
		
		return new Lexem(start, parseIndex-1, LexemType.EOL);
	}

	void checkNotEOF(String type) throws LexerException{
		if(parseIndex == inputString.length()){
			throw new LexerException(String.format("В строке %d, позиции %d, при обработке %s найден конец строки",
				numberOfLine+1, parseIndex-startOfLine+1, type));
		}
	}
	void parseChar(char t, String type) throws LexerException{
		if(inputString.charAt(parseIndex) != t){
			throw new LexerException(String.format("В строке %d, позиции %d, при обработке %s ожидалось '%с', найден '%c'",
				numberOfLine+1, parseIndex-startOfLine+1, type, t, inputString.charAt(parseIndex) ));
		}
		parseIndex++;
	}
	
	enum LexemType{
		to("->"), 		// '->'
		T("term"), 		// [a-z]
		NT("nonterm"),	// [A−Z][0−9]?|[[A−z]+([0−9])∗]
		blank("blank"),	// space | tab
		EOL("eol");		// (\n)+
		
		String msg;
		LexemType(String s){msg=s;}
		public String toString(){return msg;}
	}
	public class Lexem{
		public int start;
		public int end;
		public LexemType type;
		
		Lexem(int s, int e, LexemType t){start=s;end=e;type=t;}
	}
}


class LexerException extends Exception{
	LexerException(String msg){
		super(msg);
	}
}