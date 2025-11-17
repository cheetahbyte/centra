import {Hono} from "hono"
import { cors } from "hono/cors"
import { logger } from 'hono/logger'
import { requestId } from 'hono/request-id'
import { timeout } from 'hono/timeout'
import {getCollection, getEntry} from "./content-loader"

const app = new Hono()

app.use("/api/*", cors())
app.use("/api/*", timeout(3000))
app.use(logger())
app.use(requestId())

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
