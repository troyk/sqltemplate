package sqltemplate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

var reNameHolder = regexp.MustCompile(":[A-Za-z0-9_-]+")

type SqlTemplate struct {
	Templates *template.Template
}

func (t *SqlTemplate) Add(name string, sql string) *template.Template {
	return template.Must(t.Templates.New(name).Parse(sql))
}

// Executes a template and returns sql and arguments suitable for database/sql usage
// data is exctracted using JSON and will convert :name to $1,$2,holders
func (t *SqlTemplate) ToSql(name string, data interface{}) (sql string, args []interface{}, err error) {
	buf := &bytes.Buffer{}
	err = t.Templates.ExecuteTemplate(buf, name, data)
	sql = buf.String()
	if err != nil {
		return sql, nil, err
	}
	var dataMap map[string]interface{}
	var dataJson []byte
	dataJson, err = json.Marshal(data)
	if err != nil {
		return sql, nil, err
	}
	if err = json.Unmarshal(dataJson, &dataMap); err != nil {
		return sql, nil, err
	}

	nameHolders := reNameHolder.FindAllString(sql, -1)
	for _, name := range nameHolders {
		value, exists := dataMap[name[1:]]
		if !exists {
			continue
		}
		args = append(args, value)
		dollar := fmt.Sprintf("$%d", len(args))
		sql = strings.Replace(sql, name, dollar, -1)
	}
	return sql, args, err
}

func NewSqlTemplate() *SqlTemplate {
	t := &SqlTemplate{Templates: template.New("sql")}
	funcMap := template.FuncMap{
		"sqlBuilder": func() *sqlBuilder {
			return &sqlBuilder{}
		},
	}

	t.Templates.Funcs(funcMap)
	return t
}
