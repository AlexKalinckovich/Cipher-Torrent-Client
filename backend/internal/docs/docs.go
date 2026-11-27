package docs

import (
	_ "embed"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed openapi.yaml
var openAPISpec []byte

const (
	UIPath   = "/swagger-ui"
	SpecPath = "/swagger-ui/openapi.yaml"

	htmlContentType = "text/html; charset=utf-8"
	yamlContentType = "application/yaml; charset=utf-8"

	specURLPlaceholder = "{{SPEC_URL}}"

	swaggerUITemplate = `<!doctype html>
<html>
<head>
  <title>API</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist/swagger-ui-bundle.js"></script>
  <script>SwaggerUIBundle({url:{{SPEC_URL}},dom_id:"#swagger-ui"})</script>
</body>
</html>`
)

func RegisterRoutes(rg *gin.RouterGroup) {
	specURL := path.Join("/", rg.BasePath(), SpecPath)
	page := []byte(strings.Replace(swaggerUITemplate, specURLPlaceholder, strconv.Quote(specURL), 1))

	rg.GET(SpecPath, func(c *gin.Context) {
		c.Data(http.StatusOK, yamlContentType, openAPISpec)
	})

	rg.GET(UIPath, func(c *gin.Context) {
		c.Data(http.StatusOK, htmlContentType, page)
	})
}
