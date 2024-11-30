import { defineConfig } from "astro/config";
import vercel from "@astrojs/vercel/serverless";
import mdx from "@astrojs/mdx";
import tailwind from "@astrojs/tailwind";
import rehypeKatex from "rehype-katex";
import remarkMath from "remark-math";
import db from "@astrojs/db";

// https://astro.build/config
export default defineConfig({
  site: "https://reluekiss.com",
  output: "hybrid",
  build: {
    inlineStylesheets: "always",
    assets: "a",
  },
  adapter: vercel(),
  integrations: [
    tailwind(),
    mdx({
      remarkPlugins: [remarkMath],
      rehypePlugins: [rehypeKatex],
    }),
    db(),
  ],
  server: { port: 3000, host: "127.0.0.1" },
});
