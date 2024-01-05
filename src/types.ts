import type { AstroInstance } from "astro";


export type BlogDetails = {
  title: string;
  date: string;
  author: string;
  overrideHref?: string;
  overrideLayout?: boolean;
  description?: string;
  image?: string;
  tags?: string[];
};

export interface Frontmatter {
  title: string;
  date: string;
  author: string;
  overrideHref?: string;
  description?: string;
  image?: string;
  tags?: string[];
}

export type BlogPage = AstroInstance & {
  details: BlogDetails;
};
