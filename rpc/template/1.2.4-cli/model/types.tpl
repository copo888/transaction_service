
type (
	{{.upperStartCamelObject}}Model interface{
		{{.method}}
	}

	default{{.upperStartCamelObject}}Model struct {
		Conn *gorm.DB
		table string
	}

	{{.upperStartCamelObject}} struct {
		{{.fields}}
	}
)
