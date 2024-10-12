package runtime

type GDTypeAliasType struct {
	GDIdent
	GDTypable
}

func (t *GDTypeAliasType) GetCode() GDTypableCode { return GDTypeAliasTypeCode }
func (t *GDTypeAliasType) ToString() string       { return t.GDIdent.ToString() }

func NewGDTypeAliasType(ident GDIdent, typ GDTypable) *GDTypeAliasType {
	return &GDTypeAliasType{ident, typ}
}
func NewGDStrTypeAliasType(ident string, typ GDTypable) *GDTypeAliasType {
	return NewGDTypeAliasType(NewGDStrIdent(ident), typ)
}
