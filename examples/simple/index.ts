import * as concourse from "@pulumi/concourse";

const pipeline = new concourse.Pipeline("my-random", { length: 24 });

const provider = new concourse.Provider("concourse", {url: "http://localhost:8080", username: "test", password: "test"})

export const output = pipeline.result;