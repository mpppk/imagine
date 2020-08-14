import {Toolbar} from "@material-ui/core";
import Paper from '@material-ui/core/Paper';
import {makeStyles} from '@material-ui/core/styles';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import * as React from 'react';
import {useEffect} from "react";
import {useDispatch, useSelector} from "react-redux";
import AutoSizer from 'react-virtualized-auto-sizer';
import {FixedSizeList} from "react-window";
import InfiniteLoader from 'react-window-infinite-loader';
import {assetActionCreators} from "../actions/asset";
import {State} from "../reducers/reducer";

const useStyles = makeStyles({
  table: {
    height: 500,
    minWidth: 650,
  },
});

// tslint:disable-next-line:variable-name
const AssetTable = () => {
  const classes = useStyles();
  const dispatch = useDispatch();
  const globalState = useSelector((s: State) => s.global);

  const loadNextPage = () => {
    if (globalState.currentWorkSpace !== null) {
      dispatch(assetActionCreators.requestAssets({
        requestNum: 10,
        workSpaceName: globalState.currentWorkSpace.name,
      }));
    }
    return null;
  };

  // Fixme use redux-saga
  useEffect(() => {
    const n = globalState.currentWorkSpace?.name;
    if (n !== undefined) {
      loadNextPage();
    }
  }, [globalState.currentWorkSpace]);

  // If there are more items to be loaded then add an extra row to hold a loading indicator.
  const itemCount = globalState.hasMoreAssets ? globalState.assets.length + 1 : globalState.assets.length;

  // Only load 1 page of items at a time.
  // Pass an empty callback to InfiniteLoader in case it asks us to load more than once.
  const loadMoreItems = globalState.isScanningAssets ? (() => null) : loadNextPage;
  // Every row is loaded except for our loading indicator row.
  const isItemLoaded = (index: number) => !globalState.hasMoreAssets || index < globalState.assets.length;

  // Render an item or a loading indicator.
  // tslint:disable-next-line:variable-name
  const Item = ({index, style}: any) => {
    let content;
    if (!isItemLoaded(index)) {
      content = "Loading...";
    } else {
      content = globalState.assets[index].path;
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
}

export default () => {
  return (
    <div>
      <AssetTable/>
    </div>
  );
}

