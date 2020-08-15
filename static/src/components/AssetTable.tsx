import * as React from "react";
import {Toolbar} from "@material-ui/core";
import TableContainer from "@material-ui/core/TableContainer";
import Paper from "@material-ui/core/Paper";
import Table from "@material-ui/core/Table";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";
import TableCell from "@material-ui/core/TableCell";
import TableBody from "@material-ui/core/TableBody";
import {makeStyles} from "@material-ui/core/styles";
import {VirtualizedAssetProps} from "../services/virtualizedAsset";
import {AssetListItemProps, VirtualizedAssetList} from "./VirtualizedAssetList";
import AutoSizer from "react-virtualized-auto-sizer";

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

  // tslint:disable-next-line:variable-name
  const Item: React.FC<AssetListItemProps> = ({asset, isLoaded, style}) => {
    if (!isLoaded) {
      return (<div style={style}>Loading...</div>);
    }

    return (
      <TableRow key={asset.path} style={style}>
        <TableCell component="th" scope="row">
          {asset.id}
        </TableCell>
        <TableCell component="th" scope="row">
          {asset.name}
        </TableCell>
        <TableCell component="th" scope="row">
          {asset.boundingBoxes}
        </TableCell>
        <TableCell component="th" scope="row">
          {asset.path}
        </TableCell>
      </TableRow>
    );
  };

  return (
    <>
      <Toolbar/>
      <TableContainer component={Paper}>
        <Table className={classes.table} size="small">
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>Name</TableCell>
              <TableCell>Tags</TableCell>
              <TableCell>Path</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            <AutoSizer>
              {({height, width}) => {
                return (
                  <VirtualizedAssetList
                    assets={props.assets}
                    hasMoreAssets={props.hasMoreAssets}
                    isScanningAssets={props.isScanningAssets}
                    onRequestNextPage={props.onRequestNextPage}
                    workspace={props.workspace}
                    height={height}
                    itemSize={35}
                    width={width}
                  >
                    {Item}
                  </VirtualizedAssetList>
                );
              }}
            </AutoSizer>

          </TableBody>
        </Table>
      </TableContainer>
    </>
  );
};