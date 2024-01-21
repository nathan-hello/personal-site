// import rss from '@astrojs/rss';
// import { globBlogs } from '@utils/glob';
// import { type RSSFeedItem } from '@astrojs/rss';
//
// export async function GET(context) {
//   const posts = await globBlogs(undefined, undefined, true)
//
//   const rssPostsasdf: RSSFeedItem = {
//     author:"",
//     source: "",
//
//
//   }
//   const rssPosts: RSSFeedItem = posts.map((p) => {return {
//
//   }})
//   return rss({
//     title: 'Nathan & Natalie',
//     description: 'We are boingus.',
//     site: context.site,
//     // Array of `<item>`s in output xml
//     // See "Generating items" section for examples using content collections and glob imports
//     items: posts.map((p) => {
//       return {
//         title: p.props.c.title,
//         pubDate: p.props.c.dateObj,
//         description: "",
//         author: p.props.c.author, // this is supposed to be an email
//         content: p.props.c.Component.
//       }
//     }) ,
//     // (optional) inject custom xml
//     customData: `<language>en-us</language>`,
//   });
// }
// TODO: follow this tutorial
// https://scottwillsey.com/rss-pt2/
