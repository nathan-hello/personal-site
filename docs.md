# Overview

---

To make a new post:

-   Make a new document in `src/blogs/`
-   This document can be `.astro` (html) or `.mdx` (markdown)
-   If `.astro`, then export details as such:

```ts
---
export const details = {
    title: 'New Post',
    date: '2024-jan-04',
    author: 'natalie',
}
---
<b>html goes here</b>
```

If `.mdx`, then export details like this:

```
---
title: New Post
date: 04-Jan-2024
author: natalie
---
Markdown goes here
```

> Note the triple-dashes in either file. They are necessary.

-   By default, posts will be made with the default layout.

    -   For Natalie, that is in `src/layouts/natalie/Root.astro`

-   After you're finished, run:

```bash
git add .
git commit -m "new post"
git push
```

This will push it to prod.

Blogs will be at `.com/author/yyyy/mm/dd/title`.

If you want to verify that there are no errors, run `npm run build`, fix any errors, and then push. There are going to be a bunch of warns, don't worry about that.

## Things that matter:

-   The `date` in the details has to be able to be parsed by Javascript. It doesn't really matter what the format is as long as it's unambiguous, so I've elected for YYYY-MMM-DD. If javascript can't figure out the date, then it will throw an error and the build will fail.
-   The author name has to be one of those listed in `consts.ts`. This is a build-time check so it will make sure that each post knows where to go. This was done so we don't have to rely on filesystem locations (more on that later). If an author name is used that isn't in `possibleAuthors`, then the build will fail.

## What is Astro/Mdx?

Astro is a JS library with the schtick of having zero client-side javascript (unless given a <\script> tag). It does this by compiling all of your .astro files into html.

Astro files have two sides: Javascript and html.

```astro
[label src/index.astro]

export const text = "here"
<div>This is html land, and here is that text from above: {text}</div>
```

The Javascript, using it the way that we do, is more like a helper for creating html. The routes we're making are in "static" mode, which means that the javascript doesn't even run on the server, it only runs during the compile step.

This means that after build, the sites will not change. For example:

```astro
[label src/index.astro]

const randomNumber = Math.random()
<div>The random number: {randomNumber}</div>
```

This random number will be generated once on build, and then after that it will be 'cached' in the html, never to be rendered again (until a subsequent build).

### Mdx

Full docs: [mdx-docs](https://mdxjs.com/docs/)

Mdx is a superset of markdown, for our usecase it's cool for two reasons:

-   We can statically analyze the title/date/author/etc so we can guarantee that everything is in the right shape
-   Components can be imported into them just the same as astro

I think we'll find it often that markdown is preferable over html just because it's quicker to write.

```mdx
## [label src/blog/nathan/asdf.mdx]

title: asdf
author: nathan
date: Jan-06-2024

---

import Latex from '@components/natalie/Latex.astro'

_Valid Markdown!_

<Latex formula="E=mc^2" />

_More valid markdown_
```

Note that the import statement is under the triple dashes for this one. Not sure why, it's kinda weird.

### Components & Layouts

The big power of Astro is the composability of these without client-side js. You can import a file like so:

```astro
[label src/index.astro]

import Component from "@components/default/CenterText.astro"
<Component center={true}>
    <div>This div will go in the slot tag</div>
</Component>
```

Layouts work the same exact way, they're just in a different folder by convention.

When you're making a component or layout that expects to have children, you use a slot tag

```astro
[label src/components/default/CenterText.astro]

interface Props { center: boolean } const c = center ? "text-center" : "text-left"
<div class={c}>
    <slot />
</div>
```

#### Layout for Blog post

Again, you have creative freedom over what it looks like. The blog post's body (html or mdx) will be in the `<slot/>`. But the rest of the details object you need to have access to. Here's how you do that:

```astro
---
import Meta from '@components/default/Meta.astro'
import Header from '@components/default/Header.astro'
import Footer from '@components/default/Footer.astro'
import type { BlogDetails } from '@/types'
interface Props {
    details: BlogDetails
}
const { details } = Astro.props
---

<Meta title={details.title} description={details.description} image={details.image} />
<Header />

<h1>{details.author}</h1>

<main class="min-h-screen bg-pink-800 max-w-3xl">
    <slot />
</main>
```

If you're making a custom layout, and you want to use it for a post, you use it like this:

```astro
---
import Layout from '@layouts/natalie/CustomLayout.astro'
export const details = {
    title: 'fdsa',
    author: 'natalie',
    date: 'Jan-06-2024',
}
---

<Layout details={details}>
    <h1>Content Here</h1>
</Layout>
```

## For Natalie

I made a Latex component already.

```astro
[label src/blog/natalie/asdf.astro]

import Latex from "@components/natalie/Latex.astro"
<Latex formula="E=mc^2" />
```

By default, this element will have a class of name "latex-default", which means you can style it like so:

```html
[label src/blog/natalie/asdf.astro]
<style is:global>
    .latex-default {
        text-align: center;
        background-color: cyan;
    }
</style>
```

Note the `is:global` directive. This is the first Astro-specific thing mentioned here. It just means that you want to affect the styling of components outside of my scope.

This means that all css you define is scoped to that component or page. You can safely do something like

```html
[label src/blog/natalie/asdf.astro]
<style>
    h1 {
        font-weight: 700;
    }
</style>
```

And because the CSS is scoped to that page or component, only the `h1` tags within that component will change. To be tested: What about children? I'm not sure.

**# Additional Options**

The details object, defined in `types.ts`, has a few optional configurations for special usecases.

| Key            | Required? | Description                                                                                                                                                                                                                                |
| -------------- | --------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| title          | X         | Defines part of the URL                                                                                                                                                                                                                    |
| date           | X         | Defines part of the URL                                                                                                                                                                                                                    |
| author         | X         | Defines part of the URL                                                                                                                                                                                                                    |
| tags           |           | However you want to use tags, we can figure it out. Actually using them is not implemented yet                                                                                                                                             |
| overrideHref   |           | If you don't want the url to be `YYYY/MM/DD/title`, you can use this to make a custom route that lives in, for example `/nathan/myOverride`                                                                                                |
| overrideLayout |           | Because `.astro` is just straight up html, and `.mdx` also supports importing components from Astro, you can make a route that does not inherit from the default layout. This means that the page is a blank canvas, do whatever you want. |
| description    |           | Used for <meta> tags, for sharing on social media                                                                                                                                                                                          |
| image          |           | Used for <meta> tags, for sharing on social media                                                                                                                                                                                          |
