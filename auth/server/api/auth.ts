import { auth } from "server/auth";
import type { ActionFunctionArgs, LoaderFunctionArgs } from "react-router";

export async function loader({ request }: LoaderFunctionArgs) {
  try {
    const response = await auth.handler(request);
    return response;
  } catch (error) {
    console.error("Auth loader error:", error);
    return new Response("Auth error", { status: 500 });
  }
}

export async function action({ request }: ActionFunctionArgs) {
  try {
    const response = await auth.handler(request);
    return response;
  } catch (error) {
    console.error("Auth action error:", error);
    return new Response("Auth error", { status: 500 });
  }
}
