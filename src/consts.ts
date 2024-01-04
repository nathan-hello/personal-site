// Place any global data in this file.
// You can import this data from anywhere in your site by using the `import` keyword.

export const SITE_TITLE = "Nathan & Natalie";
export const SITE_DESCRIPTION = "Welcome to my website!";

export const possibleAuthors = {
  nathan: [
    "nathan", "nate"
  ],
  natalie: [
    "natalie", "nat", "relue"
  ],
  nathanAndNatalie: [
    "both", "nathan & natalie"
  ]
}

export type PossibleAuthors = keyof typeof possibleAuthors

