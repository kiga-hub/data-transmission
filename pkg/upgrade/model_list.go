package upgrade

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// SourceConfig -
type SourceConfig struct {
	ProjectName string `json:"project_name"`
	Source      string `json:"source"`
}

// ModelConfig -
type ModelConfig struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// GetSourceConfig -
func GetSourceConfig(localPath string) []*SourceConfig {
	jsonFile, err := os.Open(localPath)
	if err != nil {
		log.Fatalf("Failed to open JSON file: %s", err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("Failed to read JSON file: %s", err)
	}

	var config []*SourceConfig
	if err := json.Unmarshal(byteValue, &config); err != nil {
		log.Fatalf("Failed to unmarshal JSON data: %s", err)
	}

	return config
}

// GetSourceConfigByProjectName -
func GetSourceConfigByProjectName(projectName string) *SourceConfig {
	configs := GetSourceConfig("./source.json")
	for _, config := range configs {
		if config.ProjectName == projectName {
			return config
		}
	}
	return nil
}

// GetModelConfig -
func GetModelConfig() []ModelConfig {
	jsonFile, err := os.Open("model_config.json")
	if err != nil {
		log.Fatalf("Failed to open JSON file: %s", err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("Failed to read JSON file: %s", err)
	}

	var configItems []ModelConfig
	if err := json.Unmarshal(byteValue, &configItems); err != nil {
		log.Fatalf("Failed to unmarshal JSON data: %s", err)
	}

	return configItems
}

// SetModelConfig -
func SetModelConfig(sourceConfig SourceConfig) error {
	// Read the JSON file
	file, err := os.ReadFile("model_config.json")
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Parse the JSON data
	var configs []ModelConfig
	err = json.Unmarshal(file, &configs)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// // Modify the value of the key "host"
	// for i, config := range configs {
	// 	if config.Key == "host" {
	// 		configs[i].Value = sourceConfig.HostIP
	// 		break
	// 	}
	// }

	// Marshal the modified data back to JSON
	modifiedData, err := json.MarshalIndent(configs, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Write the modified data back to the file
	err = os.WriteFile("model_config.json", modifiedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
