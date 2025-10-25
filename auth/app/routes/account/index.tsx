import { AccountView } from "@daveyplate/better-auth-ui"
import { useParams } from "react-router"

export default function AccountPage() {
  const params = useParams()
  const pathname = params["*"] || ""

    return (
        <main className="container p-4 md:p-6">
            <AccountView path={pathname} />
        </main>
    )
}
