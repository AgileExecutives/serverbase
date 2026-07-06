package serverbase

import "github.com/gin-gonic/gin"

// Server is a minimal server abstraction used by modules to register routes.
type Server interface {
    Engine() *gin.Engine
    Start(addr string) error
}

// Module is the interface each module should implement to register routes and run migrations.
type Module interface {
    Name() string
    RegisterRoutes(rg *gin.RouterGroup)
    Migrate() error
}
