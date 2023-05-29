
func (m *default{{.upperStartCamelObject}}Model) FindOneBy{{.upperField}}({{.in}}) (result *{{.upperStartCamelObject}}, err error) {
    err = m.Conn.Table(m.table).Where("{{.originalField}}", {{.lowerStartCamelField}}).Take(&result).Error
    return
}
