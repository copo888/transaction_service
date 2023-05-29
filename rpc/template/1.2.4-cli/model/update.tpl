
func (m *default{{.upperStartCamelObject}}Model) Update(data *{{.upperStartCamelObject}}) (err error) {
	return m.Conn.Table(m.table).Updates(data).Error
}
