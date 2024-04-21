import type { APIRoute } from "astro";
import { db, Comment, eq, NOW } from "astro:db";

export const prerender = false;

function escapeHTML(str: string) {
    return str
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
}

function renderRefrences(str: string): string {
    let ref = "";
    for (let i = 8; i < str.length; i++) {
        if (str[i] === " ") {
            break;
        }
        ref += str[i];
    }
    let num = Number(ref);
    if (Number.isNaN(num) || num === 0) {
        return str;
    }
    let reference = `<a href="#${num}" class="reply">&gt;&gt;${num}</a>`;
    let result = reference + str.substring(8 + num.toString().length);

    return result;
}

function renderGreentext(str: string): string {
    let text = str.substring(4);
    if (str.at(5) === " ") {
        return str;
    }
    let result = `<span class="greentext">&gt;${text}</span>`;
    return result;
}

function renderComment(str: string): string {
    str = escapeHTML(str);
    let lines = str.split("\n");
    let renderedLines = [];

    for (const l of lines) {
        let currentLine = l;
        if (l.substring(0, 8) === "&gt;&gt;") {
            currentLine = renderRefrences(l);
        } else if (l.substring(0, 4) === "&gt;") {
            currentLine = renderGreentext(l);
        }
        renderedLines.push(currentLine);
    }

    return renderedLines.join("<br/>");
}

export const POST: APIRoute = async (ctx) => {
    const referer = ctx.request.headers.get("Referer") || "/";
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
                body: renderComment(body),
                created_at: NOW,
                parentId: id
            });

        return new Response(null, {
            status: 204,
            headers: {
                "HX-Redirect": referer,
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
        .sort((a, b) => a.created_at.getTime() - b.created_at.getTime()) // sorted from oldest -> newest
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
        .greentext {
            color: #789922;
        }
        .reply {
            color: #d00;
        }
        </style>
        </head>
        <div>
        <article id="${c.id}" class="w-fit">
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
