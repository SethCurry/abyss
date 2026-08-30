---
title: "Developing"
description: "Advice for developing on Abyss"
summary: ""
date: 2023-09-07T16:13:18+02:00
lastmod: 2023-09-07T16:13:18+02:00
draft: false
weight: 1000
toc: true
params:
  seo:
    title: "" # custom title (optional)
    description: "" # custom description (recommended)
    canonical: "" # custom canonical URL (optional)
    robots: "" # custom robot tags (optional)
---

## Tasks

This project uses [task](https://taskfile.dev/) as its task runner.
You can find a list of the tasks [here](https://github.com/SethCurry/abyss/blob/main/Taskfile.yaml).

## The Site

You can run the dev version of the site by installing Task above and running

```bash
task site
```

Or you can just copy the bash command out of it.

The docs are auto-rebuilt by Netlify on push, so you don't need to do anything to get changes live on the site.
