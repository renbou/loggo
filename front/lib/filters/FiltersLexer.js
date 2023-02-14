// Generated from java-escape by ANTLR 4.11.1
// jshint ignore: start
import antlr4 from 'antlr4';


const serializedATN = [4,0,2,26,6,-1,2,0,7,0,2,1,7,1,2,2,7,2,1,0,4,0,9,8,
0,11,0,12,0,10,1,0,1,0,4,0,15,8,0,11,0,12,0,16,3,0,19,8,0,1,1,1,1,1,2,1,
2,1,2,1,2,0,0,3,1,1,3,0,5,2,1,0,2,1,0,48,57,10,0,9,13,32,32,133,133,160,
160,5760,5760,8192,8202,8232,8233,8239,8239,8287,8287,12288,12288,27,0,1,
1,0,0,0,0,5,1,0,0,0,1,8,1,0,0,0,3,20,1,0,0,0,5,22,1,0,0,0,7,9,3,3,1,0,8,
7,1,0,0,0,9,10,1,0,0,0,10,8,1,0,0,0,10,11,1,0,0,0,11,18,1,0,0,0,12,14,5,
46,0,0,13,15,3,3,1,0,14,13,1,0,0,0,15,16,1,0,0,0,16,14,1,0,0,0,16,17,1,0,
0,0,17,19,1,0,0,0,18,12,1,0,0,0,18,19,1,0,0,0,19,2,1,0,0,0,20,21,7,0,0,0,
21,4,1,0,0,0,22,23,7,1,0,0,23,24,1,0,0,0,24,25,6,2,0,0,25,6,1,0,0,0,4,0,
10,16,18,1,6,0,0];


const atn = new antlr4.atn.ATNDeserializer().deserialize(serializedATN);

const decisionsToDFA = atn.decisionToState.map( (ds, index) => new antlr4.dfa.DFA(ds, index) );

export default class FiltersLexer extends antlr4.Lexer {

    static grammarFileName = "Filters.g4";
    static channelNames = [ "DEFAULT_TOKEN_CHANNEL", "HIDDEN" ];
	static modeNames = [ "DEFAULT_MODE" ];
	static literalNames = [  ];
	static symbolicNames = [ null, "Number", "Whitespace" ];
	static ruleNames = [ "Number", "Digit", "Whitespace" ];

    constructor(input) {
        super(input)
        this._interp = new antlr4.atn.LexerATNSimulator(this, atn, decisionsToDFA, new antlr4.PredictionContextCache());
    }

    get atn() {
        return atn;
    }
}

FiltersLexer.EOF = antlr4.Token.EOF;
FiltersLexer.Number = 1;
FiltersLexer.Whitespace = 2;



