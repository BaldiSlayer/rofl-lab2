import java.util.LinkedList; // for stack


class PDA{
	boolean[] states; // isFinal
	PDARule[] rules;
	
	LinkedList<String> stack;
	
	PDA(){}
	
	public boolean perform(String msg){
		return false;
	}
	
	public String toString(){ // DOT representation
		String result = "digraph{\n\tnode[shape=circle]\n\tpoint[shape=point]\n\t%s\npoint->0%s\n}";
		
		return String.format(result, getFinal(), getTransitions());
	}
	
	String getFinal(){
		String result = "";
		for(int i =0; i< states.length(); i++){
			if(states[i]){
				if(result.equal("")) result = String.format("%d", i);
				else result += ", "+i;
			}
		}
		if(!result.equal("")){
			result += "[shape=doublecircle]";
		}
		return result;
	}
	String getTransitions(){
		String result = "";
		return result;
	}
	
	class PDARule{
		char goBy;
		String popSymbol;
		String[] putSymbols;
		int goTo;
		
		public boolean check(char c){
			return c == goBy && ((stack.peek() == null && popSymbol.equal("Z0")) || popSymbol.equals(stack.peek()));
		}
	}
}