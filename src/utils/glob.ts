import type { PossibleAuthors } from "@/consts";
import { type BlogAstro, type BlogMdx, type BlogDetails, type Post } from "@/types";
import { extractMetadata, generateHref, } from "./parseMetadata";
import type { AstroComponentFactory } from "astro/runtime/server/index.js";

export async function globImage(path: string) {

}

export async function globBlogs(limit: number, filterAuthor: PossibleAuthors) {

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

    combined = combined
        .filter((c) => c.author === filterAuthor)
        .sort((a, b) => b.dateObj.getTime() - a.dateObj.getTime())
        .slice(0, limit);

    return combined.map((c) => {
        return {
            params: { slug: `${c.author}/${c.href}` },
            props: { c }
        };
    });

}