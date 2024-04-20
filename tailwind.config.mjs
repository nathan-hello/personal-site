/** @type {import('tailwindcss').Config} */
export default {
  content: ["./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}"],
  theme: {
    extend: {
      colors: {
        "nat-white-c01": "#b3b8c3",
        "nat-black-c01": "#0c090a",
        "nat-pink-c01": "#440e38",
        "nat-green-c01": "#0a241e",

    },
    
      },
      fontFamily: {
        ibmplexserief: ["IBMPlexSerif", "sans-serif"],
        vcrosdmono: ["VCROSDMono", "sans-serif"],
      },
    },
  plugins: [],
  corePlugins: {
    preflight: false,
  },
};
