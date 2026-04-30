package main

import (
	"backend-gql/graph/generated"
	graph "backend-gql/graph/resolvers"
	"backend-gql/internal/auth"
	"backend-gql/internal/cache"
	"backend-gql/internal/db"
	"backend-gql/internal/logs"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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
	err := logs.InitFilesLogs()
	if err != nil {
		logs.Error("cmd/server/main", fmt.Sprintf("Error: %v", err))
	}
	logs.Debug("server.go", "Iniciando server")

	err = env.Load("cmd/.env")
	if err != nil {
		logs.Error("cmd/server/main", fmt.Sprintf("Error al cargar el archivo .env: %v", err))
	}



	err = db.InitPostgresAuth()
	if err != nil {
		logs.Error("cmd/server/main", fmt.Sprintf("Error al iniciar la base de datos Auth: %v", err))
		log.Fatalf("Error al iniciar la base de datos Auth: %v", err)
	}
	if db.DBAuth != nil {
		defer db.DBAuth.Close()
	}

	err = db.InitPostgresApp()
	if err != nil {
		logs.Error("cmd/server/main", fmt.Sprintf("Error al iniciar la base de datos App: %v", err))
		log.Fatalf("Error al iniciar la base de datos App: %v", err)
	}
	if db.DBApp != nil {
		defer db.DBApp.Close()
	}

	logs.Debug("main", "Base de datos cargada")
	err = cache.InitCache(db.DBApp)
	if err != nil {
		logs.Error("cmd/server/main", fmt.Sprintf("Error al iniciar cache: %v", err))
		log.Fatalf("Error al iniciar cache: %v", err)
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
			DB:              db.DBApp,
			DBAuth:          db.DBAuth,
		},
	}

	c.Directives.HasPermission = func(ctx context.Context, obj interface{}, next graphql.Resolver, actions []string) (interface{}, error) {
		permsData := auth.PermissionsFromContext(ctx)
		if permsData == nil || permsData.Permissions.ToolID == nil {
			logs.Error("server.go", "Access Denied: invalid auth payload")
			return nil, fmt.Errorf("access denied: invalid auth payload")
		}

		tID := strings.TrimSpace(*permsData.Permissions.ToolID)
		if tID == "" {
			logs.Error("server.go", "Access Denied: tool_id is empty")
			return nil, fmt.Errorf("access denied: tool_id is empty")
		}

		for _, p := range permsData.Permissions.ProfileID {
			pID := strings.TrimSpace(p)
			if pID == "" {
				continue
			}

			assignedPerm := cache.GlobalMatrixProfilesTools.GetProfileByKey(pID, tID, actions)
			if assignedPerm != "" {
				return next(ctx)
			}
		}

		logs.Error("server.go", "Access Denied: Role without permissions")
		return nil, fmt.Errorf("access denied: role without permissions")
	}

	srv := handler.New(generated.NewExecutableSchema(c))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{Cache: lru.New[string](100)})

	router := auth.Middleware()(srv)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", router)

	logs.Debug("main.go", fmt.Sprintf("connect to http://localhost:%s/ for GraphQL playground", port))
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
