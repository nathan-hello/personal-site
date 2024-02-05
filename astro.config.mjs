import { defineConfig } from "astro/config";
import vercel from "@astrojs/vercel/serverless";
import mdx from "@astrojs/mdx";
import tailwind from "@astrojs/tailwind";
import MDXCodeBlocks, { mdxCodeBlockAutoImport } from "astro-mdx-code-blocks";
import AutoImport from "astro-auto-import";
import rehypeKatex from "rehype-katex";
import remarkMath from "remark-math";

// https://astro.build/config
export default defineConfig({
  site: "https://reluekiss.com",
  output: "hybrid",
  build: {
    inlineStylesheets: "always",
    assets: "a",
  },
  adapter: vercel({
    webAnalytics: true,
  }),
  vite: {
    define: {
      'import.meta.env.PUBLIC_VERCEL_ANALYTICS_ID': JSON.stringify(process.env.VERCEL_ANALYTICS_ID),
    },
  },
  integrations: [
    tailwind(),
    AutoImport({
      imports: [mdxCodeBlockAutoImport("src/components/Code.astro")],
    }),
    MDXCodeBlocks(),
    mdx({
      remarkPlugins: [remarkMath],
      rehypePlugins: [rehypeKatex],
    }),
  ],
  server: { port: 3000, host: "127.0.0.1" },
});
