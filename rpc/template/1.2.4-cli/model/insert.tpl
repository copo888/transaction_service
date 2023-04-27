
func (m *default{{.upperStartCamelObject}}Model) Insert(data *{{.upperStartCamelObject}}) (err error) {
	return m.Conn.Table(m.table).Create(&data).Error
}
