import * as concourse from "@pulumi/concourse";

const pipeline = new concourse.Pipeline(
    "my-random",
    {
        pipelineName: "lel",
        resources: [
            {
                name: "test",
                type: "git",
                source: {
                    "uri": "",
                    "branch": "",
                    "private_key": "",
                }
            },
            {
                name: "image",
                type: "registry-image",
                source: {

                }
            }
        ],
        jobs: [
            {
                name: "lel",
                plan: [
                    {

                    }
                ]
            }
        ]
    }
);

const provider = new concourse.Provider("concourse", {url: "http://localhost:8080", username: "test", password: "test"})

export const output = pipeline.name;