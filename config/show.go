package config

// ShowConfig represents the configuration for the show command
type ShowConfig struct {
	Region    string `json:"region"`
	AccountID string `json:"accountID"`
}

// GetShowConfig returns the configuration for the show command
func GetShowConfig() (config ShowConfig, err error) {
	awsInfo, err := GetAWSInfo()
	if err != nil {
		return config, err

	}
	config = ShowConfig{
		Region:    awsInfo.Region,
		AccountID: awsInfo.AccountID,
	}
	return config, nil
}
