/*
[rg] ::= [rg][rg] | [rg]|[rg]
 | ([rg]) | (? :[rg]) | [rg]*
 | ([num]) | [a−z]
 [num] ::= [1−9]
 [rg] ::= \[num]


lexems
letter
num
lb
rb
la
star
*/


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
			//System.out.println(nowChar + " "+nowChar)
		}
		return arr;
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