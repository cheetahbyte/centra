# Variables 
Centra includes a few quality of life variables which makes it easier to do certain things.

## `$rel(<path>)` — Relative Files

`$rel()` lets you reference files relative to the current content file without hard-coding absolute URLs.
Centra resolves this variable **server-side** before sending content to the client.

### Example

Without `$rel` (hard-coded URL):

```md
![Banana Bread](https://cms.centra.test/api/images/banana-bread.jpg)
```

With `$rel`:

```md
![Banana Bread]($rel(images/banana-bread.jpg))
```

Centra will replace this with the correct absolute URL automatically.

### Requirement

For `$rel()` to work, **`CENTRA_PUBLIC_URL` must be set**:

```env
CENTRA_PUBLIC_URL=https://cms.centra.example
```

This value is used as the base URL when resolving relative paths.

### Notes

* Paths are resolved **relative to CONTENT_ROOT directory**
* Works only in YAML and Markdown files
