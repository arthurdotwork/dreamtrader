import {useQuery} from "@tanstack/react-query";

const assets: Asset[] = [
    { id: 1, name: 'Microsoft', ticker: 'MSFT' },
    { id: 2, name: 'Apple', ticker: 'AAPL' },
    { id: 3, name: 'Google', ticker: 'GOOGL' },
]

type Asset = {
    id: number
    name: string
    ticker: string
}

export const useAssets = () => useQuery({queryKey: ['assets'], queryFn: () => assets})
export const useAsset = (assetId: string) => useQuery({queryKey: ['asset', assetId], queryFn: () => assets.find((a) => a.id.toString() === assetId)})
