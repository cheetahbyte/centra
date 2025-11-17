# What is this?
This is a dead-simple headless cms. 

# Why?
I feel like most headless cms are geared towards marketing teams.
I just want something simple that gets out of my way to serve data for some websites.
I feel like this is it.

# Using
To use centra, you first have to create a folder (/git repo)
This repository will contain all your content files aswell as a dockerfile.

## Setup
This guide takes the following structure for granted
```
- Dockerfile
- content/
-- pages
--- home.yaml
-- sections
--- about.yaml
```

### Dockerfile
To get your CMS running, you have to create the following Dockerfile
```Dockerfile
FROM ghcr.io/cheetahbyte/centra:main
COPY content/ /content
ENV CONTENT_ROOT=/content
```

### Other files
#### `content/pages/home.yaml`
This file contains the main "page". 
```yaml
updated_at: 2025-11-17

sections:
  - $ref: "sections/about/why"
```

#### `content/sections/about.yaml`
```yaml
type: "text-section"
heading: "About US"
body: |
  Bla
  Bla
  Bla
```

> 🎉 Congrats!
> Build & deploy your Dockerfile, and your CMS is live.

# Environment Variables
There a few env vars that can be configured

| Variable | Description | Default |
|----------|-----------|----------|
| `PORT`   | sets the port on which the webserver starts.    | `3000`   |
| `CONTENT_ROOT`  | sets the directory where the server will look for content.  | `/content`   |
