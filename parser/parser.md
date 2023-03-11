# Parsing

Given a myriad of possible parsing techniques for parsing ECMAScript, we'll go over the industry-proven approach [adopted by V8](https://chromium.googlesource.com/v8/v8/+/2893b9fbd61dc7e50e679d21e7850e8486d4320d/src/parsing/preparser.h#19): recursive descent.

The ECMA standard defines its syntactic grammar entirely through [left recursion productions](https://en.wikipedia.org/wiki/Left_recursion). Specifically, **direct** left-recursion, which means that the productions have the form:

$$$
A -> A \alpha | \beta
$$$

in which $\alpha$ represents whatever comes after the non-terminal $A$, and $\beta$ represents the other alternative to $A$.

This rule produces the set of strings

$$$
\beta
\beta\alpha
\beta\alpha\alpha
\beta\alpha\alpha\alpha
...
$$$

It is not as efficient to parse a left-recursive grammar with a recursive descent parser because we'd have to descend into the left recursion _for every token_, as illustrated in the example above. Given the size of the ECMAScript grammar, this would be _very_ inefficient. Instead, we'll use a technique called [left factoring](https://en.wikipedia.org/wiki/Left_factoring)[^1] to remove the left recursion.

In our previous example, we could derive non-recursive rules that would produce the same strings:

$$$
A -> \beta A'
A' -> \alpha A'
$$$

Which could be further simplified into a one-step transformation:

So this is the transformation technique we're going to apply to remove the left recursion from the ECMAScript grammar.

## The ECMAScript Grammar

The following rule

$$$
Expr -> Expr "," AssignExpr | AssignExpr
$$$

Becomes

$$$
Expr -> AssignExpr Expr'
Expr' -> "," AssignExpr Expr'
$$$


### Expressions

Epsilon expressions are removed for brevity.

```hs
--
-- Expr
--
Expr  = AssignExpr | Expr'
Expr' = "," AssignExpr Expr' | Epsilon

-- AssignExpr
AssignExpr = CondExpr | YieldExpr | ArrowFunc | LeftHandSideExpr | LeftHandSideExpr "=" AssignExpr | LeftHandSideExpr ("=" | "+=" | "*=" | "/=" | "|=" | "&=" | "&&=") AssignExpr

-- ..branching off CondExpr ::
-- CondExpr
--
CondExpr  = LogOrExpr | LogOrExpr "?" Expr ":" AssignExpr
LogOrExpr = LogOrExpr "||" LogAndExpr | LogAndExpr

-- LogOrExpr
LogOrExpr  = LogAndExpr LogOrExpr'
LogOrExpr' = "||" LogAndExpr LogOrExpr'

-- LogAndExpr
LogAndExpr  = BitOrExpr LogAndExpr'
LogAndExpr' = "&&" BitOrExpr LogAndExpr'

-- BitOrExpr
BitOrExpr  = BitXorExpr BitOrExpr'
BitOrExpr' = "|" BitXorExpr BitOrExpr'

-- BitXorExpr
BitXorExpr  = BitAndExpr BitXorExpr'
BitXorExpr' = "^" BitAndExpr BitXorExpr'

-- BitAndExpr
BitAndExpr = EqualityExpr BitAndExpr'
BitAndExpr = "&" EqualityExpr BitAndExpr'

-- EqualityExpr
EqualityExpr  = RelationalExpr EqualityExpr'
EqualityExpr' = ("==" | "!=" | "===" | "!==") RelationalExpr EqualityExpr'

-- RelationalExpr
RelationalExpr  = ShiftExpr RelationalExpr'
RelationalExpr' = ("<" | ">" | "<=" | ">=" | "instanceof" | "in") | ShiftExpr RelationalExpr'

-- ShiftExpr
ShiftExpr  = AdditiveExpr ShiftExpr'
ShiftExpr' = ("<<" | ">>" | ">>>") AdditiveExpr ShiftExpr'

-- AdditiveExpr
AdditiveExpr = MultiplicativeExpr AdditiveExpr'
AdditiveExpr' = ("+" | "-") MultiplicativeExpr AdditiveExpr'

-- MultiplicativeExpr 
MultiplicativeExpr = ExponentialExpr MultiplicativeExpr'
MultiplicativeExpr' = ("*" | "/", "%") ExponentialExpr MultiplicativeExpr'

-- ExponentialExpr
ExponentialExpr = UnaryExpr | UpdateExpr "**" UnaryExpr

-- UnaryExpr
UnaryExpr = UpdateExpr | ("delete" | "void" | "typeof" | "+" | "-" | "~" | "!") UnaryExpr

-- UpdateExpr
UpdateExpr = LeftHandSideExpr | LeftHandSideExpr ("++" | "--") | ("++" | "--") LeftHandSideExpr

-- LeftHandSideExpr
LeftHandSideExpr = NewExpr | CallExpr

-- ..branching into NewExpr ::
-- NewExpr
NewExpr = MemberExpr | "new" NewExpr

-- MemberExpr
MemberExpr = PrimaryExpr | MemberExpr "[" Expr "]" | MemberExpr "." IdentifierName | MemberExpr TemplateLiteral | SuperProperty | SuperCall

-- ..branching into PrimaryExpr ::
-- PrimaryExpr
PrimaryExpr = "this" | Identifier | Literal | ArrayLiteral | ObjectLiteral | FunctionExpr | ClassExpr | GeneratorExpr | RegularExpr | TemplateLiteral | "(" Expr ")"

-- and we are back to expanding Expr again
--
-- what is left to expand:
-- YieldExpr
-- ArrowFunc
-- CallExpr
-- SuperProperty
-- SuperCall
-- FunctionExpr
-- ClassExpr
-- GeneratorExpr
-- RegularExpr
```