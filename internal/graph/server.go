package graph

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

type Server struct {
	gqlServer *handler.Server
}

const defaultPort = "8080"

func (s *Server) Run(resolver *Resolver, wg sync.WaitGroup) {
	defer wg.Done()
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolver}))

	s.gqlServer = srv

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
