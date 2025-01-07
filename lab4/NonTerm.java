import java.util.ArrayList;

class NonTerm extends Term{
	public ArrayList<ArrayList<Term>> rewriteRules;
	
	NonTerm(String representation){
		super(representation);
		rewriteRules = new ArrayList<ArrayList<Term>>();
	}
	
	public void addRule(){
		rewriteRules.add(new ArrayList<Term>());
		//return rewriteRules.get(rewriteRules.size()-1);
	}
	/*public ArrayList<Term> getLastTerm(){
		return rewriteRules.get(rewriteRules.size()-1);
	} */
	
	public void addToLastRule(Term t){
		if(rewriteRules.size() == 0){
			addRule();
		}
		rewriteRules.get(rewriteRules.size()-1).add(t);
	}
	public String alternativeString(int i){
		String result = "";
		for(Term j: rewriteRules.get(i)){
			result += j.getName();
		}
		return result;
	}
	
	public boolean checkEqv(NonTerm nterm){
		System.out.println(this);
		System.out.println(nterm);
		if(rewriteRules.size() != nterm.rewriteRules.size())
			return false;
		
		for(ArrayList<Term> i : nterm.rewriteRules){
			boolean r = false;
			for(ArrayList<Term> j : this.rewriteRules){
				if(i.equals(j)){
					r = true;
					break;
				}
			}
			if(!r)
				return false;
		}
		return true;
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