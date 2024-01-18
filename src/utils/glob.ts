import type { PossibleAuthors } from "@/consts";
import { type BlogAstro, type BlogMdx, type BlogDetails, type Post, type Image } from "@/types";
import { extractMetadata, formatBytes, generateHref, parseAuthorName, } from "./parseMetadata";
import type { AstroComponentFactory } from "astro/runtime/server/index.js";
import fs from "fs"
import path from "path"

export const globImages = async (imgs: string[]): Promise<Image[]> => {
  const globber = import.meta.glob("/public/**/*.{jpg,gif,png,jpeg,bmp,webp}", { as: "url" });

  console.log("globber", globber)
  let images: Image[] = []

  for (const img of imgs) {
    console.log("given url", img)
    const i = `/public/images/${img}`
    console.log("i", i)
    const url = await globber[i]()

    const fsPath = `.${i}`
    const size = fs.statSync(fsPath).size
    const ext = path.extname(fsPath)
    const file = path.basename(fsPath, path.extname(fsPath));

    // const href = i.substring(6, undefined)

    if (!url) {
      throw new Error(`ERROR: ${url} undefined from ${imgs}`)
    }


    images.push({
      size: formatBytes(size),
      ext: ext,
      url: url,
      filename: file,
      fullname: `${file}${ext}`
    })

  }

  return images;
}

export type RGlobBlogs = {
  params: {
    post: string
  };
  props: {
    c: Post;
  }
}

export async function globBlogs(
  limit: number | undefined,
  author: PossibleAuthors | undefined,
  hideHidden: boolean | undefined
): Promise<RGlobBlogs[]> {

  let combined: Post[] = [];
  const interim: { details: BlogDetails, component: AstroComponentFactory; }[] = [];
  const blogs = import.meta.glob<BlogAstro>("/src/blog/**/*.astro");

  let count = 100001;

  for (const post in blogs) {
    const f = await blogs[post]();
    const g = extractMetadata(f);
    interim.push(g);
    count++;
    const newCount = count + 1; count++
  }

  const mdxs = import.meta.glob<BlogMdx>("/src/blog/**/*.mdx");

  for (const post in mdxs) {
    const f = await mdxs[post]();
    const g = extractMetadata(f);
    interim.push(g);
    count++;
    const newCount = count + 1; count++
  }

  for (const p of interim) {
    const id = count;
    let href;
    if (p.details.overrideHref) {
      href = p.details.overrideHref;
    } else {
      href = generateHref(p.details.date, id);
    }

    let imgs: Image[] = []
    if (typeof p.details.image === "string") {
      imgs = await globImages([p.details.image])
    } 

    if (Array.isArray(p.details.image)) {
      imgs = await globImages(p.details.image)
    }

    combined.push({
      ...p.details,
      id: id,
      Component: p.component,
      dateObj: new Date(p.details.date),
      relativeUrl: href,
      absoluteUrl: `${p.details.author}/${href}`,
      globbedImgs: imgs
    });
  }

  combined = combined.sort((a, b) => b.dateObj.getTime() - a.dateObj.getTime());

  if (author) {
    combined = combined.filter((c) => parseAuthorName(c.author) === author);
  }
  if (limit) {
    combined = combined.slice(0, limit);
  }
  if (hideHidden) {
    combined = combined.filter((c) => c.hidden !== true)
  }


  return combined.map((c) => {
    console.log(c);
    return {
      params: { post: c.relativeUrl },
      props: { c, }
    };
  });
}
