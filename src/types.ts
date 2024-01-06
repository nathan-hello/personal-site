import type { AstroInstance, MDXInstance } from "astro";


export type BlogDetails = {
  title: string;
  date: string;
  author: string;
  overrideHref?: string;
  overrideLayout?: boolean;
  description?: string;
  image?: string;
  tags?: string[];
  hidden?: boolean;
};


export type BlogAstro = AstroInstance & {
  details: BlogDetails;
};

export type BlogMdx = MDXInstance<BlogDetails>;