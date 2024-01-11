import type { PossibleAuthors } from "@/consts";
import { type BlogAstro, type BlogMdx, type BlogDetails, type Post, type Image } from "@/types";
import { extractMetadata, generateHref, } from "./parseMetadata";
import type { AstroComponentFactory } from "astro/runtime/server/index.js";
import path from "path";
import fs from "fs";

export const globImage = async (img: string): Promise<Image> => {
    const imgPath = path.join(process.cwd(), "public", img);
    const imgSize = fs.statSync(imgPath).size;
    const extension = path.extname(imgPath);
    const globber = import.meta.glob("/public/**/*.{jpg, gif, png, jpeg}", { as: "url" });
    const href = await globber[img]().then((g) => g);

    return {
        href: href,
        size: imgSize.toString(),
        type: extension
    };
};

export async function globBlogs(limit: number | undefined, author: PossibleAuthors | undefined) {

    let combined: Post[] = [];
    const interim: { details: BlogDetails, component: AstroComponentFactory; }[] = [];
    const blogs = import.meta.glob<BlogAstro>("/src/blog/**/*.astro");

    let count = 100001;

    for (const post in blogs) {
        const f = await blogs[post]();
        const g = extractMetadata(f);
        interim.push(g);
        count++;
    }

    const mdxs = import.meta.glob<BlogMdx>("/src/blog/**/.mdx");

    for (const post in mdxs) {
        const f = await mdxs[post]();
        const g = extractMetadata(f);
        interim.push(g);
        count++;
    }

    for (const p of interim) {
        const id = count;
        let href;
        if (p.details.overrideHref) {
            href = p.details.overrideHref;
        } else {
            href = generateHref(p.details.date, id);
        }

        combined.push({
            ...p.details,
            id: id,
            href: href,
            Component: p.component,
            dateObj: new Date(p.details.date),
        });
    }

    combined = combined.sort((a, b) => b.dateObj.getTime() - a.dateObj.getTime());

    if (author) {
        combined = combined.filter((c) => c.author === author);
    }
    if (limit) {
        combined = combined.slice(0, limit);
    }


    return combined.map((c) => {
        console.log(c);
        return {
            params: { post: `${c.author}/${c.href}` },
            props: { c }
        };
    });

}