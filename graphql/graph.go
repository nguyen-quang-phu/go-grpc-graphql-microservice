package main

import "github.com/99designs/gqlgen/graphql"

type Server struct {
	// accountClient *account.Client
	// catalogClient *catalog.Client
	// orderClient   *order.Client
}

func NewGraphQLServer(accountUrl, catalogUrl, orderUrl string) (*Server, error) {
	// accountClient, err := account.NewClient(accountUrl)
	// if err != nil {
	// 	return nil, err
	// }
	// catalogClient, err := catalog.NewClient(catalogUrl)
	// if err != nil {
	// 	accountClient.Close()
	// 	return nil, err
	// }
	// orderClient, err := order.NewClient(orderUrl)
	// if err != nil {
	// 	accountClient.Close()
	// 	catalogClient.Close()
	// 	return nil, err
	// }

	return &Server{
		// accountClient,
		// catalogClient,
		// orderClient,
	}, nil
}

func (server *Server) Mutation() MutationResolver {
  return &mutationResolver{server}
}

func (server *Server) Query() QueryResolver {
	return &queryResolver{server}
}

func (server *Server) Account() AccountResolver {
	return &accountResolver{server}
}

func (server *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: server,
	})

}
