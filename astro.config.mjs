import { defineConfig } from "astro/config";
import vercel from "@astrojs/vercel/serverless";
import mdx from "@astrojs/mdx";
import tailwind from "@astrojs/tailwind";
import MDXCodeBlocks, { mdxCodeBlockAutoImport } from 'astro-mdx-code-blocks';
import AutoImport from "astro-auto-import"

// https://astro.build/config
export default defineConfig({
  site: "https://reluekiss.com",
  output: "hybrid",
  build: {
    inlineStylesheets: "always"
  },
  adapter: vercel({
    webAnalytics: false
  }),
  integrations: [
    tailwind(),
    AutoImport({
      imports: [mdxCodeBlockAutoImport("src/components/Code.astro")]
    }),
    MDXCodeBlocks(),
    mdx(), 
  ],
  server: { port: 3000, host: "127.0.0.1" },
  
});
