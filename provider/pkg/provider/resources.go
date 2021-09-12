package provider

import (
	"fmt"

	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/atc/configvalidate"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	"sigs.k8s.io/yaml"
)

// makeOrUpdatePipeline creates the pipeline in concourse and returns the pipeline name
func (k *concourseProvider) makeOrUpdatePipeline(name string, configMap map[string]interface{}) error {
	// TODO: why doesn't this work
	logging.V(3).Infof("creating concourse pipeline")

	teamClient := k.client.Team(k.team)
	pipelineRef := atc.PipelineRef{Name: name}

	// no error is returned in case no pipeline is found
	_, existingConfigVersion, _, err := teamClient.PipelineConfig(pipelineRef)
	if err != nil {
		return err
	}

	configBytes, err := yaml.Marshal(configMap)
	if err != nil {
		return err
	}

	var newConfig atc.Config
	err = yaml.Unmarshal(configBytes, &newConfig)
	if err != nil {
		return err
	}

	atcWarnings, errorMessages := configvalidate.Validate(newConfig)
	if len(errorMessages) > 0 {
		errString := "errs during execution: "
		for _, errMessage := range errorMessages {
			errString += fmt.Sprintf("%s, ", errMessage)
		}
		return fmt.Errorf(errString)
	}

	for _, w := range atcWarnings {
		logging.Warningf("%s: %s", w.Type, w.Message)
	}

	_, _, warnings, err := teamClient.CreateOrUpdatePipelineConfig(pipelineRef, existingConfigVersion, configBytes, false)
	if err != nil {
		return err
	}

	if len(warnings) > 0 {
		warnString := "warnings during execution: "
		for _, w := range warnings {
			logging.Warningf("%s: %s", w.Type, w.Message)
			warnString += fmt.Sprintf("{type: %s, message: %s}, ", w.Type, w.Message)
		}
		// TODO: replace this with proper warnings, otherwise it will leave behind orphaned pipelines
		return fmt.Errorf(warnString)
	}

	return nil
}

func (k *concourseProvider) deletePipeline(name string) error {
	_, err := k.client.Team(k.team).DeletePipeline(atc.PipelineRef{Name: name})
	return err
}
