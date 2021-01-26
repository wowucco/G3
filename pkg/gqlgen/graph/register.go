package graph

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/wowucco/G3/internal/delivery"
	"github.com/wowucco/G3/internal/menu"
	"github.com/wowucco/G3/internal/product"
	"github.com/wowucco/G3/pkg/gqlgen/graph/generated"
)

func RegisterGraphql(router *gin.RouterGroup, uc product.UseCase, r product.ReadRepository, m menu.ReadRepository, d delivery.DeliveryReadRepository)  {
	//srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{useCase: uc}}))

	cnf := generated.Config{Resolvers: &Resolver{useCase: uc, productRead: r, menuRead: m, deliveryRead: d}}

	gql := router.Group("/graphql")
	{
		gql.GET("/", playgroundHandler(gql.BasePath() + "/query"))
		gql.POST("/query", graphqlHandler(cnf))
	}
}

// Defining the Graphql handler
func graphqlHandler(cfg generated.Config) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(generated.NewExecutableSchema(cfg))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler(queryPath string) gin.HandlerFunc {
	h := playground.Handler("GraphQL", queryPath)

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
