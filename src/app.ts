import {Hono} from "hono"
import { cors } from "hono/cors"
import {getCollection, getEntry} from "./content-loader"

const app = new Hono()

app.use("/api/*", cors())

app.get("/health", (c) => c.json({ status: "ok" }));

app.get("/api/:collection", async (c) => {
  const collection = c.req.param("collection");
  const items = await getCollection(collection);

  return c.json({ collection, items });
});

app.get("/api/:collection/:slug", async (c) => {
  const collection = c.req.param("collection");
  const slug = c.req.param("slug");

  const item = await getEntry(collection, slug);

  if (!item) {
    return c.json(
      { error: "Not found", collection, slug },
      404,
    );
  }

  return c.json(item);
});

export default {
  port: 3000,
  fetch: app.fetch,
};
