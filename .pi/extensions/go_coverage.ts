import {
    createAgentSession,
    type ExtensionAPI,
} from "@earendil-works/pi-coding-agent";

import util from "node:util";
import child_process from "node:child_process";
const exec = util.promisify(child_process.exec);

type CoverageLine = {
    filePath: string;
    line: number;
    functionName: string;
    coverage: number;
};
type CoverageResult = {
    lines: CoverageLine[];
    total: number;
};

async function runTestsAndCollectCoverage(): Promise<CoverageResult> {
  await exec(
      "go test ./... -coverprofile=/tmp/abyss.pi.cover.out",
  );

    const { stdout } = await exec(
        "go tool cover -func=/tmp/abyss.pi.cover.out",
    );

    const lines = stdout.split("\n")
        .map((line) => line.split("\t").filter((cell) => cell !== ""))
        .filter((line) => line.length === 3 && line[0].startsWith("github.com/SethCurry/abyss") && !line[0].startsWith("github.com/SethCurry/abyss/pkg/protobyss"))
        .map((line) => {
            const firstParts = line[0].split(":")
            return {
                filePath: firstParts[0],
                line: parseInt(firstParts[1]),
                functionName: line[1],
                coverage: parseFloat(line[2].replaceAll("%","")),
            }
        });


    console.log(lines);

    return {
        total: 0,
        lines: [],
    }
}


export default function (pi: ExtensionAPI) {
    pi.registerCommand("increase_coverage", {
        description: "Run tests, collect coverage, write tests for the lowest coverage.",
        handler: async (name, ctx) => {
            ctx.ui.notify("running tests, collecting coverage...", "info");

            const coverage = await runTestsAndCollectCoverage();

            ctx.ui.notify(
                `done running tests: ${coverage.total}`,
                "info",
            );
        },
    });
}
