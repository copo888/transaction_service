
func (m *default{{.upperStartCamelObject}}Model) FindOne({{.lowerStartCamelPrimaryKey}} {{.dataType}}) (result *{{.upperStartCamelObject}},err error) {
	 err = m.Conn.Table(m.table).Take(&result, {{.lowerStartCamelPrimaryKey}}).Error
	 return
}
