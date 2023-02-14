grammar Filters;

filter
  : Number
  ;

Number: Digit+ ('.' Digit+)? ;
fragment Digit : [0-9] ;
Whitespace : [\p{White_Space}] -> skip ;
