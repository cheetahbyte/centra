# What is this?

This is a dead-simple headless CMS. No fancy dashboards, no bloated page builders, no “enterprise-ready omnichannel synergy”.
Just a clean backend that stores your content and serves it fast. That’s it.

# Why?
Most headless CMS feel like they’re built for marketing teams with five layers of approval and a mandatory “Content Strategist” role. I don’t need that.
I just want something lightweight that stays out of my way, doesn’t force workflows on me, and simply delivers structured data to a couple of websites and apps.
This project is exactly that:
> A minimal, predictable CMS that does its job without drama.

# Setup
The most basic setup you can do to try out centra locally:
1. Create a directory `content/` and add a file to it.
`content/home.yaml`
```yaml
abc: 1
```
2. Run the docker container and mount the directory to the global content dir
```
docker run -v ./content:/content -p 3000:3000 ghcr.io/cheetahbyte/centra:latest
```
3. Access the API
```sh
 curl http://localhost:3000/api/home
```
This will return the data we just put into the file as json. Yay!
```json
{"abc":1}
```

Further [reading](https://github.com/cheetahbyte/centra/docs/Setup.md)

# Configuration
Everything is configured via environment variables.
You can read all about the configuration [here](https://github.com/cheetahbyte/centra/docs/Configuration.md)

# Content Model
If you are interested in how Centra interacts with your content, you can read about it [here](https://github.com/cheetahbyte/centra/docs/ContentModel.md).
