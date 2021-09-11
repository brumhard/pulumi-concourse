module github.com/brumhard/pulumi-concourse/provider

go 1.16

require (
	github.com/concourse/concourse v1.6.1-0.20210910133157-8f01f1264ce7
	github.com/golang/protobuf v1.5.2
	github.com/mitchellh/mapstructure v1.4.1
	github.com/pkg/errors v0.9.1
	github.com/pulumi/pulumi/pkg/v3 v3.12.0
	github.com/pulumi/pulumi/sdk/v3 v3.12.0
	github.com/stretchr/testify v1.7.0
	golang.org/x/oauth2 v0.0.0-20210628180205-a41e5a781914
	google.golang.org/grpc v1.40.0
	sigs.k8s.io/yaml v1.2.0
)
