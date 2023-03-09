package parser

// Expr -> AssignExpr | (Expr, AssignExpr)
// AssignExpr -> CondExpr | YieldExpr | ArrowFunc | LeftHandSideExpr | LeftHandSideExpr "=" AssignExpr | LeftHandSideExpr ("=" | "+=" | "*=" | "/=" | "|=" | "&=" | "&&=") AssignExpr

// ..branching off CondExpr ::
//
// CondExpr -> LogOrExpr | (LogOrExpr "?" Expr ":" AssignExpr)*
// LogOrExpr -> LogAndExpr | (LogOrExpr "||" LogAndExpr)*
// LogOrExpr -> (LogOrExpr "||" LogAndExpr)

// LogAndExpr -> BitOrExpr | (LogAndExpr "&&" BitOrExpr)*
// BitOrExpr -> BitXorExpr | (BitOrExpr "|" BitXorExpr)*
// BitXorExpr -> BitAndExpr | (BitXorExpr "^" BitAndExpr)*
// BitAndExpr -> EqualityExpr | (BitAndExpr "&" EqualityExpr)*
// EqualityExpr -> RelationalExpr | (EqualityExpr ("==" | "===" | "!=" | "!==") RelationalExpr)*
// RelationalExpr -> ShiftExpr | (RelationalExpr ("<" | ">" | "<=" | ">=" | "instanceof" | "in") ShiftExpr)*
// ShiftExpr -> AdditiveExpr | (ShiftExpr ("<<" | ">>" | ">>>") AdditiveExpr)*
// AdditiveExpr -> MultiplicativeExpr | (AdditiveExpr ("+" | "-") MultiplicativeExpr)*
// MultiplicativeExpr -> ExponentialExpr | (MultiplicativeExpr ("*" | "/" | "%") ExponentialExpr)*
// ExponentialExpr -> UnaryExpr | UpdateExpr "**" UnaryExpr
//
// ..branching into UpdateExpr
//
// UnaryExpr -> UpdateExpr | (("delete" | "void" | "+" | "-" | "~" | "!") UnaryExpr)
// UpdateExpr -> LeftHandSideExpr | (LeftHandSideExpr ("++" | "--"))
// LeftHandSideExpr -> NewExpr | CallExpr
//
// ..branching into NewExpr
//
// NewExpr -> MemberExpr | ("new" NewExpr)
// MemberExpr -> PrimaryExpr | (MemberExpr "[" Expr "]") | (MemberExpr "." Identifier) | (MemberExpr TemplateLiteral) | SuperProperty | MetaProperty | ("new" MemberExpr Arguments)
//
// ..branching into PrimaryExpr
//
// PrimaryExpr -> "this" | IdentifierReference | Literal | ArrayLiteral | ObjectLiteral | FunctionExpr | ClassExpr | GeneratorExpr | AsyncFuncExpr | RegularExpr | TemplateLiteral | ParenExpr

func (p *Parser) parseExpr() {}
