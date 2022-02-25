package authorization

var (
	client OpaClient
	config *Config
)

func GetClient() (OpaClient, error) {
	if client == nil {
		var err error
		config, err = NewEnvConfig()
		if err != nil {
			return nil, err
		}
		if config.Mocked {
			client = NewOpaClientMock()
		} else {
			client = NewOpaClient(config)
		}
	}

	return client, nil
}
