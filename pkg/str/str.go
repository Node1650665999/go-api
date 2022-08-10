package str

import (
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

// Plural 单数转复数： user -> users
func Plural(word string) string {
	return pluralize.NewClient().Plural(word)
}

// Singular 复数转单数： users -> user
func Singular(word string) string {
	return pluralize.NewClient().Singular(word)
}

// Snake 驼峰转下划线： TopicComment -> topic_comment
func Snake(s string) string {
	return strcase.ToSnake(s)
}

// Camel 下划线转驼峰： topic_comment -> TopicComment
func Camel(s string) string {
	return strcase.ToCamel(s)
}

// LowerCamel 驼峰首字母转小写：TopicComment -> topicComment
func LowerCamel(s string) string {
	return strcase.ToLowerCamel(s)
}
