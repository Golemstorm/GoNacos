package nacos_registry

import (
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/web"
)

func WebRegistry(r registry.Registry) web.Option {
	return func(o *web.Options) {
		o.Registry = r
		if o.Metadata == nil {
			o.Metadata = make(map[string]string)
		}
		o.Metadata["broker"] = "http"
		o.Metadata["http"] = "http"
		o.Metadata["registry"] = "nacos"
		o.Metadata["server"] = "http"
		o.Metadata["transport"] = "http"
	}
}
