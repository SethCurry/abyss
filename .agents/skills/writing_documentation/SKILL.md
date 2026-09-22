---
name: writing_documentation
description: This skill contains instructions on how to update abyss' documentation, the writing style to use, and how to build and interact with the hugo site.
---

## Documentation Structure

All of abyss' documentation is contained within the `site` directory.  Most of your edits will be inside `site/content`, which
is where all of the Markdown content is stored so it can be rendered.

There are 2 major directories within `site/content`:

- `site/content/blog` contains blog posts, which are mainly release announcements
- `site/content/docs` contains all of the documentation pages

## Writing Style

The writing style to use depends on what the page is.  There are 3 writing styles I would like you to use:

- Marketing, for pages that are meant to lure users in
- Tutorials, for pages that explain how to do something to users
- Reference, for pages that explain how something works to other developers

### Tutorial Writing Style

Pages in `site/content/docs/guides` should be written in a tutorial writing style.
Use a light and cheerful tone, provide a lot of examples with in-depth explanations, and assume
that the user works in HR and accounting and knows very little about computers.

### Reference Writing Style

Pages in `site/content/docs/reference` should be written in a reference writing style.
Be brief, but thorough.  You need to provide in-depth explanations, but assume that the user
is a software developer and familiar with the basics of computers, Docker, LLMs, etc.

### Marketing Writing Style

This is the writing style to use for the main `README.md`, as well as any pages that are in `site` but not
in `site/content/docs/reference` or `site/content/docs/guides`.

These pages should be written as if they were written by a marketing agency.  Show highlights of the product,
but do not write long explanations of how to set them up.  Your goal is to get users interested in using
or contributing to `abyss`.
