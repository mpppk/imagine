import AutoSizer from "react-virtualized-auto-sizer";
import * as React from "react";
import {FixedSizeList, ListChildComponentProps} from "react-window";
import InfiniteLoader from "react-window-infinite-loader";
import {getVirtualizedAssetsProps, VirtualizedAssetProps} from "../services/virtualizedAsset";
import {Asset} from "../models/models";
import {useEffect} from "react";

export interface AssetListItemProps extends Pick<ListChildComponentProps, 'style'>{
  asset: Asset
  isLoaded: boolean
}

interface Props extends VirtualizedAssetProps {
  children: React.FC<AssetListItemProps>
}

// tslint:disable-next-line:variable-name
export const VirtualizedAssetList: React.FC<Props> = (props) => {
  // Fixme use redux-saga
  useEffect(() => {
    if (props.workspace !== null) {
      props.onRequestNextPage();
    }
  }, [props.workspace]);

  const assetInfo = getVirtualizedAssetsProps(props);

  // tslint:disable-next-line:variable-name
  const AssetItem: React.FC<ListChildComponentProps> = ({index, style}) => {
    // tslint:disable-next-line:variable-name
    const Children = props.children;

    return (
      <Children
        style={style}
        asset={props.assets[index]}
        isLoaded={assetInfo.isAssetLoaded(index)}
      />
    );
  };

  return (
    <AutoSizer>
      {({height, width}) => (
        <InfiniteLoader
          isItemLoaded={assetInfo.isAssetLoaded}
          itemCount={assetInfo.assetCount}
          loadMoreItems={assetInfo.loadMoreAssets}
        >
          {({onItemsRendered, ref}) => (
            <FixedSizeList
              className="List"
              height={height}
              itemCount={assetInfo.assetCount}
              itemSize={30}
              onItemsRendered={onItemsRendered}
              ref={ref}
              width={width}
            >
              {AssetItem}
            </FixedSizeList>
          )}
        </InfiniteLoader>
      )}
    </AutoSizer>
  );
};