import type { AstroInstance, MDXInstance } from 'astro';
import type { AstroComponentFactory } from 'astro/runtime/server/index.js';

export type ImageAccessibility = {
    alt: string; // a description of the image
    role?: astroHTML.JSX.AriaRole; // list of image roles: https://developer.mozilla.org/en-US/docs/Web/Accessibility/ARIA/Roles#roles_defined_on_mdn
    ariaDescribedby?: string; // if you describe the image in an HTML element, use give it an it like id="carpark-description". that way the screen reader can say "this div describes the image"
    loading?: astroHTML.JSX.ImgHTMLAttributes["loading"]; // set to "eager" if image is essential to the post, "lazy" if it is not. default of this is lazy.
};

export type Image = {
    url: string;
    size: string;
    ext: string;
    filename: string;
    fullname: string;
    accessibility: ImageAccessibility;
};

export type BlogDetails = {
    title: string;
    date: string;
    author: string;
    overrideHref?: string;
    overrideLayout?: boolean;
    description?: string;
    image?: string | string[];
    tags?: string[];
    hidden?: boolean;
    aria?: { [x: string]: ImageAccessibility; };
};

export type BlogAstro = AstroInstance & {
    details: BlogDetails;
};

export type BlogMdx = MDXInstance<BlogDetails>;

export type Post = BlogDetails & {
    id: number;
    globbedImgs: Image[];
    relativeUrl: string;
    absoluteUrl: string;
    dateObj: Date;
    Component: AstroComponentFactory;
};