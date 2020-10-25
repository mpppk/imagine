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
    textAlign: 'center',
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
  height: number
}

// tslint:disable-next-line:variable-name
const MemoizedImg = React.memo((props: {src: string}) => <img src={props.src} height={'100%'}/>);

const generateImageGridTile = (selectedIndex: number, handleClickImage: (path: string, index: number) => void): React.FC<AssetListItemProps> => ({asset, index, isLoaded, style}) => {
  if (!isLoaded) {
    return (<div style={style}>Loading...</div>);
  }

  const classes = useStyles();

  const pathUrl = assetPathToUrl(asset.path);

  const genClickImageHandler = (imgPath: string, i: number) => () => {
    handleClickImage(imgPath, i);
  };

  const selectedTileStyle: CSSProperties = {border: "solid"};
  const newStyle = index === selectedIndex ? {...style, ...selectedTileStyle} : style;


  return (
    <GridListTile
      style={newStyle}
      key={asset.path}
      cols={1}
      onClick={genClickImageHandler(pathUrl, index)}
      className={classes.gridListTile}
    >
      <MemoizedImg src={pathUrl} />
    </GridListTile>
  );
};

// tslint:disable-next-line:variable-name
export const ImageGridList: React.FC<Props> = (props) => {
  const classes = useStyles();

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
          height={props.height}
          itemSize={props.cellHeight}
          width={props.width}
        >
          {generateImageGridTile(props.selectedIndex, props.onClickImage)}
        </VirtualizedAssetList>
      </GridList>
    </div>
  );
}