
func New{{.upperStartCamelObject}}Model(conn *gorm.DB) {{.upperStartCamelObject}}Model {
	return &default{{.upperStartCamelObject}}Model{
		Conn:conn,
		table:      {{.table}},
	}
}
