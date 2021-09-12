import * as concourse from "@pulumi/concourse";

const provider = new concourse.Provider("concourse", { url: "http://localhost:8080", username: "test", password: "test" })

const pipeline = new concourse.Pipeline(
  "testing-pipeline",
  {
    jobs: [
      {
        name: "lel",
        plan: [
          {
            task: "lel",
            config: {
              platform: "linux",
              image_resource: {
                type: "registry-image",
                source: { "repository": "debian" }
              },
              run: {
                path: "bash",
                args: [
                  "-cex",
                  `
                  echo "this is garbage"
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