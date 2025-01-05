import java.util.LinkedList;

class Main{
	public static void main(String[] args){
		LinkedList<Integer> l = new LinkedList<Integer>();
		
		l.add(1);
		l.add(2);
		l.add(3);
		
		System.out.println(l.size());
	}
}

class GrammarException extends Exception{
	GrammarException(String msg){
		super(msg);
	}
}