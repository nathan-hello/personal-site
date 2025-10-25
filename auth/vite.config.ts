import { reactRouterHonoServer } from "react-router-hono-server/dev";
import { reactRouter } from "@react-router/dev/vite";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";
import tsconfigPaths from "vite-tsconfig-paths";
import path from "path";

export default defineConfig({
  plugins: [reactRouterHonoServer({ runtime: "bun" }), tailwindcss(), reactRouter(), tsconfigPaths()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./app"),
      "@server": path.resolve(__dirname, "./server"),
    },
  },
  build: {
    target: "esnext",
  },
});
