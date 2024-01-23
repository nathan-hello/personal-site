import { createClient } from "@supabase/supabase-js";
import "dotenv/config";
import type { Database } from "./supabase.types";
import type { Post } from "@/types";

if (!process.env.DB_URL || !process.env.DB_KEY) {
   throw new Error((`No URL or key provided: ${JSON.stringify({DB_URL: process.env.DB_URL, DB_KEY: process.env.key})}`))

}
const s = createClient<Database>(process.env.DB_URL, process.env.DB_KEY);

export async function selectBlogById(id: number) {
  const { data, error } = await s.from("blogs").select().eq("id", id).limit(1).single();
  if (error !== null) {
    if (error.code !== "PGRST116") {
      throw new Error(JSON.stringify(error));
    }
    console.error(`blog with id ${id} not found`);
  }
  return data;
}

export async function insertBlog(p: Post) {
  const f = await s.from("blogs").insert({ author: p.author, created_at: p.dateObj.toISOString(), id: p.id, url: p.absoluteUrl });
  if (f.error !== null) {
    if (f.error.code !== "23505") {
      throw new Error(JSON.stringify(f.error));
    }
    console.error(`blog with id ${p.id} is alrady in db`);
  }
  return f.status;
}

export async function selectBlogs() {
  const { data, error } = await s.from('blogs').select();
  if (data === null) {
    throw new Error("DB returned no blogs");
  }
  if (error !== null) {
    throw new Error(JSON.stringify(error));
  }
  return data;
}

export async function selectCommentsFromBlogId(id: number) {
  const { data, error } = await s.from('comments').select('author, body, created_at').eq('blog_id', id);
  if (error !== null) {
    console.error(error);
  }
  return data;
}


export async function insertComment(body: string, author: string, blogId: number) {
  const { data, error } = await s.from("comments").insert({ blog_id: blogId, body: body, author: author }).limit(1).select().single();

  if (error !== null) {
    console.error(error);
  }

  return data;

} 
