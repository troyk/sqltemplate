package sqltemplate

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTemplate(t *testing.T) {

	templates := NewSqlTemplate()

	templates.Add("test", `
		{{ $sql := sqlBuilder }}
		{{ $sql.If .foo "foo = :foo" }}}
		{{ $sql.And "(upper(cow)) = upper(:cow)" }}

	`)

	Convey("Render", t, func() {

		data := map[string]interface{}{
			"foo": "bar",
			"cow": "bell",
		}

		sql, args, err := templates.ToSql("test", data)

		So(err, ShouldBeNil)
		So(args, ShouldResemble, []interface{}{"bar", "bell"})
		So(sql, ShouldContainSubstring, "foo = $1")
		So(sql, ShouldContainSubstring, "upper($2)")

	})

}
