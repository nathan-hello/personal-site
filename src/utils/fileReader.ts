import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";

// Call with import.meta.url and a relative path from that file to import the contents of a file.
export function readRelative(importUrl: string, relativePath: string) {
  const __dirname = path.dirname(fileURLToPath(importUrl));
  const filePath = path.join(__dirname, relativePath);
  return fs.readFileSync(filePath, "utf8");
}
