# Content Model

This document explains how **Centra** processes, categorizes, and serves your content.

---

## Caching

To achieve high performance, Centra selectively caches content depending on file type.

Structured content such as `.yaml` and `.md` files is **fully loaded at startup**, parsed, converted to JSON, and stored internally as JSON bytes. This allows Centra to return these entries extremely fast without touching the filesystem again.

Binary assets (for example images or other media files) are handled differently. For these files, **only metadata is cached**, while the actual file contents are read from disk on demand. This keeps memory usage low and ensures the cache remains lightweight.

More advanced caching strategies (such as LRU-based eviction) may be introduced in the future, but the current approach provides a good balance between performance and simplicity.

---

## First-Class Content

Markdown (`.md`) and YAML (`.yaml`) files are treated as **first-class content** in Centra.

This means they can expose structured metadata, which Centra makes available:

* on **collection level**
* during **filtering**
* and when accessing individual entries

> [!WARNING]
> Metadata is never stripped from responses. If a file defines metadata, it will always be included when the file is requested.

---

### YAML Files

```yaml
author: Leo
category: baking
---
ingredients:
  - name: Sugar
```

Everything **above the fence (`---`)** is treated as metadata.
Everything **below** is considered the content body.

---

### Markdown Files

```md
---
author: Leo
category: baking
---
# Recipe
Firstly, scramble your eggs.
```

For Markdown files, all YAML frontmatter is automatically extracted as metadata and exposed at the collection level.

---

## Collections

In Centra, **collections are simply folders** inside your content directory.

Given the following structure:

```
- content
  - recipes
    - scrambled_eggs.yaml
  - blogs
    - first.md
    - second.md
```

Centra will automatically create two collections:

* `recipes`
* `blogs`

---

### Files in a Collection

**first.md**

```yaml
---
author: leo
title: This is the first blog post
state: released
---
# This is my stuff.
```

**second.md**

```yaml
---
author: maik
title: My deep-dive into monitors
state: not-released
---
# This is my stuff.
```

Requesting the `blogs` collection:

```sh
curl http://localhost:3000/api/blogs
```

Returns:

```json
{
  "collection": "blogs",
  "items": [
    {
      "slug": "blogs/first",
      "meta": {
        "author": "leo",
        "contentType": "application/json",
        "kind": "markdown",
        "size": 104,
        "state": "released",
        "title": "This is the first blog post"
      }
    },
    {
      "slug": "blogs/second",
      "meta": {
        "author": "maik",
        "contentType": "application/json",
        "kind": "markdown",
        "size": 108,
        "state": "not-released",
        "title": "My deep-dive into monitors"
      }
    }
  ]
}
```

---

### What Are Collections Used For?

Collections are ideal for building **overviews** of related content.
Typical use cases include:

* blog listings
* recipe indexes
* documentation pages
* filtered or paginated content views

---

## Everything Else (Raw Files)

All other file types are treated as **raw assets** and streamed directly to the client.

These files are **not fully cached**. Only limited metadata is stored, and the file contents are read from disk when requested.

---

### Example: Image Collection

```sh
curl http://localhost:3000/api/images
```

```json
{
  "collection": "images",
  "items": [
    {
      "slug": "images/one",
      "meta": {
        "contentType": "image/jpeg",
        "ext": ".jpg",
        "kind": "binary_ref",
        "mtime": 1767700524,
        "size": 64480
      }
    }
  ]
}
```

The `kind: "binary_ref"` field indicates that this file is **not cached in memory** and will be streamed from disk when accessed.
