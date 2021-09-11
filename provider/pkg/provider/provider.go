// Copyright 2016-2020, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/go-concourse/concourse"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"

	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	pbempty "github.com/golang/protobuf/ptypes/empty"
)

// TODO: get this from concourse package somewhere?
var concourseDefaultTeam = "main"

type concourseProvider struct {
	host        *provider.HostClient
	name        string
	version     string
	schemaBytes []byte // initialize in makeProvider
	config      map[string]string
	client      concourse.Client
	team        string
}

func makeProvider(host *provider.HostClient, name, version string, schemaBytes []byte) (pulumirpc.ResourceProviderServer, error) {
	// Return the new provider
	return &concourseProvider{
		host:        host,
		name:        name,
		version:     version,
		schemaBytes: schemaBytes,
		config:      map[string]string{},
	}, nil
}

// Configure configures the resource provider with "globals" that control its behavior.
// It sets all required properties on the provider that are not set in makeProvider and come from the
func (k *concourseProvider) Configure(ctx context.Context, req *pulumirpc.ConfigureRequest) (
	*pulumirpc.ConfigureResponse, error,
) {
	for key, val := range req.GetVariables() {
		k.config[strings.TrimPrefix(key, "concourse:config:")] = val
	}

	k.setLoggingContext(ctx)

	// TODO: setup logging
	if err := k.getClient(); err != nil {
		return nil, err
	}

	k.getTeam()

	return &pulumirpc.ConfigureResponse{}, nil
}

func (k *concourseProvider) setLoggingContext(ctx context.Context) {
	log.SetOutput(NewLogRedirector(ctx, k.host))
}

func (k *concourseProvider) getClient() error {
	url := k.getConfig("url", "CONCOURSE_URL")
	username := k.getConfig("username", "CONCOURSE_USERNAME")
	password := k.getConfig("password", "CONCOURSE_PASSWORD")

	client, err := NewClient(url, username, password)
	if err != nil {
		return err
	}

	k.client = client
	return nil
}

func (k *concourseProvider) getTeam() {
	team := k.getConfig("team", "CONCOURSE_TEAM")
	if team == "" {
		team = concourseDefaultTeam
	}

	k.team = team
}

// TODO: the env vars queried with this function should match the ones defined in schema.json provider section
func (k *concourseProvider) getConfig(configName, envName string) string {
	if val, ok := k.config[configName]; ok {
		return val
	}

	return os.Getenv(envName)
}

// Check validates that the given property bag is valid for a resource of the given type and returns
// the inputs that should be passed to successive calls to Diff, Create, or Update for this
// resource. As a rule, the provider inputs returned by a call to Check should preserve the original
// representation of the properties as present in the program inputs. Though this rule is not
// required for correctness, violations thereof can negatively impact the end-user experience, as
// the provider inputs are using for detecting and rendering diffs.
//
// This can be used to also apply defaults to the resources if there are any.
// TODO: should autoNaming be moved here to apply it as default or should the feature be removed in general
// 	- probably remove it, the only reason is faster replace actions for things that need some time for deletion, which is not the case here
func (k *concourseProvider) Check(ctx context.Context, req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if ty != "concourse:index:Pipeline" {
		return nil, fmt.Errorf("Unknown resource type '%s'", ty)
	}

	news, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	// make sure that final name does not contain any colon, since that is used in teamResourceID
	explicitlySet, name := autoName(urn, news)

	property := "name"
	if explicitlySet {
		property = "pipelineName"
	}

	if strings.ContainsRune(name, ':') {
		return &pulumirpc.CheckResponse{Inputs: req.News, Failures: []*pulumirpc.CheckFailure{
			{
				Property: property,
				Reason:   "must not contain a colon",
			},
		}}, nil
	}

	return &pulumirpc.CheckResponse{Inputs: req.News, Failures: nil}, nil
}

// Diff checks what impacts a hypothetical update will have on the resource's properties.
// TODO: apparently there is a diff function in the pipeline config struct
// https://github.com/concourse/concourse/blob/master/atc/config_diff.go#L277
func (k *concourseProvider) Diff(ctx context.Context, req *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if ty != "concourse:index:Pipeline" {
		return nil, fmt.Errorf("Unknown resource type '%s'", ty)
	}

	olds, err := plugin.UnmarshalProperties(req.GetOlds(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	news, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	d := olds.Diff(news)
	changes := pulumirpc.DiffResponse_DIFF_NONE
	if d.Changed("length") {
		changes = pulumirpc.DiffResponse_DIFF_SOME
	}

	return &pulumirpc.DiffResponse{
		Changes:  changes,
		Replaces: []string{"length"},
	}, nil
}

// Create allocates a new instance of the provided resource and returns its unique ID afterwards.
// TODO: could be somewhat equal to https://github.com/concourse/concourse/blob/91bc30439da46c104c223d7530e9ffcbff285bba/fly/commands/internal/setpipelinehelpers/atc_config.go#L48
func (k *concourseProvider) Create(ctx context.Context, req *pulumirpc.CreateRequest) (*pulumirpc.CreateResponse, error) {
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if ty != "concourse:index:Pipeline" {
		return nil, fmt.Errorf("Unknown resource type '%s'", ty)
	}

	inputs, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	_, name := autoName(urn, inputs)

	inputMap := inputs.Mappable()
	var config atc.Config
	if err := mapstructure.Decode(inputMap, &config); err != nil {
		return nil, err
	}

	// Actually "create" the pipeline
	// TODO: check if you could also just use a map here
	if err := k.makePipeline(name, config); err != nil {
		return nil, err
	}

	outputs := map[string]interface{}{
		"name": name,
	}

	outputProperties, err := plugin.MarshalProperties(
		resource.NewPropertyMapFromMap(outputs),
		plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true},
	)
	if err != nil {
		return nil, err
	}
	return &pulumirpc.CreateResponse{
		Id:         k.teamResourceID(name),
		Properties: outputProperties,
	}, nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

// autoName returns a bool indicating whether the name was explicitly set or not and the name that should be used.
func autoName(urn resource.URN, propertyMap resource.PropertyMap) (bool, string) {
	var propertyName string

	switch urn.Type() {
	case "concourse:index:Pipeline":
		propertyName = "pipelineName"
	default:
		panic(fmt.Sprintf("type %s does not support autoNaming", urn.Type()))
	}

	propKey := resource.PropertyKey(propertyName)

	if propertyMap.HasValue(propKey) {
		return true, propertyMap[propKey].StringValue()
	}

	b := make([]byte, 8)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}

	return false, fmt.Sprintf("%s-%s", urn.Name().String(), b)
}

const separator = ':'

func (k *concourseProvider) teamResourceID(name string) string {
	return fmt.Sprintf("%s%c%s", k.team, separator, name)
}

func nameFromTeamResourceID(teamResourceID string) string {
	parts := strings.Split(teamResourceID, string(separator))
	if len(parts) != 2 {
		panic("teamResourceID should be in the format %s")
	}

	return parts[1]
}

// Read the current live state associated with a resource.
func (k *concourseProvider) Read(ctx context.Context, req *pulumirpc.ReadRequest) (*pulumirpc.ReadResponse, error) {
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if ty != "concourse:index:Pipeline" {
		return nil, fmt.Errorf("Unknown resource type '%s'", ty)
	}
	return nil, status.Error(codes.Unimplemented, "Read is not yet implemented for 'concourse:index:Pipeline'")
}

// Update updates an existing resource with new values.
func (k *concourseProvider) Update(ctx context.Context, req *pulumirpc.UpdateRequest) (*pulumirpc.UpdateResponse, error) {
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if ty != "concourse:index:Pipeline" {
		return nil, fmt.Errorf("Unknown resource type '%s'", ty)
	}

	// Our Random resource will never be updated - if there is a diff, it will be a replacement.
	return nil, status.Error(codes.Unimplemented, "Update is not yet implemented for 'concourse:index:Pipeline'")
}

// Delete tears down an existing resource with the given ID.  If it fails, the resource is assumed
// to still exist.
func (k *concourseProvider) Delete(ctx context.Context, req *pulumirpc.DeleteRequest) (*pbempty.Empty, error) {
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if ty != "concourse:index:Pipeline" {
		return nil, fmt.Errorf("Unknown resource type '%s'", ty)
	}

	pipelineName := nameFromTeamResourceID(req.GetId())

	if err := k.deletePipeline(pipelineName); err != nil {
		return nil, err
	}

	// Note that for our Random resource, we don't have to do anything on Delete.
	return &pbempty.Empty{}, nil
}

// Invoke dynamically executes a built-in function in the provider.
func (k *concourseProvider) Invoke(_ context.Context, req *pulumirpc.InvokeRequest) (*pulumirpc.InvokeResponse, error) {
	tok := req.GetTok()
	// TODO: implement functions defined in schema functions e.g. concourse:index:getPipeline
	return nil, fmt.Errorf("Unknown Invoke token '%s'", tok)
}

// StreamInvoke dynamically executes a built-in function in the provider. The result is streamed
// back as a series of messages.
func (k *concourseProvider) StreamInvoke(req *pulumirpc.InvokeRequest, server pulumirpc.ResourceProvider_StreamInvokeServer) error {
	// NOTE: same as in azure-native provider
	return status.Error(codes.Unimplemented, "StreamInvoke is not yet implemented")
}

// GetPluginInfo returns generic information about this plugin, like its version.
func (k *concourseProvider) GetPluginInfo(context.Context, *pbempty.Empty) (*pulumirpc.PluginInfo, error) {
	// NOTE: same as in azure-native provider
	return &pulumirpc.PluginInfo{
		Version: k.version,
	}, nil
}

// TODO: define the schema in go code to get away the oneOf repeats
// GetSchema returns the JSON-serialized schema for the provider.
func (k *concourseProvider) GetSchema(ctx context.Context, req *pulumirpc.GetSchemaRequest) (*pulumirpc.GetSchemaResponse, error) {
	// NOTE: same as in azure-native provider
	if v := req.GetVersion(); v != 0 {
		return nil, fmt.Errorf("unsupported schema version %d", v)
	}

	uncompressed, err := gzip.NewReader(bytes.NewReader(k.schemaBytes))
	if err != nil {
		return nil, errors.Wrap(err, "expand compressed bytes for schema")
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, uncompressed)
	if err != nil {
		return nil, errors.Wrap(err, "closing read stream for schema")
	}

	if err = uncompressed.Close(); err != nil {
		return nil, errors.Wrap(err, "closing uncompress stream for schema")
	}

	return &pulumirpc.GetSchemaResponse{Schema: buf.String()}, nil
}

// Cancel signals the provider to gracefully shut down and abort any ongoing resource operations.
// Operations aborted in this way will return an error (e.g., `Update` and `Create` will either a
// creation error or an initialization error). Since Cancel is advisory and non-blocking, it is up
// to the host to decide how long to wait after Cancel is called before (e.g.)
// hard-closing any gRPC connection.
func (k *concourseProvider) Cancel(context.Context, *pbempty.Empty) (*pbempty.Empty, error) {
	// NOTE: same as in azure-native provider
	// TODO: implement
	return &pbempty.Empty{}, nil
}

// Call dynamically executes a method in the provider associated with a component resource.
func (k *concourseProvider) Call(ctx context.Context, req *pulumirpc.CallRequest) (*pulumirpc.CallResponse, error) {
	// NOTE: same as in azure-native provider
	return nil, status.Error(codes.Unimplemented, "Call is not yet implemented")
}

// Construct creates a new component resource.
func (k *concourseProvider) Construct(ctx context.Context, req *pulumirpc.ConstructRequest) (*pulumirpc.ConstructResponse, error) {
	/// NOTE: same as in azure-native provider
	return nil, status.Error(codes.Unimplemented, "Construct is not yet implemented")
}

// CheckConfig validates the configuration for this provider.
func (k *concourseProvider) CheckConfig(ctx context.Context, req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	// NOTE: same as in azure-native provider
	return &pulumirpc.CheckResponse{Inputs: req.GetNews()}, nil
}

// DiffConfig diffs the configuration for this provider.
func (k *concourseProvider) DiffConfig(ctx context.Context, req *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	// NOTE: same as in azure-native provider
	return &pulumirpc.DiffResponse{
		Changes:             0,
		Replaces:            []string{},
		Stables:             []string{},
		DeleteBeforeReplace: false,
	}, nil
}
