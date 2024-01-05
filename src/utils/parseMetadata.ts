import type { PossibleAuthors } from "@/consts";
import { possibleAuthors } from "@/consts";
import type { BlogPage, Frontmatter } from "@/types";
import type { AstroInstance, MDXInstance } from "astro";
import { z } from "zod";

export function parseAuthorName(s: string, file: string): PossibleAuthors {

  if (typeof s !== "string") { throw Error("Author name at " + file + " is  not string!"); }

  for (const [key, names] of Object.entries(possibleAuthors)) {
    if (names.includes(s.toLowerCase())) {
      console.log(key);
      return key as PossibleAuthors;
    }
  }

  throw Error("Author of name: " + s + " not found");

}

export function parseDateString(s: string, file: string): string {
  const d = new Date(s);

  if (!d) { throw Error("Date from file " + file + " not valid"); }

  return d.toLocaleDateString("eu", { month: '2-digit', day: '2-digit', year: 'numeric', });

}

export function parseAstroPostMetadata(i: BlogPage) {
  const dateStr = parseDateString(i.details.date, i.file);
  let href;
  if (i.details.overrideHref) {
    href = i.details.overrideHref;
  } else {
    href = dateStr + '/' + encodeURIComponent(i.details.title);
  }

  return {
    params: { slug: href },
    props: {
      type: "blog-astro",
      Component: i.default,
      details: i.details
    }
  };
}

export function parseMdxPostMetadata(m: MDXInstance<Frontmatter>) {
  const dateStr = parseDateString(m.frontmatter.date, m.file);

  let href;
  if (m.frontmatter.overrideHref) {
    href = m.frontmatter.overrideHref;
  } else {
    href = dateStr + '/' + encodeURIComponent(m.frontmatter.title);
  }

  return {
    params: { slug: href },
    props: {
      type: 'blog-markdown',
      Component: m.Content,
      details: m.frontmatter,
    },
  };

}