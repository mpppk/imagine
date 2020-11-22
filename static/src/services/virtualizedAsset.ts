import { Asset, WorkSpace } from '../models/models';

export interface VirtualizedAssetProps {
  workspace: WorkSpace | null;
  hasMoreAssets: boolean;
  assets: Asset[];
  isScanningAssets: boolean;
  onRequestNextPage: () => null;
}

export const getVirtualizedAssetsProps = (props: VirtualizedAssetProps) => {
  // If there are more items to be loaded then add an extra row to hold a loading indicator.
  const assetCount = props.hasMoreAssets
    ? props.assets.length + 1
    : props.assets.length;

  // Only load 1 page of items at a time.
  // Pass an empty callback to InfiniteLoader in case it asks us to load more than once.
  const loadMoreAssets = props.isScanningAssets
    ? () => null
    : props.onRequestNextPage;

  // Every row is loaded except for our loading indicator row.
  const isAssetLoaded = (index: number) =>
    !props.hasMoreAssets || index < props.assets.length;

  return {
    assetCount,
    loadMoreAssets,
    isAssetLoaded,
  };
};
