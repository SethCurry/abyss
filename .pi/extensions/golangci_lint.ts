import {
    createAgentSession,
    type ExtensionAPI,
} from "@earendil-works/pi-coding-agent";

import util from "node:util";
import child_process from "node:child_process";
const exec = util.promisify(child_process.exec);

type GolangciLintIssue = {
    FromLinter: string;
    Text: string;
    Severity: string;
    SourceLines: string[];
    Pos: {
        Filename: string;
        Offset: number;
        Line: number;
        Column: number;
    };
};
type GolangciLintResult = {
    Issues: GolangciLintIssue[];
};

export default function (pi: ExtensionAPI) {
    pi.registerCommand("golangci_lint", {
        description: "Run golangci-lint",
        handler: async (num_issues, ctx) => {
            ctx.ui.notify("running linter", "info");

            const { stdout } = await exec(
                "golangci-lint run --output.json.path=stdout --show-stats=false --issues-exit-code=0",
            );
            const lintResult = JSON.parse(stdout) as GolangciLintResult;

            ctx.ui.notify(
                `done running linter, found ${lintResult.Issues.length} issues`,
                "info",
            );

            var issues = lintResult.Issues;

            const numIssues = num_issues ? parseInt(num_issues) : 5;

            if (issues.length > numIssues) {
                issues = issues.slice(0, numIssues);
            }

            const issueSummaries = issues.map(
                (issue) =>
                    `${issue.Pos.Filename}:${issue.Pos.Line}:${issue.Pos.Column} ${issue.Text}`,
            );
            const issueSummary = issueSummaries.join("\n");

            const shouldProceed = await ctx.ui.confirm(
                "Have agent fix issues?",
                `I am going to try to fix ${issues.length} issues, should I proceed?\n\n${issueSummary}`,
            );

            if (!shouldProceed) {
                return;
            }

            for (const issue of issues) {
                const message = `Trying to fix: ${issue.Pos.Filename}:${issue.Pos.Line}:${issue.Pos.Column} ${issue.Text}\n`;
                ctx.ui.notify(message, "info");
                const { session } = await createAgentSession();

                try {
                    const promptMessage = `I need to fix the following linter error, can you help me?  This is my error:
        ${issue.Pos.Filename}:${issue.Pos.Line}:${issue.Pos.Column} ${issue.Text}
        `;
                    await session.prompt(promptMessage);
                    const lastMessage = session.getLastAssistantText();
                    if (lastMessage) {
                        ctx.ui.notify(lastMessage, "info");
                    }
                } finally {
                    session.dispose();
                }
            }
            ctx.ui.notify("Done fixing lint issues", "info");
        },
    });
}
