import java.util.ArrayList;
import java.util.LinkedList;

/*
[rule] ::= [NT] −> ([NT]|[T])+ [EOL]
[T] ::= [a−z]
[NT] ::= [A−Z][0−9]? | [ [A−z]+([0−9])* ]
[EOL] ::= EOL+
  
lexems
"->"
"T"
"NT"
"["
"]"
"blank"
"EOL"
*/

public class Parser{
	int parseIndex;
	String inputString;
	LinkedList<Lexer.Lexem> arr;
	
	ArrayList<NonTerm> parseResult;
	ArrayList<Term> terms;
	
	public static void main(String[] args){
		String test="A -> aA2\n A -> a\n A2 -> bb[aa47]\n\n\n [aa47]-> c\n";
		
		Parser r = new Parser();
		try{
			for(NonTerm i: r.parse(test)){
				System.out.println(i);
			}
		}catch(ParserException e){
			System.err.println(e);
		}catch(LexerException e){
			System.err.println(e);
		}
	}
	
	public ArrayList<NonTerm> parse(String msg) throws ParserException, LexerException{
		inputString=msg;
		
		Lexer l = new Lexer();
		arr = l.lex(msg);
		
		parseResult = new ArrayList<NonTerm>();
		terms = new ArrayList<Term>();
		
		while(arr.peek() != null){
			parseRule();
		}
		return parseResult;
	}
	
	NonTerm parseRule() throws ParserException{
		NonTerm rule = (NonTerm)parseLexem(Lexer.LexemType.NT);
		rule.addRule();
		
		Term tmp;
		parseLexem(Lexer.LexemType.to);
		if(checkLexem(Lexer.LexemType.T)){
			tmp=parseLexem(Lexer.LexemType.T);
		}else{
			tmp=parseLexem(Lexer.LexemType.NT);
		}
		rule.addToLastRule(tmp);
		
		while(checkLexem(Lexer.LexemType.T) || checkLexem(Lexer.LexemType.NT)){
			if(checkLexem(Lexer.LexemType.T)){
				tmp=parseLexem(Lexer.LexemType.T);
			}else{
				tmp=parseLexem(Lexer.LexemType.NT);
			}
			rule.addToLastRule(tmp);
		}
		
		parseLexem(Lexer.LexemType.EOL);
		return rule;
	}
	
	Term parseLexem(Lexer.LexemType type) throws ParserException{
		if(!checkLexem(type)){
			if(arr.peek() == null){
				throw new ParserException(String.format("При парсинге %s, найден конец файла",
					type));
			}else{
				throw new ParserException(String.format("При парсинге ожидалось %s, найден %s",
					type, arr.peek().type));
			}
		}
		Lexer.Lexem l = arr.pop();
		if(l.type == Lexer.LexemType.T || l.type==Lexer.LexemType.NT){
			String m = inputString.substring(l.start, l.end+1);
			Term t = findIfExists(m);
			if(t != null){
				return t;
			}else{
				Term tmp;
				if(l.type == Lexer.LexemType.T){
					tmp = new Term(m);
				}else{
					tmp = new NonTerm(m);
					parseResult.add((NonTerm)tmp);
				}
				terms.add(tmp);
				return tmp;
			}
		}else{
			return null;
		}
	}
	
	boolean checkLexem(Lexer.LexemType type){
		return arr.peek() != null && arr.peek().type == type;
	}
	
	Term findIfExists(String name){
		if(terms.size()==0)
			return null;
		
		for(Term i: terms){
			if(name.equals(i.getName())){
				return i;
			}
		}
		
		return null;
	}
}

class Term{
	String name;
	
	Term(String representation){
		this.name = representation;
	}
	
	public boolean isNonTerm(){
		return (this instanceof NonTerm);
	}
	
	public String getName(){return this.name;}
	
	@Override
	public String toString(){return this.name;}
	
	@Override
	public boolean equals(Object o){
		// If the object is compared with itself then return true  
        if (o == this) {
            return true;
        }
 
        /* Check if o is an instance of Complex or not
          "null instanceof [type]" also returns false */
        if (!(o instanceof Term)) {
            return false;
        }
		
		Term r = (Term) o;
		return r.getName().equals(this.getName());
	}
}

class NonTerm extends Term{
	public ArrayList<ArrayList<Term>> rewriteRules;
	
	NonTerm(String representation){
		super(representation);
		rewriteRules = new ArrayList<ArrayList<Term>>();
	}
	
	public void addRule(){
		rewriteRules.add(new ArrayList<Term>());
	}
	public void addToLastRule(Term t){
		if(rewriteRules.size() == 0){
			addRule();
		}
		rewriteRules.get(rewriteRules.size()-1).add(t);
	}
	
	@Override
	public String toString(){
		String result = this.name;
		for(int i = 0; i<rewriteRules.size(); i++){
			if(i==0){
				result+=" -> ";
			}else{
				result+="| ";
			}
			for(Term j: rewriteRules.get(i)){
				result += j.getName()+" ";
			}
		}
		return result;
	}
}

class ParserException extends Exception{
	ParserException(String msg){
		super(msg);
	}
}
