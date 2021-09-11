package provider

import "github.com/concourse/concourse/atc"
import "sigs.k8s.io/yaml"

// TODO: check if schema with camelCase works
// makePipeline creates the pipeline in concourse and returns the pipeline name
func (k *concourseProvider) makePipeline(name string, config atc.Config) error {
	// TODO: check if it works with version set to empty string, otherwise put it into provider config
	// TODO: check how checkCredentials works
	// concourse/concourse/atc/api/configserver/save.go
	// concourse/concourse/atc/config.go

	configBytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	_, _, _, err = k.client.Team(k.team).CreateOrUpdatePipelineConfig(atc.PipelineRef{Name: name}, "", configBytes, false)
	return err
}

func (k *concourseProvider) deletePipeline(name string) error {
	_, err := k.client.Team(k.team).DeletePipeline(atc.PipelineRef{Name: name})
	return err
}
