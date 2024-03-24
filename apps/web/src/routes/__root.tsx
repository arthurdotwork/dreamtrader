import {
    createRootRouteWithContext,
    Outlet
} from '@tanstack/react-router'
import {TanStackRouterDevtools} from '@tanstack/router-devtools'
import {Toaster} from "@/components/ui/toaster.tsx";
import {type AuthContextType} from "@/contexts/auth.tsx";

interface RouterContext {
    auth: AuthContextType;
}

export const Route = createRootRouteWithContext<RouterContext>()({
    component: () => {
        return (
            (
                <>
                    <Outlet/>
                    <Toaster/>
                    <TanStackRouterDevtools/>
                </>
            )
        )
    },
})
