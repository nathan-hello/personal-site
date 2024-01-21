// Place any global data in this file.
// You can import this data from anywhere in your site by using the `import` keyword.

import type { ImageAccessibility } from "./types";

export const possibleAuthors = {
  nathan: ["nathan", "nate"],
  natalie: ["natalie", "nat", "relue"],
  nathanAndNatalie: ["us", "nathan & natalie"],
};

export type PossibleAuthors = keyof typeof possibleAuthors;

export const StandaloneImgAccessibility: { [x: string]: ImageAccessibility } = {
  "carpark.png": {
    alt: "site's background image of a girl standing in a carpark",
    role: "presentation",
  },
  "cyber.jpg": { alt: "a nighttime cityscape" },
};
