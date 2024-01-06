import type { PossibleAuthors } from "@/consts";
import { possibleAuthors } from "@/consts";
import type { BlogAstro, BlogDetails, BlogMdx } from "@/types";
import type { AstroComponentFactory } from "astro/runtime/server/index.js";

export function parseAuthorName(s: string, file: string): PossibleAuthors {

  if (typeof s !== "string") { throw Error("Author name at " + file + " is  not string!"); }

  for (const [key, names] of Object.entries(possibleAuthors)) {
    if (names.includes(s.toLowerCase())) {
      return key as PossibleAuthors;
    }
  }

  throw Error("Author of name: " + s + " not found");

}

export function parseDateString(s: string): { dateStr: string, dateObj: Date; } {
  const d = new Date(s);

  if (!d) { throw Error("Date: " + s + " not valid"); }

  const year = d.getFullYear();
  const month = (d.getMonth() + 1).toString().padStart(2, "0");
  const day = d.getDate().toString().padStart(2, "0");


  return {
    dateStr: year + '/' + month + '/' + day,
    dateObj: d
  };

}

export function extractMetadata(i: BlogAstro | BlogMdx): [BlogDetails, AstroComponentFactory] {
  if ("details" in i) {
    return [i.details, i.default];
  }
  if ("frontmatter" in i) {
    return [i.frontmatter, i.Content];
  }

  throw new Error("Input: " + i + " is not a valid BlogAstro or BlogMdx");
}

export function shapeForRendering([details, component]: [BlogDetails, AstroComponentFactory]) {

  const { dateObj, dateStr } = parseDateString(details.date);
  let href;
  if (details.overrideHref) {
    href = details.overrideHref;
  } else {
    href = dateStr + '/' + encodeURIComponent(details.title);
  }

  return {
    params: { slug: href },
    props: {
      Component: component,
      details: details,
    }
  };
}