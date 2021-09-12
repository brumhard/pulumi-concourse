// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import { input as inputs, output as outputs } from "../types";

export interface AnonymousResourceArgs {
    params?: pulumi.Input<{[key: string]: pulumi.Input<string>}>;
    source: pulumi.Input<{[key: string]: pulumi.Input<string>}>;
    type: pulumi.Input<string>;
}

export interface DisplayOptionsArgs {
    /**
     * Allows users to specify a custom background image which is put at 30% opacity, grayscaled and blended into existing background. Must be an http, https, or relative URL.
     */
    background_image?: pulumi.Input<string>;
}

export interface GetStepArgs {
    get: pulumi.Input<string>;
    params?: pulumi.Input<{[key: string]: pulumi.Input<string>}>;
    passed?: pulumi.Input<pulumi.Input<string>[]>;
    resource?: pulumi.Input<string>;
    trigger?: pulumi.Input<boolean>;
}

export interface GroupArgs {
    /**
     * A list of jobs that should appear in this group. A job may appear in multiple groups. Neighbours of jobs in the current group will also appear on the same page in order to give context of the location of the group in the pipeline. You may also use any valid glob to represent several jobs.
     */
    jobs?: pulumi.Input<pulumi.Input<string>[]>;
    /**
     * A unique name for the group. This should be short and simple as it will be used as the tab name for navigation.
     */
    name: pulumi.Input<string>;
}

export interface JobArgs {
    /**
     * Step to execute regardless of whether the job succeeds, fails, errors, or aborts.
     */
    ensure?: pulumi.Input<inputs.StepArgs>;
    /**
     * If set, specifies a maximum number of builds to run at a time. If serial or serial_groups are set, they take precedence and force this value to be 1.
     */
    max_in_flight?: pulumi.Input<number>;
    /**
     * The name of the job. This should be short; it will show up in URLs.
     */
    name: pulumi.Input<string>;
    /**
     * Step to execute when the job aborts.
     */
    on_abort?: pulumi.Input<inputs.StepArgs>;
    /**
     * Step to execute when the job errors.
     */
    on_error?: pulumi.Input<inputs.StepArgs>;
    /**
     * Step to execute when the job fails.
     */
    on_failure?: pulumi.Input<inputs.StepArgs>;
    /**
     * Step to execute when the job succeeds.
     */
    on_success?: pulumi.Input<inputs.StepArgs>;
    plan: pulumi.Input<pulumi.Input<inputs.TaskStepArgs | inputs.GetStepArgs>[]>;
    /**
     * Default false. If set to true, the build log of this job will be viewable by unauthenticated users. Unauthenticated users will always be able to see the inputs, outputs, and build status history of a job. This is useful if you would like to expose your pipeline publicly without showing sensitive information in the build log.
     */
    public?: pulumi.Input<boolean>;
    /**
     * Default false. If set to true, builds will queue up and execute one-by-one, rather than executing in parallel.
     */
    serial?: pulumi.Input<boolean>;
}

export interface ResourceArgs {
    /**
     * Default 1m. The interval on which to check for new versions of the resource. Acceptable interval options are defined by the time.ParseDuration function. If set to never the resource will not be automatically checked. The resource can still be checked manually via the web UI, fly, or webhooks.
     */
    check_every?: pulumi.Input<string>;
    /**
     * The name of the resource. This should be short and simple. This name will be referenced by build plans of jobs in the pipeline.
     */
    name: pulumi.Input<string>;
    /**
     * Default false. If set to true, the metadata for each version of the resource will be viewable by unauthenticated users (assuming the pipeline has been exposed).
     */
    public?: pulumi.Input<boolean>;
    /**
     * The configuration for the resource. This varies by resource type, and is a black box to Concourse; it is blindly passed to the resource at runtime.
     */
    source: pulumi.Input<{[key: string]: pulumi.Input<string>}>;
    /**
     * Default []. A list of tags to determine which workers the checks will be performed on. You'll want to specify this if the source is internal to a worker's network, for example.
     */
    tags?: pulumi.Input<pulumi.Input<string>[]>;
    /**
     * The resource type implementing the resource.
     */
    type: pulumi.Input<string>;
    /**
     * If specified, web hooks can be sent to trigger an immediate check of the resource, specifying this value as a primitive form of authentication via query params.
     */
    webhook_token?: pulumi.Input<string>;
}

export interface ResourceTypeArgs {
    /**
     * Default 1m. The interval on which to check for new versions of the resource. Acceptable interval options are defined by the time.ParseDuration function. If set to never the resource will not be automatically checked. The resource can still be checked manually via the web UI, fly, or webhooks.
     */
    check_every?: pulumi.Input<string>;
    /**
     * The default configuration for the resource type. This varies by resource type, and is a black box to Concourse; it is merged with (duplicate fields are overwritten by) resource.source and passed to the resource at runtime.
     */
    defaults?: pulumi.Input<{[key: string]: pulumi.Input<string>}>;
    /**
     * TThe name of the resource type. This should be short and simple. This name will be referenced by pipeline.resources defined within the same pipeline, and task.image_resources used by tasks running in the pipeline.
     */
    name?: pulumi.Input<string>;
    /**
     * Arbitrary config to pass when running the get to fetch the resource type's image.
     */
    params?: pulumi.Input<{[key: string]: pulumi.Input<string>}>;
    /**
     * Default false. If set to true, the resource's containers will be run with full capabilities, as determined by the worker backend the task runs on.
     */
    privileged?: pulumi.Input<boolean>;
    /**
     * The configuration for the resource. This varies by resource type, and is a black box to Concourse; it is blindly passed to the resource at runtime.
     */
    source?: pulumi.Input<{[key: string]: pulumi.Input<string>}>;
    /**
     * Default []. A list of tags to determine which workers the checks will be performed on. You'll want to specify this if the source is internal to a worker's network, for example.
     */
    tags?: pulumi.Input<pulumi.Input<string>[]>;
    /**
     * The resource type implementing the resource.
     */
    type?: pulumi.Input<string>;
}

export interface RunArgsArgs {
    args?: pulumi.Input<pulumi.Input<string>[]>;
    dir?: pulumi.Input<string>;
    path: pulumi.Input<string>;
    user?: pulumi.Input<string>;
}

export interface StepArgs {
}

export interface TaskConfigArgs {
    image_resource: pulumi.Input<inputs.AnonymousResourceArgs>;
    platform: pulumi.Input<string>;
    run: pulumi.Input<inputs.RunArgsArgs>;
}

export interface TaskStepArgs {
    config?: pulumi.Input<inputs.TaskConfigArgs>;
    image?: pulumi.Input<string>;
    task: pulumi.Input<string>;
}
