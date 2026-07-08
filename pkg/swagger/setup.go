package swagger

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupAndMount merges all docs from registry into a single spec and mounts
// the swagger-ui handler on the provided gin engine.  Call this after all
// modules have been initialized (so their docs are already registered).
func SetupAndMount(registry *Registry, engine *gin.Engine, info ServerInfo) {
	if err := MergeAndRegister(registry, info); err != nil {
		log.Printf("[swagger] merge failed (docs may be incomplete): %v", err)
	} else {
		log.Printf("[swagger] docs merged from %d module(s); available at /swagger/index.html", len(registry.Docs()))
	}
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.InstanceName(MergedSpecName)))
}
