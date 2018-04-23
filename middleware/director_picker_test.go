package customMiddleware

import (
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/labstack/echo"

	"github.com/vicanso/pike/proxy"
	"github.com/vicanso/pike/vars"
)

func TestUpstreamPicker(t *testing.T) {
	directors := make(proxy.Directors, 0)
	aslant := "aslant"
	d := &proxy.Director{
		Name: aslant,
		Hosts: []string{
			"(www.)?aslant.site",
		},
	}
	directors = append(directors, d)

	tiny := "tiny"
	d = &proxy.Director{
		Name: tiny,
		Prefixs: []string{
			"/api",
		},
	}
	directors = append(directors, d)
	for _, d := range directors {
		d.RefreshPriority()
	}
	sort.Sort(directors)
	t.Run("get director match host", func(t *testing.T) {
		fn := DirectorPicker(directors)(func(c echo.Context) error {
			d := c.Get(vars.Director).(*proxy.Director)
			if d.Name != aslant {
				t.Fatalf("get director match host fail")
			}
			return nil
		})
		e := echo.New()
		req := httptest.NewRequest(echo.GET, "http://aslant.site/api/users/me", nil)
		c := e.NewContext(req, nil)
		fn(c)
	})

	t.Run("get director match url prefix", func(t *testing.T) {
		fn := DirectorPicker(directors)(func(c echo.Context) error {
			d := c.Get(vars.Director).(*proxy.Director)
			if d.Name != tiny {
				t.Fatalf("get director match url prefix fail")
			}
			return nil
		})
		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/api/users/me", nil)
		c := e.NewContext(req, nil)
		fn(c)
	})

}
