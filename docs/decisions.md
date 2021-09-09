#pulumi provider for concourse pipeline definitions
    - idea
        - replace ugly yaml definition with classic yaml definition shortcomings by nice pulumi style
        - use stuff as in concourse just with pulumi
        - top level pipeline resource that sets a pipeline
        - vars from pulumi config
        - steps and tasks... as classes for intellisense
            - -> reuse concourse go packages?
        - -> everything can easily be templated...

    - login:
        - probably with https://github.com/concourse/concourse/blob/91bc30439da46c104c223d7530e9ffcbff285bba/fly/rc/target.go#L311
        - nice basic auth client: https://github.com/concourse/concourse/blob/91bc30439da46c104c223d7530e9ffcbff285bba/fly/rc/target.go#L577
        - also look here: https://github.com/concourse/concourse/blob/91bc30439da46c104c223d7530e9ffcbff285bba/fly/commands/login.go#L180
        

    - what to use:
        - custom provider
            - hard to implement
            - should be used for new cloud providers -> is probably the way to go since the others always depend on resources to be created
        - dynamic provider
            - not shareable, only usable in the current code
            - e.g. used for provisioners
            - resource type is not sth like concourse:index:pipeline
            - rather used for slight changes to existing resources like migrations or provisioners than actually new resources on a new cloud provider
        - multi language component
            - based on other components to build other components -> doesn't work for auth and api calls since create, diff, delete cannot be implemented

    - autogenerate schema from atc.Config/ atc.Pipeline in concourse source code?
        - is hard since atc.Pipeline and atc.Config contain many fields that are not intended to be set by pulumi
        - probably can migrate some time to auto generating the used types here from the ones in concourse

    - what hierarchy to use
        - is the provider initialized with the team or is it a top level resource
            - is it comparable to kubernetes namespace -> resourc
            - is it rather azure subscription -> provider
            - authentication scope? if on team level -> provider

## future possible additions
- support for other concourse resources
- login
  - add option to authenticate with ClientCert
  - could also use some oauth flow