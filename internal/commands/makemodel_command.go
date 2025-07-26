package commands

import (
	"cornyk/gin-template/pkg/global"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

type columnInfo struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Extra   string
	Comment string
}

type ColumnDef struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}

func MakeModelCmd() *cobra.Command {
	var dbConn string

	cmd := &cobra.Command{
		Use:     "makemodel <model_name>",
		Short:   "自动生成包含表注释和字段注释的GORM模型",
		Example: "makemodel user_model\nmakemodel order_model --db=secondary",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			modelName := strings.TrimSuffix(args[0], "_model")
			return generateModel(modelName, dbConn)
		},
	}

	cmd.Flags().StringVarP(&dbConn, "db", "d", "default", "数据库连接名称")
	return cmd
}

func generateModel(modelName, dbConn string) error {
	db := global.DBConn(dbConn)
	if db == nil {
		return fmt.Errorf("数据库连接 '%s' 未初始化", dbConn)
	}

	// 获取表名和注释
	tableName, tableComment, err := getTableInfo(db, modelName)
	if err != nil {
		return err
	}

	// 获取字段信息
	columns, err := getTableColumns(db, tableName)
	if err != nil {
		return err
	}

	// 准备模板数据
	data := struct {
		Package      string
		ModelName    string
		TableName    string
		TableComment string
		Columns      []ColumnDef
	}{
		Package:      "models",
		ModelName:    toCamelCase(modelName),
		TableName:    tableName,
		TableComment: tableComment,
		Columns:      columns,
	}

	// 生成文件
	fileName := toSnakeCase(modelName) + "_model.go"
	filePath := filepath.Join("internal/models", fileName)
	return renderTemplate(filePath, data)
}

// 获取表信息和注释
func getTableInfo(db *gorm.DB, modelName string) (string, string, error) {
	singularTable := toSnakeCase(modelName)
	pluralTable := toPlural(singularTable)

	// 先尝试单数表名
	if hasTable(db, singularTable) {
		comment, err := getTableComment(db, singularTable)
		return singularTable, comment, err
	}

	// 再尝试复数表名
	if hasTable(db, pluralTable) {
		comment, err := getTableComment(db, pluralTable)
		return pluralTable, comment, err
	}

	return "", "", fmt.Errorf("未找到表: %s (尝试过: %s 和 %s)", modelName, singularTable, pluralTable)
}

// 获取表注释 (MySQL)
func getTableComment(db *gorm.DB, table string) (string, error) {
	var comment string
	err := db.Raw(`
		SELECT table_comment 
		FROM information_schema.tables 
		WHERE table_schema = DATABASE() 
		AND table_name = ?
	`, table).Scan(&comment).Error

	if comment == "" {
		comment = table + " table" // 默认注释
	}
	return comment, err
}

// 获取表字段信息
func getTableColumns(db *gorm.DB, table string) ([]ColumnDef, error) {
	var columns []columnInfo
	if err := db.Raw("SHOW FULL COLUMNS FROM " + table).Scan(&columns).Error; err != nil {
		return nil, err
	}

	var result []ColumnDef
	for _, c := range columns {
		result = append(result, ColumnDef{
			Name:    toCamelCase(c.Field),
			Type:    mapDBType(c.Type, c.Null == "YES"),
			Tag:     buildGormTag(c.Field, c.Key, c.Extra),
			Comment: c.Comment,
		})
	}
	return result, nil
}

// 构建GORM和JSON标签
func buildGormTag(field, key, extra string) string {
	tag := "column:" + field
	if key == "PRI" {
		tag += ";primaryKey"
	}
	if strings.Contains(extra, "auto_increment") {
		tag += ";autoIncrement"
	}
	return fmt.Sprintf("`gorm:\"%s\" json:\"%s\"`", tag, toSnakeCase(field))
}

// 数据库类型映射到Go类型
func mapDBType(dbType string, nullable bool) string {
	// 移除括号内容如varchar(255)
	cleanType := regexp.MustCompile(`\(.*\)`).ReplaceAllString(dbType, "")

	switch strings.ToLower(cleanType) {
	case "int", "integer", "tinyint", "smallint", "mediumint":
		if nullable {
			return "*int32"
		}
		return "int32"
	case "bigint":
		if nullable {
			return "*int64"
		}
		return "int64"
	case "char", "varchar", "text", "longtext":
		return "string"
	case "datetime", "timestamp", "date":
		return "time.Time"
	case "decimal", "float", "double":
		if nullable {
			return "*float64"
		}
		return "float64"
	case "json", "jsonb":
		return "datatypes.JSON"
	case "boolean", "bool":
		if nullable {
			return "*bool"
		}
		return "bool"
	default:
		return "string"
	}
}

// 渲染模板
func renderTemplate(path string, data interface{}) error {
	const tpl = `package {{.Package}}

import (
	"gorm.io/gorm"
	{{- if hasTime .Columns}}
	"time"
	{{- end}}
	{{- if hasUUID .Columns}}
	"github.com/google/uuid"
	{{- end}}
	{{- if hasJSON .Columns}}
	"gorm.io/datatypes"
	{{- end}}
)

// {{.ModelName}} {{.TableComment}}
type {{.ModelName}} struct {
	{{- range .Columns}}
	{{.Name}} {{.Type}} {{.Tag}}{{if .Comment}} // {{.Comment}}{{end}}
	{{- end}}
}

func ({{firstChar .ModelName}} *{{.ModelName}}) TableName() string {
	return "{{.TableName}}"
}
`

	funcMap := template.FuncMap{
		"hasTime": func(columns []ColumnDef) bool {
			for _, c := range columns {
				if strings.Contains(c.Type, "time.Time") {
					return true
				}
			}
			return false
		},
		"hasUUID": func(columns []ColumnDef) bool {
			for _, c := range columns {
				if strings.Contains(c.Type, "uuid.UUID") {
					return true
				}
			}
			return false
		},
		"hasJSON": func(columns []ColumnDef) bool {
			for _, c := range columns {
				if c.Type == "datatypes.JSON" {
					return true
				}
			}
			return false
		},
		"firstChar": func(s string) string {
			return strings.ToLower(string(s[0]))
		},
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl := template.Must(template.New("model").Funcs(funcMap).Parse(tpl))
	return tmpl.Execute(file, data)
}

// --- 辅助函数 ---
func toCamelCase(s string) string {
	words := strings.Split(strings.ReplaceAll(s, "_", " "), " ")
	for i := range words {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}

func toSnakeCase(s string) string {
	var re = regexp.MustCompile("([a-z0-9])([A-Z])")
	return strings.ToLower(re.ReplaceAllString(s, "${1}_${2}"))
}

func toPlural(s string) string {
	if strings.HasSuffix(s, "y") {
		return strings.TrimSuffix(s, "y") + "ies"
	}
	return s + "s"
}

func hasTable(db *gorm.DB, table string) bool {
	var exists bool
	db.Raw("SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name = ?)", table).Scan(&exists)
	return exists
}
