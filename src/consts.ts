// Place any global data in this file.
// You can import this data from anywhere in your site by using the `import` keyword.

export const possibleAuthors = {
  nathan: ['nathan', 'nate'],
  natalie: ['natalie', 'nat', 'relue'],
  nathanAndNatalie: ['us', 'nathan & natalie'],
};

export type PossibleAuthors = keyof typeof possibleAuthors;

export type ImageAccessibility = {
  alt: string; // a description of the image
  role?: astroHTML.JSX.AriaRole; // list of image roles: https://developer.mozilla.org/en-US/docs/Web/Accessibility/ARIA/Roles#roles_defined_on_mdn
  ariaDescribedby?: string; // if you describe the image in an HTML element, use give it an it like id="carpark-description". that way the screen reader can say "this div describes the image"
  loading?: astroHTML.JSX.ImgHTMLAttributes["loading"]; // set to "eager" if image is essential to the post, "lazy" if it is not. default of this is lazy.
};

// Inside of this object, the key is going to be the same name you gave it inside the blog post.
// Then, make an object conforming to the type definition above. The only required key is alt.

export const CoversAccessibility: { [x: string]: ImageAccessibility; } = {
  "abstractions.png": { alt: "left: five plastic spiders labelled as concrete, middle: five written in tally marks labelled as represntational, right: the number five written arabic numeral, labelled as abstract" },
  "dogcar.jpg": { alt: "a google street view image of a dog seemingly in the driver's seat of a moving car without a human" },
  "excitementometer.jpg": { alt: "a gauge of excitement, towards high" },
  "LainLaugh.gif": { alt: "" },
  "linuxgraph.png": { alt: "a graph showing the sales and activations by region for linux" },
  "ncmpcpp.png": { alt: "a terminal window with a music playing program open, complete with song picker and audio visualizer" },
  "newsboat.png": { alt: "a terminal window displaying a RSS feed entry" },
  "Spectral.png": { alt: "The Adams spectral sequence for p=3, t-s <= 45" },
  "torus.png": { alt: "a cifford torus" },

};
export const StandaloneImgAccessibility: { [x: string]: ImageAccessibility; } = {
  "carpark.png": { alt: "site's background image of a girl standing in a carpark", role: "presentation" },
  "cyber.jpg": { alt: "a nighttime cityscape" },
};