import * as concourse from "@pulumi/concourse";
import { input } from "@pulumi/concourse/types";
import { Config } from "@pulumi/pulumi";

// shared

const bashTask = (taskName: string, script: string): input.TaskStepArgs => {
  return {
    task: taskName,
    config: {
      platform: "linux",
      image_resource: {
        type: "registry-image",
        source: { "repository": "debian" }
      },
      run: {
        path: "bash",
        args: [
          "-ce",
          script
        ]
      }
    }
  }
};

const multiBranchResourceTypeKey = "git-multibranch";
const multibranchResourceType: input.ResourceTypeArgs = {
  name: multiBranchResourceTypeKey,
  type: "docker-image",
  source: {
    "repository": "cfcommunity/git-multibranch-resource"
  }
};

// actual pipeline

const cfg = new Config();
const concourseCfg = new Config("concourse")

const repoUrl = cfg.require("repoUrl")
const sshKey = cfg.require("sshKey")

const provider = new concourse.Provider("concourse", {
  url: concourseCfg.require("url"),
  username: concourseCfg.requireSecret("username"),
  password: concourseCfg.requireSecret("password")
});

const feature = "feature";
const release = "release";
const resources: input.ResourceArgs[] = [
  {
    name: feature,
    type: multiBranchResourceTypeKey,
    source: {
      "uri": repoUrl,
      "branches": "feature/.*",
      "private_key": sshKey
    }
  },
  {
    name: release,
    type: "git",
    source: {
      uri: repoUrl,
      branch: "main",
      tag_filter: "*",
      "private_key": sshKey
    }
  }
];

const jobs: input.JobArgs[] = [];
[feature, release].forEach((stage: string) => {
  jobs.push({
    name: `lint-${stage}`,
    plan: [
      {
        get: stage,
        trigger: true
      },
      bashTask("lint", `
      echo "stage: ${stage}"
      echo "linting..."
      `)
    ]
  })
  jobs.push({
    name: `build-${stage}`,
    plan: [
      {
        get: stage,
        trigger: true,
        passed: [`lint-${stage}`]
      },
      bashTask("build", `echo "building..."`)
    ]
  })
});

const pipeline = new concourse.Pipeline(
  "advanced-pipeline",
  {
    display: {
      background_image: "https://c4.wallpaperflare.com/wallpaper/786/743/628/synthwave-city-evga-wallpaper-preview.jpg"
    },
    resource_types: [
      multibranchResourceType
    ],
    resources: resources,
    jobs: [
      ...jobs,
      {
        name: "sometask",
        plan: [
          bashTask("final", `echo "this is the grand finally"; echo "ciao"`)
        ]
      }
    ]
  },
  {
    provider: provider,
  }
);