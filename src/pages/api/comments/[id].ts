import type { APIRoute } from "astro";
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
                body: escapeHTML(body).replace(/\n/g, '<br>'),
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
        <head>
        <style>
        article {
            background-color: #282a2e;
            color: #c5c8c6;
            border: 1px solid #111;
            padding: 5px;
        }
        input:checked + span {
            display: none;
        }
        </style>
        </head>
        <div>
        <article class="w-fit">
            <label for="remover" class="text-md">
              ${c.author} - ${c.created_at.toLocaleString("en-GB")} - id: ${c.id}
            </label>
            <input id="remover" type="checkbox" />
            <span>
            <p style="padding-left: 10px; padding-right: 10px; word-wrap: break-word; word-break: break-all;">${c.body}</p>
            </span>
          </article>
        <br/>
        </div>
        `;
        }).join("\n");
    return new Response(divs, { status: 200, headers: { "Content-Type": "text/html" } },);
};
