package sqltemplate

import (
	"fmt"
	"reflect"
	"strconv"
)

type sqlBuilder struct {
	where bool
}

func (builder *sqlBuilder) If(condition interface{}, clause string) string {
	if !reflect.DeepEqual(condition, reflect.Zero(reflect.TypeOf(condition)).Interface()) {
		return builder.And(clause)
	}
	return ""
}

func (builder *sqlBuilder) Where(clause string) string {
	if builder.where {
		return (clause)
	}
	builder.where = true
	return fmt.Sprintf(" WHERE %s ", clause)
}

func (builder *sqlBuilder) And(clause string) string {
	if builder.where {
		return builder.Where(fmt.Sprintf(" AND %s ", clause))
	}
	return builder.Where(clause)

}

func (builder *sqlBuilder) Limit(anynum interface{}) string {
	//switch limit.(type) {
	//  case int, int8, int16, int32, int64,uint, uint8, uint16, uint32, uint64:

	//}
	numstr := fmt.Sprintf("%v", anynum)
	num, _ := strconv.Atoi(numstr)
	if num > 0 {
		return fmt.Sprintf(" LIMIT %s ", numstr)
	}
	return ""
}
