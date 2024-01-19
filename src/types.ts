import type { AstroInstance, MDXInstance } from 'astro'
import type { AstroComponentFactory } from 'astro/runtime/server/index.js'
import type { ImageAccessibility } from './consts'

export type BlogDetails = {
    title: string
    date: string
    author: string
    overrideHref?: string
    overrideLayout?: boolean
    description?: string
    image?: string | string[]
    tags?: string[]
    hidden?: boolean
}

export type Image = {
    url: string
    size: string
    ext: string
    filename: string
    fullname: string
    accessibility: ImageAccessibility
}

export type Post = BlogDetails & {
    id: number
    globbedImgs: Image[]
    relativeUrl: string
    absoluteUrl: string
    dateObj: Date
    Component: AstroComponentFactory
}

export type BlogAstro = AstroInstance & {
    details: BlogDetails
}

export type BlogMdx = MDXInstance<BlogDetails>
