import { globBlogs } from "@utils/glob";
import type { APIRoute, GetStaticPaths } from "astro";
import { db, Comment, eq, NOW } from 'astro:db';

export const prerender = false;

function escapeHTML(str: string) {
    return str.replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;');
}
export const POST: APIRoute = async (ctx) => {
    try {
        const formData = await ctx.request.formData();
        const author = formData.get("author")?.toString() || "Anonymous";
        const body = formData.get("body")?.toString();
        const id = Number(ctx.params.id);
        if (!id) {
            return new Response(null, { status: 403 });
        }
        if (!body) {
            return new Response("Body is required", { status: 400 });
        }

        await db
            .insert(Comment)
            .values({
                author: escapeHTML(author),
                body: escapeHTML(body),
                created_at: NOW,
                parentId: id
            });

        // Respond with a 204 No Content status, indicating success
        // Use the HX-Redirect header for a client-side redirect, which causes a page reload
        const referer = ctx.request.headers.get("Referer") || '/';
        return new Response(null, {
            status: 204,
            headers: {
                'HX-Redirect': referer,
            },
        });
    } catch (error) {
        console.error("Failed to process form submission", error);
        return new Response("Server error", { status: 500 });
    }
};

export const GET: APIRoute = async (ctx) => {
    console.log(`GET ${ctx.params}`);
    const id = Number(ctx.params.id);
    const comments = await db.select().from(Comment).where(eq(Comment.parentId, id));
    const divs = comments
        .sort((a, b) => b.created_at.getTime() - a.created_at.getTime())
        .map((c) => {
            return `
          <article class="bg-slate-800 text-white w-fit px-2 py-1 my-2">
            <p class="text-green-600 text-md">
              ${c.author} - ${c.created_at.toLocaleString("en-GB")}
            </p>
            <p>${c.body}</p>
          </article>
        `;
        }).join("\n");
    return new Response(divs, { status: 200, headers: { "Content-Type": "text/html" } },);
};