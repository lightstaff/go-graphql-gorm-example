package scalar

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"

	"github.com/lightstaff/go-graphql-gorm-example/dao/types"
)

// NullStringの変換定義
var NullStringScalar = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "NullString",
	Description: "Support for null string",
	Serialize: func(value interface{}) interface{} {
		switch value := value.(type) {
		case types.NullString:
			return value.String
		case *types.NullString:
			return value.String
		default:
			return nil
		}
	},
	ParseValue: func(value interface{}) interface{} {
		switch value := value.(type) {
		case string:
			return types.NewNullString(value)
		case *string:
			return types.NewNullString(*value)
		default:
			return nil
		}
	},
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			return types.NewNullString(valueAST.Value)
		default:
			return nil
		}
	},
})
