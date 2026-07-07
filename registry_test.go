package serverbase

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type dummyModule struct{ name string }

func (d *dummyModule) Name() string                      { return d.name }
func (d *dummyModule) RegisterRoutes(_ *gin.RouterGroup) {}
func (d *dummyModule) Migrate() error                    { return nil }

func TestModuleRegistry_RegisterAndList(t *testing.T) {
	mr := NewModuleRegistry()
	m1 := &dummyModule{name: "a"}
	m2 := &dummyModule{name: "b"}
	mr.Register(m1)
	mr.Register(m2)
	mods := mr.Modules()
	require.Len(t, mods, 2)
	require.Equal(t, "a", mods[0].Name())
	require.Equal(t, "b", mods[1].Name())
}
