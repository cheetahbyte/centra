import {readdir, readFile, stat} from "node:fs/promises"
import path from "node:path"
import {parse} from "yaml"

export type ContentRecord = {
  slug: string;
  [key: string]: unknown;
}

const CONTENT_ROOT = path.join(process.cwd(), "content");

export async function getCollection(
  collection: string
): Promise<ContentRecord[]> {
    const dir = path.join(CONTENT_ROOT, collection);
    let files: string[] = []
    try {
      files = await readdir(dir)
    } catch {
      return []
    }

    const entries: ContentRecord[] = []

    for (const file of files) {
      if (!file.endsWith(".yaml") && !file.endsWith(".yml")) continue
      const filePath = path.join(dir, file)

      const stats = await stat(filePath)
      if (!stats.isFile()) continue;

      const raw = await readFile(filePath, "utf8")
      const data = parse(raw) ?? {}

      let slug = (data.slug as string) ?? path.basename(file, ".yaml")
      slug = slug.toString()

      entries.push({
        slug,
        ...data
      })
    }
    return entries
}

export async function getEntry(
  collection: string,
  slug: string,
): Promise<ContentRecord | null> {
  const dir = path.join(CONTENT_ROOT, collection)

  // possible files
  const possibles = [`${slug}.yaml`, `${slug}.yml`]

  for (const filename of possibles) {
    const filepath = path.join(dir, filename)
    try {
      const raw = await readFile(filepath, "utf8")
      const data = parse(raw) ?? {}
      return {
        slug,...data
      }
    } catch {}
  }
  return null
}
