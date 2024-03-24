import { createRootRoute, Link, Outlet } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/router-devtools'

export const Route = createRootRoute({
    component: () => (
        <>
            <nav className="p-6">
                <ul className="flex justify-end gap-4">
                    <li>
                        <Link to="/" className="[&.active]:font-bold">
                            Home
                        </Link>{' '}
                    </li>
                    <li>
                        <Link to="/assets" className="[&.active]:font-bold">
                            Assets
                        </Link>
                    </li>
                </ul>

            </nav>
            <hr />
            <Outlet />
            <TanStackRouterDevtools />
        </>
    ),
})
