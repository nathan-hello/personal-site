import { insertComment, selectBlogById } from "@utils/supabase";
import type { APIRoute } from "astro";

export const POST: APIRoute = async ({ request, params }) => {
    const data = await request.formData();
    const body = data.get("name")?.toString();
    let author = data.get("author")?.toString();

    if (!author) {
        author = "anonymous";
    }
    if (!body || body.length > 3 || body.length < 1000) {
        return new Response(
            JSON.stringify({
                message: "Missing required fields",
            }),
            { status: 400 }
        );
    }

    if (!params.id
        || !Number(params.id)
        || !selectBlogById(Number(params.id))
    ) {
        return new Response(
            JSON.stringify({
                message: "bad url",
            }),
            { status: 400 }
        );
    }

    await insertComment(body, author, Number(params.id));
    return new Response(JSON.stringify({ status: 200 }));
};
