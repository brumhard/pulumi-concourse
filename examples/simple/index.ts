import * as concourse from "@pulumi/concourse";

const provider = new concourse.Provider("concourse", { url: "http://localhost:8080", username: "test", password: "test" })

const pipeline = new concourse.Pipeline(
  "my-random",
  {
    pipelineName: "my-random",
    jobs: [
      {
        name: "lel",
        plan: [
          {
            task: "lel",
            config: {
              platform: "linux",
              imageResource: {
                type: "registry-image",
                source: { "repository": "debian" }
              },
              run: {
                path: "bash",
                args: [
                  "-cex",
                  `
                  echo "this is great"
                  echo "cmon main"
                  echo "let's party"
                  `
                ]
              }
            }
          }
        ]
      }
    ]
  },
  {
    provider: provider,
  }
);


export const output = pipeline.name;