import {createFileRoute, Link} from "@tanstack/react-router";
import {useAsset} from "@/queries/asset.ts";

export const Route = createFileRoute('/assets/$assetId')({
    component: AssetPage,
})

function AssetPage() {
    const {assetId} = Route.useParams()
    const {data: asset} = useAsset(assetId)
    if (!asset) {
        return <div>Asset not found</div>
    }

    return (
        <div className="p-6">
            <h1 className="text-xl font-black mb-4">Asset Details</h1>
            <p className="mb-4">{asset.name} <span
                className="italic">({asset.ticker})</span></p>
            <Link to="/assets" className="text-blue-500 hover:underline">
                Back to assets
            </Link>
        </div>
    )
}
