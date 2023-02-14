// Generated from java-escape by ANTLR 4.11.1
// jshint ignore: start
import antlr4 from 'antlr4';
const serializedATN = [4,1,2,5,2,0,7,0,1,0,1,0,1,0,0,0,1,0,0,0,3,0,2,1,0,
0,0,2,3,5,1,0,0,3,1,1,0,0,0,0];


const atn = new antlr4.atn.ATNDeserializer().deserialize(serializedATN);

const decisionsToDFA = atn.decisionToState.map( (ds, index) => new antlr4.dfa.DFA(ds, index) );

const sharedContextCache = new antlr4.PredictionContextCache();

export default class FiltersParser extends antlr4.Parser {

    static grammarFileName = "java-escape";
    static literalNames = [  ];
    static symbolicNames = [ null, "Number", "Whitespace" ];
    static ruleNames = [ "filter" ];

    constructor(input) {
        super(input);
        this._interp = new antlr4.atn.ParserATNSimulator(this, atn, decisionsToDFA, sharedContextCache);
        this.ruleNames = FiltersParser.ruleNames;
        this.literalNames = FiltersParser.literalNames;
        this.symbolicNames = FiltersParser.symbolicNames;
    }

    get atn() {
        return atn;
    }



	filter() {
	    let localctx = new FilterContext(this, this._ctx, this.state);
	    this.enterRule(localctx, 0, FiltersParser.RULE_filter);
	    try {
	        this.enterOuterAlt(localctx, 1);
	        this.state = 2;
	        this.match(FiltersParser.Number);
	    } catch (re) {
	    	if(re instanceof antlr4.error.RecognitionException) {
		        localctx.exception = re;
		        this._errHandler.reportError(this, re);
		        this._errHandler.recover(this, re);
		    } else {
		    	throw re;
		    }
	    } finally {
	        this.exitRule();
	    }
	    return localctx;
	}


}

FiltersParser.EOF = antlr4.Token.EOF;
FiltersParser.Number = 1;
FiltersParser.Whitespace = 2;

FiltersParser.RULE_filter = 0;

class FilterContext extends antlr4.ParserRuleContext {

    constructor(parser, parent, invokingState) {
        if(parent===undefined) {
            parent = null;
        }
        if(invokingState===undefined || invokingState===null) {
            invokingState = -1;
        }
        super(parent, invokingState);
        this.parser = parser;
        this.ruleIndex = FiltersParser.RULE_filter;
    }

	Number() {
	    return this.getToken(FiltersParser.Number, 0);
	};


}




FiltersParser.FilterContext = FilterContext; 
