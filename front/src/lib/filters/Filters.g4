grammar Filters;

filter returns [filterKind]
  : expression EOF { $filterKind = $expression.filterKind; }
  ;

expression returns [filterKind]
  : 'NOT' expression { $filterKind = "not"; }
  | expression 'AND' expression { $filterKind = "and"; }
  | expression 'OR' expression { $filterKind = "or"; }
  | StringLiteral { $filterKind = "text"; }
  | FieldName ':' StringLiteral { $filterKind = "scoped"; }
  ;

StringLiteral : '"' ('\\' . | ~["\\])* '"' ;
FieldName : [\p{N}\p{L}_]+ ;
