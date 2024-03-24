import {createLazyFileRoute, Link} from '@tanstack/react-router'
import {useAssets} from "@/queries/asset.ts";

export const Route = createLazyFileRoute('/assets/')({
    component: AssetsList,
})

function AssetsList() {
    const {data: assets} = useAssets();

    return (
        <div className="p-6">
            <h1 className="text-xl font-black mb-4">Assets</h1>
            <ul>
                {assets?.map((asset) => (
                    <li key={asset.id} className="mb-4 border border-neutral-200 hover:border-neutral-400 transition hover:bg-neutral-100">
                        <Link to={`/assets/$assetId`} params={{ assetId: asset.id.toString() }} className="flex justify-between p-4">
                        <p>{asset.name} <span className="italic">({asset.ticker})</span></p>
                            <p>See details</p>
                        </Link>
                    </li>
                ))}
            </ul>
        </div>
    )
}
