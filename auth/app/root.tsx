import { Links, Meta, Outlet, Scripts, ScrollRestoration } from "react-router";
import {authClient} from '@/lib/auth'
import {AuthUIProvider} from '@daveyplate/better-auth-ui'
import {Link, useNavigate} from 'react-router'
import { Toaster } from "sonner";
import type { Route } from "./+types/root";
import "./app.css";

export default function App() {
  return (
      <AuthConfigProvider>
        <Outlet />
        <Toaster/>
      </AuthConfigProvider>
  );
}

export function AuthConfigProvider({children}: {children: React.ReactNode}) {
  const navigate = useNavigate()

  return (
    <AuthUIProvider
      authClient={authClient}
      magicLink
      redirectTo="/"
      changeEmail={false}
      navigate={href => {
        navigate(href)
      }}
      replace={href => {
        navigate(href, {
          replace: true,
        })
      }}
      Link={props => <Link {...props} to={props.href} />}
    >
      {children}
    </AuthUIProvider>
  )
}

export function Layout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <Meta />
        <Links />
      </head>
      <body>
        {children}
        <Scripts />
        <ScrollRestoration />
      </body>
    </html>
  );
}

export function ErrorBoundary({ error }: Route.ErrorBoundaryProps) {
  return null;
}
