package serverbase

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
)

// GinServer is a simple Server implementation using gin.Engine.
type GinServer struct {
    engine *gin.Engine
    srv    *http.Server
}

// NewGinServer creates a new GinServer with default middleware.
func NewGinServer() *GinServer {
    e := gin.New()
    e.Use(gin.Recovery())
    return &GinServer{engine: e}
}

// Engine returns the underlying gin engine.
func (s *GinServer) Engine() *gin.Engine { return s.engine }

// Start runs the server on the given address and handles graceful shutdown.
func (s *GinServer) Start(addr string) error {
    if addr == "" {
        addr = ":8085"
    }

    s.srv = &http.Server{
        Addr:    addr,
        Handler: s.engine,
    }

    // Start server in goroutine
    go func() {
        if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Printf("server listen error: %v", err)
        }
    }()

    // Wait for interrupt signal to gracefully shutdown the server with a timeout.
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    <-ctx.Done()
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    return s.srv.Shutdown(shutdownCtx)
}

// WireModules registers modules on the server under root group.
func WireModules(s Server, mr *ModuleRegistry) error {
    rg := s.Engine().Group("/api/v1")
    return mr.RegisterAll(rg)
}

// RunWithModulesAndConfig runs modules, optionally migrates, and starts server with graceful shutdown.
func RunWithModulesAndConfig(cfg *Config, mr *ModuleRegistry) error {
    s := NewGinServer()

    // Register routes
    if err := WireModules(s, mr); err != nil {
        return fmt.Errorf("wire modules: %w", err)
    }

    // Run migrations if requested
    if cfg != nil && cfg.MigrateOnStart {
        for _, m := range mr.Modules() {
            if err := m.Migrate(); err != nil {
                return fmt.Errorf("migrate module %s: %w", m.Name(), err)
            }
        }
    }

    // Start server (blocks until shutdown)
    return s.Start(cfg.Port)
}
