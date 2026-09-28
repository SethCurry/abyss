This is the base image that other JS-based images are built from.
It's really just `base` with `npm` and `node` installed.

It exists solely to prevent apt-get install'ing those across 3 different images.
It makes build times for testing terribly long.

If you're looking at building your own image, check out the
[guide on the site](https://abyss.scurry.io/docs/guides/custom-docker-images/).
