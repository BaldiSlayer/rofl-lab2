import java.util.LinkedList;
import java.util.Scanner;

class Main{
	public static void main(String[] args){
		Scanner scan = new Scanner(System.in);
		Parser p = new Parser();
		
		while(true){
			System.out.print("enter regex: ");
			String test = scan.nextLine();
			
			try{
				Regex r = p.parse(test);
				System.out.println(r);
				FrameGrammar fg = new FrameGrammar(r);
			}catch(GrammarException e){
				System.out.println(e);
			}
			
			System.out.println("---------------------");
		}
	}
}

class GrammarException extends Exception{
	GrammarException(String msg){
		super(msg);
	}
}