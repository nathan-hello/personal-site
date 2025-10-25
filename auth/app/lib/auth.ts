import { createAuthClient } from "better-auth/react";
import {
  passkeyClient,
  usernameClient,
  twoFactorClient,
  oneTimeTokenClient,
  magicLinkClient,
  emailOTPClient,
} from "better-auth/client/plugins";

const url =
  process.env.NODE_ENV === "development"
    ? "http://localhost:5173"
    : process.env.PRODUCTION_URL;
if (!url) {
  throw Error("process.env.PRODUCTION_URL was undefined");
}

export const authClient = createAuthClient({
  baseURL: url,
  plugins: [
    passkeyClient(),
    usernameClient(),
    emailOTPClient(),
    magicLinkClient(),
    twoFactorClient(),
    oneTimeTokenClient(),
  ],
});
