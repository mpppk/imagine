import GridList from '@material-ui/core/GridList';
import GridListTile from '@material-ui/core/GridListTile';
import { makeStyles } from '@material-ui/core/styles';
import React, {CSSProperties} from 'react';
import {VirtualizedAssetProps} from "../services/virtualizedAsset";
import {AssetListItemProps, VirtualizedAssetList} from "./VirtualizedAssetList";
import {assetPathToUrl} from "../util";

const useStyles = makeStyles((theme) => ({
  gridList: {
    height: '100%',
  },
  gridListTile: {
    cursor: 'pointer',
  },
  root: {
    backgroundColor: theme.palette.background.paper,
    display: 'flex',
    flexWrap: 'wrap',
    justifyContent: 'space-around',
    overflow: 'hidden',
  },
}));

interface Props extends VirtualizedAssetProps {
  cellHeight: number;
  width: number;
  paths: string[]
  onClickImage: (path: string, index: number) => void
  selectedIndex: number
}

// tslint:disable-next-line:variable-name
export const ImageGridList: React.FC<Props> = (props) => {
  const classes = useStyles();

  const genClickImageHandler = (imgPath: string, index: number) => () => {
    props.onClickImage(imgPath, index);
  };

  // tslint:disable-next-line:variable-name
  const ImageGridTile: React.FC<AssetListItemProps> = ({asset, index, isLoaded, style}) => {
    if (!isLoaded) {
      return (<div style={style}>Loading...</div>);
    }

    const pathUrl = assetPathToUrl(asset.path);

    const selectedTileStyle: CSSProperties = {border: "solid"};
    const newStyle = index === props.selectedIndex ? {...style, ...selectedTileStyle} : style;

    return (
      <GridListTile
        style={newStyle}
        key={asset.path}
        cols={1}
        onClick={genClickImageHandler(pathUrl, index)}
        className={classes.gridListTile}
      >
        <img src={pathUrl} />
      </GridListTile>
    );
  };

  return (
    <div className={classes.root}>
      <GridList cellHeight={props.cellHeight} className={classes.gridList} cols={1}>
        <VirtualizedAssetList
          assets={props.assets}
          selectedIndex={props.selectedIndex}
          hasMoreAssets={props.hasMoreAssets}
          isScanningAssets={props.isScanningAssets}
          onRequestNextPage={props.onRequestNextPage}
          workspace={props.workspace}
          height={window.innerHeight}
          itemSize={props.cellHeight}
          width={props.width}
        >
          {ImageGridTile}
        </VirtualizedAssetList>
      </GridList>
    </div>
  );
}