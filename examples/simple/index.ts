import * as concourse from "@pulumi/concourse";

const random = new concourse.Random("my-random", { length: 24 });

export const output = random.result;