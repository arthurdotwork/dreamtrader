import {createFileRoute, Link, useNavigate} from '@tanstack/react-router'
import {useAssets} from "@/queries/asset.ts";
import {authMiddleware, useAuth} from "@/contexts/auth.tsx";


export const Route = createFileRoute('/assets/')({
    beforeLoad: ({context}) => authMiddleware(context.auth),
    component: AssetsList,
})

function AssetsList() {
    const navigate = useNavigate();
    const auth = useAuth();
    const {data: assets} = useAssets();

    const logout = () => {
        auth.logout()
        navigate({to: '/authenticate'})
    }

    return (
        <div className="h-screen flex">
            <aside className="w-1/5 bg-neutral-100 p-6 flex flex-col">
                <nav className="flex-grow">
                    <h1 className="text-lg font-black mb-4">DreamTrader</h1>
                    <ul className="flex flex-col gap-4">
                        <li>
                            <Link to="/assets" className="[&.active]:bg-neutral-200 [&.active]:font-bold block p-2 bg-neutral-100 hover:bg-neutral-200 rounded-lg transition">Assets</Link>
                        </li>
                    </ul>
                </nav>
                <div>
                    <button onClick={() => logout()} className="text-left w-full cursor-pointer [&.active]:bg-neutral-200 [&.active]:font-bold block p-2 bg-neutral-100 hover:bg-neutral-200 rounded-lg transition">Logout</button>
                </div>
            </aside>
            <div className="p-6 grow">
                <h1 className="text-xl font-black mb-4">Assets</h1>
                <ul>
                    {assets?.map((asset) => (
                        <li key={asset.id}
                            className="mb-4 border border-neutral-200 hover:border-neutral-400 transition hover:bg-neutral-100">
                            <Link to={`/assets/$assetId`}
                                  params={{assetId: asset.id.toString()}}
                                  className="flex justify-between p-4">
                                <p>{asset.name} <span
                                    className="italic">({asset.ticker})</span>
                                </p>
                                <p>See details</p>
                            </Link>
                        </li>
                    ))}
                </ul>
            </div>
        </div>
    )
}
