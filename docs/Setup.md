# Setup

I recommend that you also take a look at the [Configuration Page](./Configuration.md)

## Setup with Docker Compose
Example Docker Compose:
```yaml
services:
  centra:
    image: ghcr.io/cheetahbyte/centra:latest
    volumes:
      - "./content:/content"
    ports:
      - "3000:3000"
```

## Baking Centra into your own Dockerfile
If you don't want to use Github Webhooks, you can also just use Dockerfiles. This setup could be useful, if you dont want to expose Centra to the internet and only talk to it locally.

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
FROM ghcr.io/cheetahbyte/centra:latest
COPY content/ /content
ENV CONTENT_ROOT=/content
```

### Other files
#### `content/pages/home.yaml`
This file contains the main "page". 
```yaml
quote: Nobody is perfect
```

#### `content/sections/about.yaml`
```yaml
team:
  - Leo
  - Maik
  - Kevin
  - Laura
  - Luca
  - Enrico Gieren Jacob
```

> 🎉 Congrats!
> Build & deploy your Dockerfile, and your CMS is live.
## Setup with Kubernetes
-- coming soon --

## Setup with Git Sync
The setup with Git Sync is not as trivial as the normal docker setups, but not hard either.

1. Setup a Centra instance and make it reachable from the outside (domain or ip)
2. Go to your repositories settings and add a webhook. 
  - Input the url of your centra instance, dont forget to add the path `/webhook`. For example: `https://centra.test/webhook`.
  - Set the content type to `application/json`. This is really important.
  - If you have configured SSL for your domain, please check the SSL verification setting.
  - Leave the events on `Just the push event`.
  - (Optionally) add a webhook secret for added security. Make sure to also expose it to centra via the `WEBHOOK_SECRET` env var.
3. Setup your deploy keys

> [!IMPORTANT]
> Centra only supports `ed25519` keys and I can't be bothered to implement others.

You have three options when using deploy keys:
a. Bring your own. You can either put them in a directory of your choice, mount it to `/keys` (or provide a different path via `KEYS_DIR`)
a,5/b. Bring your own. You can provide your keys in plain text via `SSH_PUBLIC_KEY` and `SSH_PRIVATE_KEY`.
c. Let centra do the work. Centra will automatically create the keys. I recommend to mount a volume to `/keys` to keep the keys after a restart.

After you told centra about your ssh key, you can add them to your repository as deploy key.

> Centra will always output the ssh public key (which is the one you need to add as deploy key) to the console, so you don't have to take extra steps. (especially helpful when letting centra generate them.)
4. Enjoy! Centra should now automatically update the cache when you push changes to your repository. (main branch only)
