package utils

import (
	"bon-voyage-agent/shared"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

func LoadPlugins(pluginDir string) (map[string]shared.Plugin, error) {

	plugins := make(map[string]shared.Plugin)
	// Walk through the plugin directory and find .so files.
	err := filepath.Walk(pluginDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is a .so file.
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".so") {
			// Load the plugin.
			p, err := plugin.Open(path)
			if err != nil {
				log.Printf("Error loading plugin %s: %v", path, err)
				return nil // Continue loading other plugins
			}

			// Lookup the PluginInstance symbol.
			symbol, err := p.Lookup("PluginInstance")
			if err != nil {
				log.Printf("Error looking up PluginInstance in %s: %v", path, err)
				return nil // Continue loading other plugins
			}

			// Assert the symbol as the Plugin interface.
			pluginInstance, ok := symbol.(shared.Plugin)
			if !ok {
				log.Printf("Error asserting PluginInstance in %s as Plugin interface", path)
				return nil // Continue loading other plugins
			}

			// Extract the plugin name from the file name.
			pluginName := strings.TrimSuffix(strings.TrimPrefix(info.Name(), "plugin_"), ".so")
			plugins[pluginName] = pluginInstance
			log.Printf("Loaded plugin: %s", pluginName)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return plugins, nil
}
