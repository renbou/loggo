// Generated from java-escape by ANTLR 4.11.1
// jshint ignore: start
import antlr4 from 'antlr4';
const serializedATN = [4,1,6,37,2,0,7,0,2,1,7,1,1,0,1,0,1,0,1,0,1,1,1,1,
1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,3,1,20,8,1,1,1,1,1,1,1,1,1,1,1,1,1,1,
1,1,1,1,1,1,1,5,1,32,8,1,10,1,12,1,35,9,1,1,1,0,1,2,2,0,2,0,0,38,0,4,1,0,
0,0,2,19,1,0,0,0,4,5,3,2,1,0,5,6,5,0,0,1,6,7,6,0,-1,0,7,1,1,0,0,0,8,9,6,
1,-1,0,9,10,5,1,0,0,10,11,3,2,1,5,11,12,6,1,-1,0,12,20,1,0,0,0,13,14,5,5,
0,0,14,20,6,1,-1,0,15,16,5,6,0,0,16,17,5,4,0,0,17,18,5,5,0,0,18,20,6,1,-1,
0,19,8,1,0,0,0,19,13,1,0,0,0,19,15,1,0,0,0,20,33,1,0,0,0,21,22,10,4,0,0,
22,23,5,2,0,0,23,24,3,2,1,5,24,25,6,1,-1,0,25,32,1,0,0,0,26,27,10,3,0,0,
27,28,5,3,0,0,28,29,3,2,1,4,29,30,6,1,-1,0,30,32,1,0,0,0,31,21,1,0,0,0,31,
26,1,0,0,0,32,35,1,0,0,0,33,31,1,0,0,0,33,34,1,0,0,0,34,3,1,0,0,0,35,33,
1,0,0,0,3,19,31,33];


const atn = new antlr4.atn.ATNDeserializer().deserialize(serializedATN);

const decisionsToDFA = atn.decisionToState.map( (ds, index) => new antlr4.dfa.DFA(ds, index) );

const sharedContextCache = new antlr4.PredictionContextCache();

export default class FiltersParser extends antlr4.Parser {

    static grammarFileName = "java-escape";
    static literalNames = [ null, "'NOT'", "'AND'", "'OR'", "':'" ];
    static symbolicNames = [ null, null, null, null, null, "StringLiteral", 
                             "FieldName" ];
    static ruleNames = [ "filter", "expression" ];

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

    sempred(localctx, ruleIndex, predIndex) {
    	switch(ruleIndex) {
    	case 1:
    	    		return this.expression_sempred(localctx, predIndex);
        default:
            throw "No predicate with index:" + ruleIndex;
       }
    }

    expression_sempred(localctx, predIndex) {
    	switch(predIndex) {
    		case 0:
    			return this.precpred(this._ctx, 4);
    		case 1:
    			return this.precpred(this._ctx, 3);
    		default:
    			throw "No predicate with index:" + predIndex;
    	}
    };




	filter() {
	    let localctx = new FilterContext(this, this._ctx, this.state);
	    this.enterRule(localctx, 0, FiltersParser.RULE_filter);
	    try {
	        this.enterOuterAlt(localctx, 1);
	        this.state = 4;
	        localctx._expression = this.expression(0);
	        this.state = 5;
	        this.match(FiltersParser.EOF);
	         localctx.filterKind =  localctx._expression.filterKind 
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


	expression(_p) {
		if(_p===undefined) {
		    _p = 0;
		}
	    const _parentctx = this._ctx;
	    const _parentState = this.state;
	    let localctx = new ExpressionContext(this, this._ctx, _parentState);
	    let _prevctx = localctx;
	    const _startState = 2;
	    this.enterRecursionRule(localctx, 2, FiltersParser.RULE_expression, _p);
	    try {
	        this.enterOuterAlt(localctx, 1);
	        this.state = 19;
	        this._errHandler.sync(this);
	        switch(this._input.LA(1)) {
	        case 1:
	            this.state = 9;
	            this.match(FiltersParser.T__0);
	            this.state = 10;
	            this.expression(5);
	             localctx.filterKind =  "not" 
	            break;
	        case 5:
	            this.state = 13;
	            this.match(FiltersParser.StringLiteral);
	             localctx.filterKind =  "text" 
	            break;
	        case 6:
	            this.state = 15;
	            this.match(FiltersParser.FieldName);
	            this.state = 16;
	            this.match(FiltersParser.T__3);
	            this.state = 17;
	            this.match(FiltersParser.StringLiteral);
	             localctx.filterKind =  "scoped" 
	            break;
	        default:
	            throw new antlr4.error.NoViableAltException(this);
	        }
	        this._ctx.stop = this._input.LT(-1);
	        this.state = 33;
	        this._errHandler.sync(this);
	        var _alt = this._interp.adaptivePredict(this._input,2,this._ctx)
	        while(_alt!=2 && _alt!=antlr4.atn.ATN.INVALID_ALT_NUMBER) {
	            if(_alt===1) {
	                if(this._parseListeners!==null) {
	                    this.triggerExitRuleEvent();
	                }
	                _prevctx = localctx;
	                this.state = 31;
	                this._errHandler.sync(this);
	                var la_ = this._interp.adaptivePredict(this._input,1,this._ctx);
	                switch(la_) {
	                case 1:
	                    localctx = new ExpressionContext(this, _parentctx, _parentState);
	                    this.pushNewRecursionContext(localctx, _startState, FiltersParser.RULE_expression);
	                    this.state = 21;
	                    if (!( this.precpred(this._ctx, 4))) {
	                        throw new antlr4.error.FailedPredicateException(this, "this.precpred(this._ctx, 4)");
	                    }
	                    this.state = 22;
	                    this.match(FiltersParser.T__1);
	                    this.state = 23;
	                    this.expression(5);
	                     localctx.filterKind =  "and" 
	                    break;

	                case 2:
	                    localctx = new ExpressionContext(this, _parentctx, _parentState);
	                    this.pushNewRecursionContext(localctx, _startState, FiltersParser.RULE_expression);
	                    this.state = 26;
	                    if (!( this.precpred(this._ctx, 3))) {
	                        throw new antlr4.error.FailedPredicateException(this, "this.precpred(this._ctx, 3)");
	                    }
	                    this.state = 27;
	                    this.match(FiltersParser.T__2);
	                    this.state = 28;
	                    this.expression(4);
	                     localctx.filterKind =  "or" 
	                    break;

	                } 
	            }
	            this.state = 35;
	            this._errHandler.sync(this);
	            _alt = this._interp.adaptivePredict(this._input,2,this._ctx);
	        }

	    } catch( error) {
	        if(error instanceof antlr4.error.RecognitionException) {
		        localctx.exception = error;
		        this._errHandler.reportError(this, error);
		        this._errHandler.recover(this, error);
		    } else {
		    	throw error;
		    }
	    } finally {
	        this.unrollRecursionContexts(_parentctx)
	    }
	    return localctx;
	}


}

FiltersParser.EOF = antlr4.Token.EOF;
FiltersParser.T__0 = 1;
FiltersParser.T__1 = 2;
FiltersParser.T__2 = 3;
FiltersParser.T__3 = 4;
FiltersParser.StringLiteral = 5;
FiltersParser.FieldName = 6;

FiltersParser.RULE_filter = 0;
FiltersParser.RULE_expression = 1;

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
        this.filterKind = null
        this._expression = null; // ExpressionContext
    }

	expression() {
	    return this.getTypedRuleContext(ExpressionContext,0);
	};

	EOF() {
	    return this.getToken(FiltersParser.EOF, 0);
	};


}



class ExpressionContext extends antlr4.ParserRuleContext {

    constructor(parser, parent, invokingState) {
        if(parent===undefined) {
            parent = null;
        }
        if(invokingState===undefined || invokingState===null) {
            invokingState = -1;
        }
        super(parent, invokingState);
        this.parser = parser;
        this.ruleIndex = FiltersParser.RULE_expression;
        this.filterKind = null
    }

	expression = function(i) {
	    if(i===undefined) {
	        i = null;
	    }
	    if(i===null) {
	        return this.getTypedRuleContexts(ExpressionContext);
	    } else {
	        return this.getTypedRuleContext(ExpressionContext,i);
	    }
	};

	StringLiteral() {
	    return this.getToken(FiltersParser.StringLiteral, 0);
	};

	FieldName() {
	    return this.getToken(FiltersParser.FieldName, 0);
	};


}




FiltersParser.FilterContext = FilterContext; 
FiltersParser.ExpressionContext = ExpressionContext; 
