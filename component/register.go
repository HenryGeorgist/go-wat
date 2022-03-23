package component

import (
	"errors"
)

type PluginRegistry struct {
	PluginMap map[string]*Computable
}

func (p *PluginRegistry) RegisterPlugin(pluginName string, plugin *Computable) error {
	_, exists := p.PluginMap[pluginName]
	if exists {
		return errors.New("plugin " + pluginName + " has already been registered")
	} else {
		p.PluginMap[pluginName] = plugin
		return nil
	}
}
