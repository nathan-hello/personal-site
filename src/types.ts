import type { AstroInstance, MarkdownInstance } from "astro";

type BlogDetails = {
  title: string;
  img?: string;
};

export interface Frontmatter {
  title: string;
  img?: string;
}

export type BlogPage = AstroInstance & {
  details: BlogDetails;
};
