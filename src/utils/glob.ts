import { type PossibleAuthors } from "@/consts";
import {
        type BlogAstro,
        type BlogMdx,
        type BlogDetails,
        type Post,
        type Image,
} from "@/types";
import { extractMetadata, formatBytes, parseAuthorName } from "./parseMetadata";
import fs from "fs";
import path from "path";

export async function globImages(
        imgs: string[],
        year: string,
        aria: BlogDetails["aria"]
): Promise<Image[]> {
        const globber = import.meta.glob("/public/**/*.{jpg,gif,png,jpeg,bmp,webp,svg}", {
                as: "url",
        });

        let images: Image[] = [];


        for (const img of imgs) {
                const i = `/public/images/covers/${year}/${img}`;
                const url = await globber[i]();

                const fsPath = `.${i}`;
                const size = fs.statSync(fsPath).size;
                const ext = path.extname(fsPath);
                const file = path.basename(fsPath, path.extname(fsPath));
                const urlNoPublic = url.slice("/public".length);

                if (!url || !urlNoPublic) {
                        throw new Error(`ERROR: ${url} undefined from ${imgs}`);
                }

                if (!aria || !aria[img]) {
                        console.log(`\n=====\nNo aria for the image ${img}. Consider adding one.\n=====\n`);
                }
                // else {
                // console.log(`aria for ${img}:\n  ${JSON.stringify(aria[img])}`);
                // }

                const defaultAria = { [img]: { alt: "" } };
                const accessibility = { ...defaultAria, ...aria }[img];

                images.push({
                        size: formatBytes(size),
                        ext: ext,
                        url: url,
                        filename: file,
                        fullname: `${file}${ext}`,
                        accessibility: accessibility,
                });
        }
        return images;
};

export type RGlobBlogs = {
        params: {
                post: string;
        };
        props: {
                c: Post;
        };
};

export async function globBlogs(
        limit: number | undefined,
        author: PossibleAuthors | undefined,
        hideHidden: boolean | undefined
): Promise<RGlobBlogs[]> {
        let combined: Post[] = [];
        const interim: ReturnType<typeof extractMetadata>[] = [];
        const blogs = import.meta.glob<BlogAstro>("/src/blog/**/*.astro");


        for (const post in blogs) {
                const f = await blogs[post]();
                const g = extractMetadata(f);
                interim.push(g);
        }

        const mdxs = import.meta.glob<BlogMdx>("/src/blog/**/*.mdx");

        for (const post in mdxs) {
                const f = await mdxs[post]();
                const g = extractMetadata(f);
                interim.push(g);
        }

        interim.sort((a, b) => {
                return a.dateObj.getTime() - b.dateObj.getTime();
        });

        let count = 100001;
        for (const p of interim) {
                const id = count;
                count++;
                let href;
                if (p.details.overrideHref) {
                        href = p.details.overrideHref;
                } else {
                        href = `p/${id}`;
                }

                let imgs: Image[] = [];
                if (typeof p.details.image === "string") {
                        imgs = await globImages(
                                [p.details.image],
                                p.dateObj.getFullYear().toString(),
                                p.details.aria
                        );
                }

                if (Array.isArray(p.details.image)) {
                        imgs = await globImages(
                                p.details.image,
                                p.dateObj.getFullYear().toString(),
                                p.details.aria
                        );
                }

                combined.push({
                        ...p.details,
                        id: id,
                        Component: p.component,
                        dateObj: p.dateObj,
                        relativeUrl: href,
                        absoluteUrl: `/${p.details.author}/${href}`,
                        globbedImgs: imgs,
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
                combined = combined.filter((c) => c.hidden !== true);
        }

        return combined.map((c) => {
                return {
                        params: { post: c.relativeUrl },
                        props: { c },
                };
        });
}

