
func (m *default{{.upperStartCamelObject}}Model) Delete(req interface{}) error {
    return m.Conn.Table(m.table).Delete(&req).Error
}
