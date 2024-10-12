package ast

import "gdlang/src/gd/scanner"

type NodeType struct {
	TypeTokenInfo *NodeTokenInfo
	*BaseNode
}

func (t *NodeType) GetPosition() scanner.Position {
	return t.TypeTokenInfo.GetPosition()
}

func NewNodeType(typeTokenInfo *NodeTokenInfo) *NodeType {
	return &NodeType{typeTokenInfo, &BaseNode{}}
}
