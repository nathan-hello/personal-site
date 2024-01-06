import { defineConfig } from "astro/config";
import vercel from "@astrojs/vercel/serverless";
import mdx from "@astrojs/mdx";

import tailwind from "@astrojs/tailwind";

// https://astro.build/config
export default defineConfig({
  site: "https://nathan-hello.com",
  output: "hybrid",
  adapter: vercel({
    webAnalytics: false
  }),
  integrations: [mdx(), tailwind()],
  server: { port: 3000 }
});
