import type { PossibleAuthors } from "@/consts";
import { possibleAuthors } from "@/consts";
import type { BlogAstro, BlogDetails, BlogMdx } from "@/types";
import type { AstroComponentFactory } from "astro/runtime/server/index.js";

export function formatBytes(bytes: number, decimals = 2) {
  if (bytes === 0) return '0 Bytes';

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

export function parseAuthorName(s: string): PossibleAuthors {

  for (const [key, names] of Object.entries(possibleAuthors)) {
    if (names.includes(s.toLowerCase())) {
      return key as PossibleAuthors;
    }
  }

  throw Error(`Author of name: ${s} not found`);
}

function parseDateString(s: string): string {
  const d = new Date(s);

  if (!d) { throw Error(`Date: ${s} not valid`); }

  const year = d.getFullYear();
  const month = (d.getMonth() + 1).toString().padStart(2, "0");
  const day = d.getDate().toString().padStart(2, "0");

  return `${year}/${month}/${day}`;
}

export function extractMetadata(i: BlogAstro | BlogMdx): { details: BlogDetails, component: AstroComponentFactory; } {
  if ("details" in i) {
    return { details: i.details, component: i.default };
  }
  if ("frontmatter" in i) {
    return { details: i.frontmatter, component: i.Content };
  }

  throw new Error(`Input: ${i} is not a valid BlogAstro or BlogMdx`);
}

export function generateHref(dateStr: string, i: number) {
  return `${parseDateString(dateStr)}/${(i).toString()}`;
}

