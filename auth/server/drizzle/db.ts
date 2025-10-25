import { drizzle } from "drizzle-orm/bun-sqlite";
import { Database } from "bun:sqlite";
import * as schema from "./schema";

const location = process.env.DB_FILE_NAME;
if (!location) {
        throw Error(".env: `DB_FILE_NAME` is required")
}


const sqlite = new Database(location);
export const db = drizzle(sqlite, { schema });
