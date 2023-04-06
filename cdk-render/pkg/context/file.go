package context

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
)

func WriteFile(output string, cdkData *CdkContext) error {
	data := map[string]any{}

	for _, plugin := range cdkData.Plugins {
		content, err := plugin.Render(cdkData)
		if err != nil {
			return err
		}
		for k, v := range content {
			data[k] = v
		}
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(output, b, fs.FileMode(0644))
}
