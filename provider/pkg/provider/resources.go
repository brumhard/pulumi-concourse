package provider

// makePipeline creates the pipeline in concourse and returns the pipeline name
func (k *concourseProvider) makePipeline(name string, config []byte) error {
	// TODO: check if it works with version set to empty string, otherwise put it into provider config
	// TODO:
	_, _, _, err := k.client.Team(k.team).CreateOrUpdatePipelineConfig(name, "", config, false)
	k.client.Team(k.team).Pipeline()
	return err
}
