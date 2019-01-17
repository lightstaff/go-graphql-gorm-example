package main

import (
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/lightstaff/go-graphql-gorm-example/dao"
	"github.com/lightstaff/go-graphql-gorm-example/scalar"
)

// GraphQL上のユーザー定義
var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if data, ok := p.Source.(dao.User); ok {
					return data.ID, nil
				}

				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if data, ok := p.Source.(dao.User); ok {
					return data.Name, nil
				}

				return nil, nil
			},
		},
	},
})

// GraphQL上のメール定義
var emailType = graphql.NewObject(graphql.ObjectConfig{
	Name: "EMail",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if data, ok := p.Source.(dao.EMail); ok {
					return data.ID, nil
				}

				return nil, nil
			},
		},
		"address": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if data, ok := p.Source.(dao.EMail); ok {
					return data.Address, nil
				}

				return nil, nil
			},
		},
		"remarks": &graphql.Field{
			Type: scalar.NullStringScalar,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if data, ok := p.Source.(dao.EMail); ok {
					return data.Remarks, nil
				}

				return nil, nil
			},
		},
		"userId": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if data, ok := p.Source.(dao.EMail); ok {
					return data.UserID, nil
				}

				return nil, nil
			},
		},
	},
})

// GraphQL循環参照エラー対策
func init() {
	userType.AddFieldConfig("emails", &graphql.Field{
		Type: graphql.NewList(emailType),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if data, ok := p.Source.(dao.User); ok {
				return data.EMails, nil
			}

			return nil, nil
		},
	})

	emailType.AddFieldConfig("user", &graphql.Field{
		Type: userType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if data, ok := p.Source.(dao.EMail); ok {
				if data.User != nil {
					return *data.User, nil
				}

				return nil, nil
			}

			return nil, nil
		},
	})
}

func main() {
	db, err := gorm.Open("mysql", "root:1234@tcp(localhost:3306)/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.LogMode(true)

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(
			graphql.ObjectConfig{
				Name: "query",
				Fields: graphql.Fields{
					"users": &graphql.Field{
						Type:        graphql.NewList(userType),
						Description: "Users",
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							users := make([]dao.User, 0)
							if err := db.Preload("EMails").Preload("EMails.User").Find(&users).Error; err != nil {
								return nil, err
							}

							return users, nil
						},
					},
				},
			},
		),
	})
	if err != nil {
		panic(err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", h)
	http.ListenAndServe(":8080", nil)
}
