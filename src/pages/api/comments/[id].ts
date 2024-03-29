import { globBlogs } from "@utils/glob";
import type { APIRoute, GetStaticPaths } from "astro";
import { db, Comment, eq, NOW } from 'astro:db';

export const prerender = false;

function escapeHTML(str: string) {
    return str
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;');
}

export const POST: APIRoute = async (ctx) => {
    const referer = ctx.request.headers.get("Referer") || '/';
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

        return new Response(null, {
            status: 204,
            headers: {
                'HX-Redirect': referer,
            },
        });
    } catch (error) {
        return new Response(`
            <div>Something went wrong. 
                 Please let us know what happened by filing a github issue. 
                 Link to the repo is at the top. 
                 Here is the error: ${JSON.stringify(error)}
            </div> 
        `, {
            status: 200, headers: { "Content-Type": "text/html" }
        });
        ;
    }
};

export const GET: APIRoute = async (ctx) => {
    const id = Number(ctx.params.id);
    if (!id) {
        return new Response(null, { status: 404 });
    }
    const comments = await db.select().from(Comment).where(eq(Comment.parentId, id));
    const divs = comments
        .sort((a, b) => b.created_at.getTime() - a.created_at.getTime())
        .map((c) => {
            return `
          <article class="bg-slate-800 text-white w-fit px-2 py-1 my-2">
            <p class="text-green-600 text-md">
              ${c.author} - ${c.created_at.toLocaleString("en-GB")} - id: ${c.id}
            </p>
            <p>${c.body}</p>
          </article>
        `;
        }).join("\n");
    return new Response(divs, { status: 200, headers: { "Content-Type": "text/html" } },);
};