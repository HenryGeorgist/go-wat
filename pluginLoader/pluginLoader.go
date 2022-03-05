package pluginloader

import (
	"errors"

	"github.com/HenryGeorgist/go-wat/component"
)

type PluginRegistry struct {
	PluginMap map[string]*component.Computable
}

func (p *PluginRegistry) RegisterPlugin(pluginName string, plugin *component.Computable) error {
	_, exists := p.PluginMap[pluginName]
	if exists {
		return errors.New("plugin " + pluginName + " has already been registered")
	} else {
		p.PluginMap[pluginName] = plugin
		return nil
	}
}
