import * as React from "react";
import {useEffect} from "react";
import {Toolbar} from "@material-ui/core";
import TableContainer from "@material-ui/core/TableContainer";
import Paper from "@material-ui/core/Paper";
import Table from "@material-ui/core/Table";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";
import TableCell from "@material-ui/core/TableCell";
import TableBody from "@material-ui/core/TableBody";
import AutoSizer from "react-virtualized-auto-sizer";
import InfiniteLoader from "react-window-infinite-loader";
import {FixedSizeList} from "react-window";
import {makeStyles} from "@material-ui/core/styles";
import {getVirtualizedAssetsProps, VirtualizedAssetProps} from "../services/virtualizedAsset";

const useStyles = makeStyles({
  table: {
    height: 500,
    minWidth: 650,
  },
});

type Props = VirtualizedAssetProps;

// tslint:disable-next-line:variable-name
export const AssetTable = (props: Props) => {
  const classes = useStyles();
  // Fixme use redux-saga
  useEffect(() => {
    if (props.workspace !== null) {
      props.onRequestNextPage();
    }
  }, [props.workspace]);

  const assetInfo = getVirtualizedAssetsProps(props);

  // tslint:disable-next-line:variable-name
  const Item = ({index, style}: any) => {
    let content;
    if (!assetInfo.isAssetLoaded(index)) {
      content = "Loading...";
    } else {
      content = props.assets[index].path;
    }

    return <div style={style}>{content}</div>;
  };

  return (
    <>
      <Toolbar/>
      <TableContainer component={Paper}>
        <Table className={classes.table} size="small" aria-label="a dense table">
          <TableHead>
            <TableRow>
              <TableCell>Dessert (100g serving)</TableCell>
              <TableCell align="right">Calories</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
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
                      {Item}
                    </FixedSizeList>
                  )}
                </InfiniteLoader>
              )}
            </AutoSizer>
          </TableBody>
        </Table>
      </TableContainer>
    </>
  );
};