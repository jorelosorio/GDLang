%{
/*
 * Copyright (C) 2023 The GDLang Team.
 *
 * This file is part of GDLang.
 *
 * GDLang is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * GDLang is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with GDLang.  If not, see <http://www.gnu.org/licenses/>.
 */

package ast

import (
       "gdlang/lib/runtime"
)
%}

%union {
       token                *NodeTokenInfo
       
       node                 Node
       node_list            []Node

       flag                 bool

       gd_type              runtime.GDTypable
       gd_type_list         []runtime.GDTypable
}

%token  <token>                    LAS LLSHIFT LRSHIFT LIDENT LINT LFLOAT LSTRING LIMAG LCHAR LCOMMENT LMUL
%token  <token>                    LADD LSUB LQUO LREM
%token  <token>                    LQMARK LNSAFE LADD_ASSIGN LSUB_ASSIGN LMUL_ASSIGN LQUO_ASSIGN LREM_ASSIGN
%token  <token>                    LARROW LINC LDEC
%token  <token>                    LLAND LOR LLOR LNOT
%token  <token>                    LEQL LLSS LGTR LASSIGN
%token  <token>                    LNEQ LLEQ LGEQ LELLIPSIS
%token  <token>                    LLPAREN LLBRACK LLBRACE LCOMMA LPERIOD LRPAREN LRBRACK LRBRACE LSEMICOLON LCOLON LCOLONCOLON
%token  <token>                    LUSE LTYPEALIAS LSET LPUB LCONST LELSE LFOR LIN LFUNC LIF LBREAK LRETURN
%token  <token>                    LTANY LTBOOL LTINT LTFLOAT LTCOMPLEX LTSTRING LTCHAR
%token  <token>                    LTRUE LFALSE LNIL

%type   <node>                     file_body_stmt break_stmt return_stmt stmt expr pseudocall uexpr pexpr 
%type   <node>                     set mut_collection_op literal update_obj block block_stmt func lambda tuple array
%type   <node_list>                optional_expr_list optional_file_body_stmt_list file_body_stmt_list expr_list tuple_expr_list optional_block_stmt_list block_stmt_list

%type   <node>                     struct struct_attr for_if_stmt for_in_stmt if_expr if_stmt elseif_stmt else_stmt selexpr ident file use ident_with_type ident_with_optional_type optional_assign_expr const_ident_with_optional_type
%type   <node_list>                struct_attr_list elseif_stmt_list optional_file_package_list use_list ident_access_list ident_list func_arg_list optional_func_arg_list set_expr_list const_ident_with_optional_type_list set_expr_option_list
%type   <node>                     typealias cast_expr

%type   <flag>                     safe_accessor optional_const optional_pub optional_trailing_comma

%type   <gd_type>                  optional_return_type func_type type union_type tuple_type unary_type array_type struct_type struct_attr_type obj_optional_type
%type   <gd_type_list>             struct_attr_type_list type_list tuple_attr_type_list

%error LSET LIDENT LCOLON LNIL:
       "NIL_AS_A_TYPE_ERR"

%error optional_file_package_list optional_file_body_stmt_list LUSE:
       "USE_ONLY_AT_HEADER_ERR"

%left  LAS
%left  LLSHIFT LRSHIFT
%left  LQMARK LCOLON
%left  LLOR
%left  LLAND
%left  LEQL LNEQ LLSS LGTR LLEQ LGEQ
%left  LADD LSUB
%left  LMUL LQUO LREM

%left  LLPAREN
%left  LRPAREN

%%

file:
       optional_file_package_list
       optional_file_body_stmt_list {
              file := NewNodeFile($1, $2)
              yylex.(*Ast).Root = file
              $$ = file
       }
;

// Packages

optional_file_package_list:
       use_list LSEMICOLON {
              $$ = $1
       }
       | /* empty */ {
              $$ = make([]Node, 0)
       }
;

use_list:
       use_list LSEMICOLON use {
              $1 = append($1, $3)
              $$ = $1
       }
       | use {
              $$ = make([]Node, 1)
              $$[0] = $1
       }
;

use:
       LUSE ident_access_list LLBRACE ident_list LRBRACE {
              $$ = NewNodePackage($2, $4)
       }
;

ident_list:
       ident LCOMMA ident_list {
              $3 = append($3, $1)
              $$ = $3
       }
       | ident {
              $$ = []Node{$1}
       }
;

// ident., ...
ident_access_list:
       ident_access_list LPERIOD ident {
              $1 = append($1, $3)
              $$ = $1
       }
       | ident {
              $$ = []Node{$1}
       }
;

optional_file_body_stmt_list:
       file_body_stmt_list LSEMICOLON
       | /* empty */ {
              $$ = make([]Node, 0)
       }
;

file_body_stmt_list:
       file_body_stmt_list LSEMICOLON file_body_stmt {
              $1 = append($1, $3)
              $$ = $1
       }
       | file_body_stmt {
              $$ = make([]Node, 1)
              $$[0] = $1
       }
;

file_body_stmt:
       optional_pub set {
              sets, ok := $2.(*NodeSets)
              if !ok {
                     panic("file_body_stmt: Invalid `*NodeSets` object")
              }

              for _, node := range sets.Nodes {
                     if set, ok := node.(*NodeSet); ok {
                            set.IsPub = $1
                     } else {
                            panic("file_body_stmt: Invalid `*NodeSet` object")
                     }
              }

              $$ = $2
       }
       | optional_pub func {
              $2.(*NodeFunc).IsPub = $1
              $$ = $2
       }
       | optional_pub typealias {
              $2.(*NodeTypeAlias).IsPub = $1
              $$ = $2
       }
;

// Public

optional_pub:
       LPUB {
              $$ = true
       }
       | /* empty */ {
              $$ = false
       }
;

// Type alias

typealias:
       LTYPEALIAS ident LASSIGN type {
              $$ = NewNodeTypeAlias(false, $2.(*NodeIdent), $4)
       }
;

// Statements

stmt:
       set
       | update_obj
       | mut_collection_op
       | func
       | for_in_stmt
       | for_if_stmt
       | lambda
       | typealias
       | pseudocall
       | if_stmt
;

// Optional comma for trailing comma
optional_trailing_comma:
       LCOMMA {
              $$ = true
       }
       | /* empty */ {
              $$ = false
       }
;

// Set objs

set:
       LSET set_expr_list {
              $$ = NewNodeSets($2)
       }
;

// e.g. a: int = 0, b: float = 1.0
set_expr_list:
       set_expr_list LCOMMA set_expr_option_list {
              $1 = append($1, $3...)
              $$ = $1
       }
       | set_expr_option_list
;

set_expr_option_list:
       const_ident_with_optional_type optional_assign_expr /*(= expr)?*/ {
              nodeSet, ok := $1.(*NodeSet)
              if !ok {
                     panic("set_expr_option_list: Invalid `*NodeSet` object")
              }
              nodeSet.Expr = $2
              $$ = []Node{nodeSet}
       }
       // Deconstructing set object declaration
       // e.g. set (a, b) = (1, 2)
       | optional_const /*const*/ LLPAREN const_ident_with_optional_type_list LRPAREN optional_assign_expr /*(= expr)?*/ {
              sharedExpr := NewNodeSharedExpr($5)
              for i, node := range $3 {
                     if set, ok := node.(*NodeSet); ok {
                            set.Index = byte(i)
                            set.Expr = sharedExpr
                     } else {
                            panic("set: Invalid `*NodeSet` object")
                     }
              }
              $$ = $3
       }
;

// Inline set object declaration for spread assignment
const_ident_with_optional_type_list:
       const_ident_with_optional_type_list LCOMMA const_ident_with_optional_type {
              $1 = append($1, $3)
              $$ = $1
       }
       | const_ident_with_optional_type {
              $$ = make([]Node, 1)
              $$[0] = $1
       }
;

const_ident_with_optional_type:
       optional_const /*const*/ ident_with_optional_type /*ident(:type)?*/ {
              identWithType, ok := $2.(*NodeIdentWithType)
              if !ok {
                     panic("const_ident_with_optional_type: Invalid `*NodeIdentWithType` object")
              }
              $$ = NewNodeSet(false, $1, identWithType, nil)
       }
;

optional_assign_expr:
       LASSIGN expr {
              $$ = $2
       }
       | /* empty */ {
              $$ = nil
       }
;

optional_const:
       LCONST {
              $$ = true
       }
       | /* empty */ {
              $$ = false
       }
;

update_obj:
       expr LASSIGN expr {
              $$ = NewNodeUpdateSet($1, $3)
       }
       | expr LADD_ASSIGN expr {
              $$ = NewNodeUpdateSet($1, NewNodeExprOperation(runtime.ExprOperationAdd, $1, $3))
       }
       | expr LSUB_ASSIGN expr {
              $$ = NewNodeUpdateSet($1, NewNodeExprOperation(runtime.ExprOperationSubtract, $1, $3))
       }
       | expr LMUL_ASSIGN expr {
              $$ = NewNodeUpdateSet($1, NewNodeExprOperation(runtime.ExprOperationMultiply, $1, $3))
       }
       | expr LQUO_ASSIGN expr {
              $$ = NewNodeUpdateSet($1, NewNodeExprOperation(runtime.ExprOperationQuo, $1, $3))
       }
       | expr LREM_ASSIGN expr {
              $$ = NewNodeUpdateSet($1, NewNodeExprOperation(runtime.ExprOperationRem, $1, $3))
       }
;

// Type with optional type
ident_with_optional_type:
       ident obj_optional_type {
              $$ = NewNodeIdentWithType($1.(*NodeIdent), $2)
       }
;

obj_optional_type:
       LCOLON type {
              $$ = $2
       }
       | /* empty */ {
              $$ = runtime.GDUntypedType
       }
;

// ident with type `ident: type`

// Typed identifier with mandatory type
// it is strongly typed, used commonly in function arguments
ident_with_type:
       ident LCOLON type {
              $$ = NewNodeIdentWithType($1.(*NodeIdent), $3)
       }
;

// Types

unary_type:
       LTINT                { $$ = runtime.GDIntType                                              }
       | LTFLOAT            { $$ = runtime.GDFloatType                                            }
       | LTCOMPLEX          { $$ = runtime.GDComplexType                                          }
       | LTBOOL             { $$ = runtime.GDBoolType                                             }
       | LTANY              { $$ = runtime.GDAnyType                                              }
       | LTSTRING           { $$ = runtime.GDStringType                                           }
       | LTCHAR             { $$ = runtime.GDCharType                                             }
       | LIDENT             { $$ = runtime.NewGDIdentRefType(runtime.NewGDStringIdent($1.Lit))       }
       | tuple_type         { $$ = $1                                                             }
       | array_type         { $$ = $1                                                             }
       | struct_type        { $$ = $1                                                             }
       | LFUNC func_type    { $$ = $2                                                             }
;

union_type:
       type LOR type {
              if cT, isCT := $1.(runtime.GDUnionType); isCT {
                     $$ = runtime.NewGDUnionType(append(cT, $3)...)
              } else {
                     $$ = runtime.NewGDUnionType($1, $3)
              }
       }
;

type:
       unary_type
       | LLPAREN union_type LRPAREN {
              $$ = $2
       }
;

tuple_type:
       LLPAREN tuple_attr_type_list LRPAREN {
              $$ = runtime.NewGDTupleType($2...)
       }
;

tuple_attr_type_list:
       // Empty tuple ()
       LCOMMA {
              $$ = make([]runtime.GDTypable, 0)
       }
       // One element tuple (n,)
       | type LCOMMA {
              $$ = make([]runtime.GDTypable, 1)
              $$[0] = $1
       }
       // Tuple with multiple elements (n, NewNode, ...)
       | type LCOMMA type_list {
              $3 = append([]runtime.GDTypable{$1}, $3...)
              $$ = $3
       }
;

array_type:
       LLBRACK type LRBRACK {
              $$ = runtime.NewGDArrayType($2)
       }
;

struct_type:
       LLBRACE struct_attr_type_list optional_trailing_comma LRBRACE {
              attrTypes := make([]runtime.GDStructAttrType, len($2))
              for i, attr := range $2 {
                     attrTypes[i] = attr.(runtime.GDStructAttrType)
              }
              $$ = runtime.NewGDStructType(attrTypes...)
       }
;

// Function type

func_arg_list:
       func_arg_list LCOMMA ident_with_type {
              $1 = append($1, $3)
              $$ = $1
       }
       | ident_with_type {
              $$ = make([]Node, 1)
              $$[0] = $1
       }
;

optional_func_arg_list:
       func_arg_list
       | /* empty */ {
              $$ = make([]Node, 0)
       }
;

optional_return_type:
       LARROW type {
              $$ = $2
       }
       | /* empty */ {
              $$ = runtime.GDNilType
       }
;

func_type:
       // () => type?
       LLPAREN optional_func_arg_list LRPAREN optional_return_type {
              $$ = buildFuncType($2, false, $4)
       }
       // (arg, ...?) => type?
       | LLPAREN func_arg_list LCOMMA LELLIPSIS LRPAREN optional_return_type {
              $$ = buildFuncType($2, true, $6)
       }
;

struct_attr_type_list:
       struct_attr_type_list LCOMMA struct_attr_type {
              $1 = append($1, $3)
              $$ = $1
       }
       | struct_attr_type {
              $$ = make([]runtime.GDTypable, 1)
              $$[0] = $1
       }
;

struct_attr_type:
       ident LCOLON type {
              ident := runtime.NewGDStringIdent($1.(*NodeIdent).Lit)
              $$ = runtime.GDStructAttrType{Ident: ident, Type: $3}
       }
;

type_list:
       type_list LCOMMA type {
              $1 = append($1, $3)
              $$ = $1
       }
       | type {
              $$ = make([]runtime.GDTypable, 1)
              $$[0] = $1
       }
;

// Blocks

block:
       LLBRACE optional_block_stmt_list LRBRACE {
              $$ = NewNodeBlock($2)
       }
;

block_stmt:
       stmt
       | return_stmt
       | break_stmt
;

return_stmt:
       LRETURN              { $$ = NewNodeReturn($1, nil)  }
       | LRETURN expr       { $$ = NewNodeReturn($1, $2)   }
;

break_stmt:
       LBREAK {
              $$ = NewNodeBreak($1)
       }
;

optional_block_stmt_list:
       block_stmt_list LSEMICOLON
       | /* empty */ {
              $$ = make([]Node, 0)
       }
;

// block with return stmt list
block_stmt_list:
       block_stmt_list LSEMICOLON block_stmt {
              $1 = append($1, $3)
              $$ = $1
       }
       | block_stmt {
              $$ = make([]Node, 1)
              $$[0] = $1
       }
;

// Lambda

lambda:
       LFUNC func_type block {
              $$ = NewNodeLambda($2.(*runtime.GDLambdaType), $3.(*NodeBlock))
       }
;

// Function

func:
       LFUNC ident func_type block {
              $$ = NewNodeFunc(true, $2.(*NodeIdent), $3.(*runtime.GDLambdaType), $4.(*NodeBlock))
       }
;

// Expression

if_expr:
       expr LQMARK expr LCOLON expr { // cond ? expr : expr
              $$ = NewNodeTernaryIf($1, $3, $5)
       }
;

cast_expr:
       expr LAS type {
              $$ = NewNodeCastExpr($1, $3)
       }
;

expr:
       uexpr
       | if_expr          // Ternary if (cond ? expr : expr)
       | cast_expr        // Type cast (expr as type)
       | mut_collection_op // Add or remove from a collection (<< | >>)
       | expr LLOR expr { // ||
              $$ = NewNodeExprOperation(runtime.ExprOperationOr, $1, $3)
       }
       | expr LLAND expr { // &&
              $$ = NewNodeExprOperation(runtime.ExprOperationAnd, $1, $3)
       }
       | expr LEQL expr { // ==
              $$ = NewNodeExprOperation(runtime.ExprOperationEqual, $1, $3)
       }
       | expr LNEQ expr { // !=
              $$ = NewNodeExprOperation(runtime.ExprOperationNotEqual, $1, $3)
       }
       | expr LLSS expr { // <
              $$ = NewNodeExprOperation(runtime.ExprOperationLess, $1, $3)
       }
       | expr LGTR expr { // >
              $$ = NewNodeExprOperation(runtime.ExprOperationGreater, $1, $3)
       }
       | expr LLEQ expr { // <=
              $$ = NewNodeExprOperation(runtime.ExprOperationLessEqual, $1, $3)
       }
       | expr LGEQ expr { // >=
              $$ = NewNodeExprOperation(runtime.ExprOperationGreaterEqual, $1, $3)
       }
       | expr LADD expr { // +
              $$ = NewNodeExprOperation(runtime.ExprOperationAdd, $1, $3)
       }
       | expr LSUB expr { // -
              $$ = NewNodeExprOperation(runtime.ExprOperationSubtract, $1, $3)
       }
       | expr LMUL expr { // *
              $$ = NewNodeExprOperation(runtime.ExprOperationMultiply, $1, $3)
       }
       | expr LQUO expr { // /
              $$ = NewNodeExprOperation(runtime.ExprOperationQuo, $1, $3)
       }
       | expr LREM expr { // %
              $$ = NewNodeExprOperation(runtime.ExprOperationRem, $1, $3)
       }
;

uexpr:
       pexpr
       | LADD uexpr  { $$ = NewNodeExprOperation(runtime.ExprOperationUnaryPlus, $2, nil)       }
       | LSUB uexpr  { $$ = NewNodeExprOperation(runtime.ExprOperationUnaryMinus, $2, nil)  }
       | LNOT uexpr  { $$ = NewNodeExprOperation(runtime.ExprOperationNot, $2, nil)       }
;

pexpr:
       selexpr
       | LLPAREN expr LRPAREN {
              $$ = $2
       }
;

mut_collection_op:
       expr LLSHIFT expr {
              $$ = NewNodeMutCollectionOp(MutableCollectionAddOp, $1, $3)
       }
       | expr LRSHIFT expr {
              $$ = NewNodeMutCollectionOp(MutableCollectionRemoveOp, $1, $3)
       }
;

selexpr:
       literal
       | ident
       | tuple
       | array
       | struct
       | lambda
       | pexpr LELLIPSIS {
              $$ = NewNodeEllipsisExpr($1)
       }
       // Attribute accessor expression with identifier
       // e.g. obj.attr
       | pexpr safe_accessor ident {
              $$ = NewNodeSafeDotExpr($1, $2, $3)
       }
       // Array accessor expression
       | pexpr LLBRACK expr LRBRACK {
              $$ = NewNodeIterIdxExpr(false, $1, $3)
       }
       // Pseudocall expression
       | pseudocall
;

pseudocall:
       pexpr LLPAREN optional_expr_list LRPAREN {
              $$ = NewNodeCallExpr($1, $3)
       }
;

safe_accessor:
       LPERIOD       { $$ = false  }
       | LNSAFE      { $$ = true   }
;

// Expression list comma separated
expr_list:
       expr_list LCOMMA expr {
              $1 = append($1, $3)
              $$ = $1
       }
       | expr {
              $$ = make([]Node, 1)
              $$[0] = $1
       }
;

optional_expr_list:
       expr_list optional_trailing_comma
       | /* empty */ {
              $$ = make([]Node, 0)
       }
;

// Literals

literal:
       LINT          { $$ = NewNodeLiteral($1)  }
       | LFLOAT      { $$ = NewNodeLiteral($1)  }
       | LSTRING     { $$ = NewNodeLiteral($1)  }
       | LNIL        { $$ = NewNodeLiteral($1)  }
       | LTRUE       { $$ = NewNodeLiteral($1)  }
       | LFALSE      { $$ = NewNodeLiteral($1)  }
       | LIMAG       { $$ = NewNodeLiteral($1)  }
       | LCHAR       { $$ = NewNodeLiteral($1)  }
;

ident:
       LIDENT        { $$ = NewNodeIdent($1)    }
;

// Tuple

tuple:
       LLPAREN tuple_expr_list LRPAREN {
              $$ = NewNodeTuple($2...)
       }
;

// Struct

struct:
       LLBRACE struct_attr_list optional_trailing_comma LRBRACE {
              $$ = NewNodeStruct($2...)
       }
       | LLBRACE LRBRACE {
              $$ = NewNodeStruct()
       }
;

struct_attr_list:
       struct_attr_list LCOMMA struct_attr {
              $1 = append($1, $3)
              $$ = $1
       }
       | struct_attr {
              $$ = make([]Node, 1)
              $$[0] = $1
       }
;

struct_attr:
       ident LCOLON expr {
              $$ = NewNodeStructAttr($1.(*NodeIdent), $3)
       }
;

tuple_expr_list:
       // Empty tuple ()
       LCOMMA {
              $$ = make([]Node, 0)
       }
       // One element tuple (n,)
       | expr LCOMMA {
              $$ = make([]Node, 1)
              $$[0] = $1
       }
       // Tuple with multiple elements (n, NewNode, ...)
       | expr LCOMMA expr_list {
              $3 = append([]Node{$1}, $3...)
              $$ = $3
       }
;

// Array

array:
       LLBRACK optional_expr_list LRBRACK {
              $$ = NewNodeArray($1, $3, $2)
       }
;

// Statements

for_in_stmt:
       LFOR set LIN expr block {
              $$ = NewNodeForIn($2, $4, $5.(*NodeBlock))
       }
;

for_if_stmt:
       LFOR set LIF expr_list block {
              $$ = NewNodeForIf($2, $4, $5.(*NodeBlock))
       }
       | LFOR LIF expr_list block {
              $$ = NewNodeForIf(nil, $3, $4.(*NodeBlock))
       }
       | LFOR block {
              $$ = NewNodeForIf(nil, nil, $2.(*NodeBlock))
       }
;

if_stmt:
       LIF expr_list block elseif_stmt_list else_stmt {
              nIf := NewNodeIf($2, $3.(*NodeBlock))
              $$ = NewNodeIfElse(nIf, $4, $5)
       }
;

elseif_stmt_list:
       elseif_stmt_list elseif_stmt {
              $1 = append($1, $2)
              $$ = $1
       }
       | /* empty */ {
              $$ = make([]Node, 0)
       }
;

elseif_stmt:
       LELSE LIF expr_list block {
              $$ = NewNodeIf($3, $4.(*NodeBlock))
       }
;

else_stmt:
       LELSE block {
              $$ = NewNodeIf(nil, $2.(*NodeBlock))
       }
       | /* empty */ {
              $$ = nil
       }
;

%%