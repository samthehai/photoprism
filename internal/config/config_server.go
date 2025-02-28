package config

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/photoprism/photoprism/internal/server/header"
	"github.com/photoprism/photoprism/pkg/fs"
)

const (
	HttpModeProd  = "release"
	HttpModeDebug = "debug"
)

// DetachServer checks if server should detach from console (daemon mode).
func (c *Config) DetachServer() bool {
	return c.options.DetachServer
}

// TrustedProxy returns the list of trusted proxy servers as comma-separated list.
func (c *Config) TrustedProxy() string {
	return strings.Join(c.options.TrustedProxies, ", ")
}

// TrustedProxies returns proxy server ranges from which client and protocol headers can be trusted.
func (c *Config) TrustedProxies() []string {
	return c.options.TrustedProxies
}

// HttpsProxyHeader returns the proxy protocol header names.
func (c *Config) HttpsProxyHeader() []string {
	return c.options.HttpsProxyHeaders
}

// HttpsProxyProto returns the proxy protocol header HTTPS values.
func (c *Config) HttpsProxyProto() []string {
	return c.options.HttpsProxyProto
}

// HttpsProxyHeaders returns a map with the proxy https protocol headers.
func (c *Config) HttpsProxyHeaders() map[string]string {
	p := len(c.options.HttpsProxyHeaders)
	h := make(map[string]string, p+1)

	if p == 0 {
		h[header.ForwardedProto] = header.ProtoHttps
		return h
	}

	for k, v := range c.options.HttpsProxyHeaders {
		if l := len(c.options.HttpsProxyProto); l == 0 {
			h[v] = header.ProtoHttps
		} else if l > k {
			h[v] = c.options.HttpsProxyProto[k]
		} else {
			h[v] = c.options.HttpsProxyProto[0]
		}
	}

	return h
}

// HttpMode returns the server mode.
func (c *Config) HttpMode() string {
	if c.Prod() {
		return HttpModeProd
	} else if c.options.HttpMode == "" {
		if c.Debug() {
			return HttpModeDebug
		}

		return HttpModeProd
	}

	return c.options.HttpMode
}

// HttpCompression returns the http compression method (none or gzip).
func (c *Config) HttpCompression() string {
	return strings.ToLower(strings.TrimSpace(c.options.HttpCompression))
}

// HttpHost returns the built-in HTTP server host name or IP address (empty for all interfaces).
func (c *Config) HttpHost() string {
	if c.options.HttpHost == "" {
		return "0.0.0.0"
	}

	return c.options.HttpHost
}

// HttpPort returns the HTTP server port number.
func (c *Config) HttpPort() int {
	if c.options.HttpPort == 0 {
		return 2342
	}

	return c.options.HttpPort
}

// TemplatesPath returns the server templates path.
func (c *Config) TemplatesPath() string {
	return filepath.Join(c.AssetsPath(), "templates")
}

// CustomTemplatesPath returns the path to custom templates.
func (c *Config) CustomTemplatesPath() string {
	if p := c.CustomAssetsPath(); p != "" {
		return filepath.Join(p, "templates")
	}

	return ""
}

// TemplateFiles returns the file paths of all templates found.
func (c *Config) TemplateFiles() []string {
	results := make([]string, 0, 32)

	tmplPaths := []string{c.TemplatesPath(), c.CustomTemplatesPath()}

	for _, p := range tmplPaths {
		matches, err := filepath.Glob(regexp.QuoteMeta(p) + "/[A-Za-z0-9]*.*")

		if err != nil {
			continue
		}

		for _, tmplName := range matches {
			results = append(results, tmplName)
		}
	}

	return results
}

// TemplateExists checks if a template with the given name exists (e.g. index.tmpl).
func (c *Config) TemplateExists(name string) bool {
	if found := fs.FileExists(filepath.Join(c.TemplatesPath(), name)); found {
		return true
	} else if p := c.CustomTemplatesPath(); p != "" {
		return fs.FileExists(filepath.Join(p, name))
	} else {
		return false
	}
}

// TemplateName returns the name of the default template (e.g. index.tmpl).
func (c *Config) TemplateName() string {
	if s := c.Settings(); s != nil {
		if c.TemplateExists(s.Templates.Default) {
			return s.Templates.Default
		}
	}

	return "index.tmpl"
}

// StaticPath returns the static assets' path.
func (c *Config) StaticPath() string {
	return filepath.Join(c.AssetsPath(), "static")
}

// StaticFile returns the path to a static file.
func (c *Config) StaticFile(fileName string) string {
	return filepath.Join(c.AssetsPath(), "static", fileName)
}

// BuildPath returns the static build path.
func (c *Config) BuildPath() string {
	return filepath.Join(c.StaticPath(), "build")
}

// ImgPath returns the path to static image files.
func (c *Config) ImgPath() string {
	return filepath.Join(c.StaticPath(), "img")
}
