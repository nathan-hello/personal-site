import { defineConfig } from "astro/config";
import vercel from "@astrojs/vercel/serverless";
import mdx from "@astrojs/mdx";
import path from "path"
import tailwind from "@astrojs/tailwind";

// https://astro.build/config
export default defineConfig({
  site: "https://nathan-hello.com",
  output: "hybrid",
  build: {
    inlineStylesheets: "always"
  },
  adapter: vercel({
    webAnalytics: false
  }),
  integrations: [mdx(), tailwind()],
  server: { port: 3000 },
  vite: {
    resolve: {
      alias: {
        "@images": path.resolve(__dirname, "public/images"),
        "@videos": path.resolve(__dirname, "public/videos"),
        "@audios": path.resolve(__dirname, "public/audios"),
      }
    }
  }
});
