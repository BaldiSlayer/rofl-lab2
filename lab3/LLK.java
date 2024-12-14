import java.util.ArrayList;
import java.util.LinkedList;

class LLK{
	ArrayList<NonTerm> rules;

	LLK(ArrayList<NonTerm> llk) throws LLKException{
		rules=llk;
		toCorrect();
	}
	
	public ArrayList<NonTerm> getRules(){return rules;}

	public void toCorrect() throws LLKException{
		LinkedList<NonTerm> deleteList = new LinkedList<NonTerm>();
		for(NonTerm rule : rules){
			if(rule.rewriteRules.size()==0)// throw new LLKException(String.format(
				deleteList.push(rule);//"Отсутсвует раскрытие нетерминала %s", rule.getName()));
		}
		for(NonTerm rule : deleteList){
			rules.remove(rule);
		}

		for(NonTerm rule : rules){
			for(ArrayList<Term> i: rule.rewriteRules){
				for(Term t: i){
					if(t.isNonTerm() && !rules.contains(t))
						throw new LLKException(String.format(
							"В правиле переписывание %s присутсвует нетерминал %s без правил переписывания",rule, t));
				}
			}
		}
	}

	public boolean checkLLK(int k) throws LLKException{
		if(k < 0 || k > 3) throw new LLKException(
			String.format("Неправильный параметр LL(k):%d", k));
		return false;
	}


	public String toString(){ // to table first follow
		String result = "";
		for(int i = 0; i< rules.size(); i++){
			if(i != 0){
				result += "\n";
			}
			result += rules.get(i);
		}
		return result;
	}
	

	public LLK toHomsk(){
		
	}
	/*public PDA toPDA(){
		return null;
	}*/

	public static void main(String[] args){
		String test = "S->aSa\nS->b\n";//"S -> SS\nS->a\n";

		Parser r = new Parser();
		try{
			LLK system = new LLK(r.parse(test));
			System.out.println(system);
		}catch(Exception e){
			System.err.println(e);
		}
	}
	
	class LLKCheck{
		ArrayList<NonTerm> ruleList;
		ArrayList<ArrayList<LLKResult>> LLKtable;
		ArrayList<Integer> DFScolor;
	
		LLKCheck(LLK gram){
			ruleList = gram.getRules();
		}
	
		public boolean isLLk(int k){
			this.getFirstK(k);
			
			return false;
		}
	
		public ArrayList<ArrayList<NonTerm>> getRecursiveDependences(){
			ArrayList<Object> result = new ArrayList<ArrayList<NonTerm>>();
			for(NonTerm nt: ruleList){
				ArrayList<NonTerm> tmp = new ArrayList<NonTerm>();
				for(ArrayList<Term> test: nt.rewriteRules){
					for(Term i: test){
						if(i.isNonTerm()){
							tmp.add(i);
						}
					}
				}
				result.add(tmp);
			}
			return result;
		}
		

		public ArrayList<ArrayList<LLKResult>> getFirstK(int k){
			ArrayList<ArrayList<NonTerm>> rd = this.getRecursiveDependences();
			
			LLKtable = new ArrayList<ArrayList<LLKResult>>(ruleList.size());
			ArrayList<Integer> DFScolor = new ArrayList<integer>(ruleList.size());
			for(int i = 0 ; i< ruleList.size(); i++){
				result.add(null);
				DFScolor.add(0);
			}
			
			for(int i = 0; i< ruleList.size(); i++){ // dfs check
				if(result.get(i) == null){
					result.add(DFS(ruleList.get(i), rd));
				}
			}
			return result;
		}
		
		ArrayList<LLKResult> DFS(NonTerm t, ArrayList<ArrayList<NonTerm>> graph){
			
		}
	
		class LLKResult{
			public String pref;
			public int alternative;
			
			LLKResult(String s, int a){pref = s; alternative=a;}
		}
	}
}

class LLKException extends Exception{
	LLKException(String msg){
		super(msg);
	}
}