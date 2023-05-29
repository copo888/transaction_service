import (
	"database/sql"
	"strings"
	{{if .time}}"time"{{end}}

	"gorm.io/gorm"
)
