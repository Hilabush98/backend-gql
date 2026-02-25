package main

import (
	"backend-gql/graph/generated"
	graph "backend-gql/graph/resolvers"
	"backend-gql/internal/auth"
	"backend-gql/internal/cache"
	"backend-gql/internal/db"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	env "github.com/joho/godotenv"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	err := env.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env,ERROR: ", err)
	}
	err = db.InitOracleERP("PROD")
	if err != nil {
		log.Fatal(err)
	}
	defer db.DB.Close()
	err = cache.InitCache(db.DB)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	c := generated.Config{
		Resolvers: &graph.Resolver{
			ToolsCache:      cache.GlobalMatrixcTool,
			ProfilesCache:   cache.GlobalMatrixcProfiles,
			RProfilesTools:  cache.GlobalMatrixProfilesTools,
			RGroupsProfiles: cache.GlobalMatrixrGroupsProfiles,
			GroupsCache:     cache.GlobalMatrixcGroups,
			DB:              db.DB,
		},
	}
	c.Directives.HasPermission = func(ctx context.Context, obj interface{}, next graphql.Resolver, actions []string) (interface{}, error) {
		permsData := auth.PermissionsFromContext(ctx)
		log.Println("Action 1:", actions[0], "Action 2:", actions[1])
		if permsData == nil {
			return nil, fmt.Errorf("acceso denegado: cabecera de permisos vacía")
		}

		tID := *permsData.Permissions.ToolID

		for _, p := range permsData.Permissions.ProfileID {

			pID := p
			assignedPerm := cache.GlobalMatrixProfilesTools.GetProfileByKey(pID, tID, actions)
			if assignedPerm != "" {
				log.Printf("[AUTH] Match dinámico: Profiles %s - Tool %s", tID, pID)
				return next(ctx)

			}

		}

		return nil, fmt.Errorf("acceso denegado: no se encontró relación en la matriz de permisos")
	}
	srv := handler.New(generated.NewExecutableSchema(c))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})
	router := auth.Middleware()(srv)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", router)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
