import * as React from "react";
import {useEffect, useRef} from "react";
import {FixedSizeList, ListChildComponentProps} from "react-window";
import InfiniteLoader from "react-window-infinite-loader";
import {getVirtualizedAssetsProps, VirtualizedAssetProps} from "../services/virtualizedAsset";
import {Asset} from "../models/models";

export interface AssetListItemProps {
  style: React.CSSProperties
  asset: Asset
  index: number
  isLoaded: boolean
}

interface Props extends VirtualizedAssetProps {
  height: number
  width: number
  itemSize: number
  children: React.FC<AssetListItemProps>
  selectedIndex: number
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
          style={style as React.CSSProperties}
          asset={props.assets![index]}
          index={index}
          isLoaded={assetInfo.isAssetLoaded(index)}
        />
      );
    };

    const loaderRef = useRef<any>(null);
    useEffect(() => {
      if (loaderRef.current !== null) {
        loaderRef.current._listRef.scrollToItem(props.selectedIndex);
      }
    }, [props.selectedIndex])

    return (
      <InfiniteLoader
        ref={loaderRef}
        isItemLoaded={assetInfo.isAssetLoaded}
        itemCount={assetInfo.assetCount}
        loadMoreItems={assetInfo.loadMoreAssets}
      >
        {({onItemsRendered, ref}) => (
          <FixedSizeList
            className="List"
            height={props.height}
            itemCount={assetInfo.assetCount}
            itemSize={props.itemSize}
            onItemsRendered={onItemsRendered}
            ref={ref}
            width={props.width}
          >
            {AssetItem}
          </FixedSizeList>
        )}
      </InfiniteLoader>
    );
  }
;