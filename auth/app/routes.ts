import { type RouteConfig, route } from "@react-router/dev/routes";

export default [
  route("auth/:pathname", "./routes/auth/index.tsx"),
  route("account/*", "./routes/account/index.tsx"),
  route("api/auth/*", "../server/api/auth.ts"),
] satisfies RouteConfig;
