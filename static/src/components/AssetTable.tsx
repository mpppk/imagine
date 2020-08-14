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
import {Asset} from "../models/models";

const useStyles = makeStyles({
  table: {
    height: 500,
    minWidth: 650,
  },
});

interface Props {
  workspaceName: string | null
  hasMoreAssets: boolean
  assets: Asset[]
  isScanningAssets: boolean
  onRequestNextPage: () => null
}

// tslint:disable-next-line:variable-name
export const AssetTable = (props: Props) => {
  const classes = useStyles();
  // Fixme use redux-saga
  useEffect(() => {
    if (props.workspaceName !== null) {
      props.onRequestNextPage();
      // loadNextPage();
    }
  }, [props.workspaceName]);

  // If there are more items to be loaded then add an extra row to hold a loading indicator.
  const itemCount = props.hasMoreAssets ? props.assets.length + 1 : props.assets.length;

  // Only load 1 page of items at a time.
  // Pass an empty callback to InfiniteLoader in case it asks us to load more than once.
  const loadMoreItems = props.isScanningAssets ? (() => null) : props.onRequestNextPage;
  // Every row is loaded except for our loading indicator row.
  const isItemLoaded = (index: number) => !props.hasMoreAssets || index < props.assets.length;

  // Render an item or a loading indicator.
  // tslint:disable-next-line:variable-name
  const Item = ({index, style}: any) => {
    let content;
    if (!isItemLoaded(index)) {
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
                  isItemLoaded={isItemLoaded}
                  itemCount={itemCount}
                  loadMoreItems={loadMoreItems}
                >
                  {({onItemsRendered, ref}) => (
                    <FixedSizeList
                      className="List"
                      height={height}
                      itemCount={itemCount}
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