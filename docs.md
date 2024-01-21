# Overview

---

To make a new post:

- Make a new document in `src/blogs/`
- This document can be `.astro` (html) or `.mdx` (markdown)
- If `.astro`, then export details as such:
- After you're finished, run:

```bash
git commit -am "author: new post"
git push
```

This will push it to prod.

Blogs will be at `.com/author/p/id`.

If you want to verify that there are no errors, run `npm run build`, fix any errors, and
then push. There are going to be a bunch of warns, don't worry about that. If it fails,
try running `npm install` so add any new dependencies, otherwise a nice error message
should appear.

## Things that matter:

-`date`: make it `YYYY-MMM-DD`, or at least a 3 letter month so it's unambiguous to
javascript.  
-`author`: make it something as defined in `consts.ts` -`image`: can be a string of the
filename, or an array of strings for multiple images. place the images in
`/public/images/covers/[year]/` -`aria`: an object defined in `types.ts`, mainly used for
the `alt` key for alt text on images.

To use aria in `.astro`, it looks like this

```js
export const details = {
  title: "i got a cat",
  date: "1970-Jan-01T00:00",
  author: "nathan",
  image: "cat.jpg",
  aria: {
    "cat.jpg": {
      alt: "a really cool kitty cat",
    },
  },
};
```

If used in `.mdx`, it looks like this NOTE: if you have a colon in the string, you need to
put quotes around it NOTE: don't put dashes behind the keys

```yml
title: i got a cat
date: "1970-Jan-01T00:00"
author: nathan
image: cat.jpg
aria:
  cat.jpg:
    alt: a really cool kitty cat
```
