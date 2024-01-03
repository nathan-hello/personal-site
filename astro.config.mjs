import { defineConfig } from 'astro/config';
import mdx from '@astrojs/mdx';
import sitemap from '@astrojs/sitemap';

import tailwind from "@astrojs/tailwind";

// https://astro.build/config
export default defineConfig({
  integrations: [
    mdx(), 
    sitemap(), 
    tailwind()
  ],
  server: {
    host: true,
    port: 3000
  },
  output: "static",
  site: "http://localhost:3000/"
});
