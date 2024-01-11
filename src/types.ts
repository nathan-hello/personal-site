import type { AstroInstance, MDXInstance } from "astro";
import type { AstroComponentFactory } from "astro/runtime/server/index.js";


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

export type Image = {
  href: string,
  size: string,
  type: string,
};

export type Post = BlogDetails & {
  id: number,
  images?: Image[];
  href: string;
  dateObj: Date;
  Component: AstroComponentFactory;
};


export type BlogAstro = AstroInstance & {
  details: BlogDetails;
};

export type BlogMdx = MDXInstance<BlogDetails>;
