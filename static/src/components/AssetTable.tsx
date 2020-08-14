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
import {makeStyles} from "@material-ui/core/styles";
import {getVirtualizedAssetsProps, VirtualizedAssetProps} from "../services/virtualizedAsset";
import {AssetListItemProps, VirtualizedAssetList} from "./VirtualizedAssetList";

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

  // tslint:disable-next-line:variable-name
  const Item: React.FC<AssetListItemProps> = ({asset, isLoaded, style}) => {
    let content;
    if (!isLoaded) {
      content = "Loading...";
    } else {
      content = asset.path;
    }

    return (<div style={style as React.CSSProperties}>{content}</div>);
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
            <VirtualizedAssetList
              assets={props.assets}
              hasMoreAssets={props.hasMoreAssets}
              isScanningAssets={props.isScanningAssets}
              onRequestNextPage={props.onRequestNextPage}
              workspace={props.workspace}
            >
              {Item}
            </VirtualizedAssetList>
          </TableBody>
        </Table>
      </TableContainer>
    </>
  );
};