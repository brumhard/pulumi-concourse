package provider

import (
	"github.com/concourse/concourse/atc"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
)
import "sigs.k8s.io/yaml"

// TODO: check if schema with camelCase works
// makePipeline creates the pipeline in concourse and returns the pipeline name
func (k *concourseProvider) makePipeline(name string, configMap map[string]interface{}) error {
	// TODO: check if it works with version set to empty string, otherwise put it into provider config
	// TODO: check how checkCredentials works
	// concourse/concourse/atc/api/configserver/save.go
	// concourse/concourse/atc/config.go

	// TODO: why doesn't this work
	logging.V(3).Infof("creating concourse pipeline")

	configBytes, err := yaml.Marshal(configMap)
	if err != nil {
		return err
	}

	_, _, warnings, err := k.client.Team(k.team).CreateOrUpdatePipelineConfig(atc.PipelineRef{Name: name}, "", configBytes, false)
	if err != nil {
		return err
	}

	if len(warnings) > 0 {
		warnString := ""
		for _, w := range warnings {
			logging.V(3).Infof("%s: %s", w.Type, w.Message)
			warnString += "     " + w.Type + ":" + w.Message
		}
		// TODO: replace this with proper warnings, otherwise it will leave behind orphaned pipelines
		return errors.New("pipeline creation had errors my dude:" + warnString)
	}

	return nil
}

func (k *concourseProvider) deletePipeline(name string) error {
	_, err := k.client.Team(k.team).DeletePipeline(atc.PipelineRef{Name: name})
	return err
}
