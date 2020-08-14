import AutoSizer from "react-virtualized-auto-sizer";
import * as React from "react";
import {FixedSizeList, ListChildComponentProps} from "react-window";
import InfiniteLoader from "react-window-infinite-loader";
import {getVirtualizedAssetsProps, VirtualizedAssetProps} from "../services/virtualizedAsset";

// export type OriginalAssetItem = React.FC<Pick<ListChildComponentProps, 'index' | 'style'>>
export type OriginalAssetItemProps = Pick<ListChildComponentProps, 'index' | 'style'>

interface Props extends VirtualizedAssetProps {
  children: React.FC<OriginalAssetItemProps>
}

// tslint:disable-next-line:variable-name
export const VirtualizedAssetList: React.FC<Props> = (props) => {
  const assetInfo = getVirtualizedAssetsProps(props);

  // tslint:disable-next-line:variable-name
  const AssetItem: React.FC<ListChildComponentProps> = ({index, style}) => {
    // tslint:disable-next-line:variable-name
    const Children = props.children;

    return (
      <Children
        style={style}
        index={index}
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