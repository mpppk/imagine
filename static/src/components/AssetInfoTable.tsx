import {TableHead, Theme} from "@material-ui/core";
import {makeStyles} from "@material-ui/core/styles";
import React from "react";
import {Asset} from "../models/models";
import TableContainer from "@material-ui/core/TableContainer";
import Paper from "@material-ui/core/Paper";
import Table from "@material-ui/core/Table";
import TableRow from "@material-ui/core/TableRow";
import TableCell from "@material-ui/core/TableCell";
import TableBody from "@material-ui/core/TableBody";

const useStyles = makeStyles((_theme: Theme) => {
  return {
    table: {},
  }
});

interface Props {
  asset: Asset
}

// tslint:disable-next-line:variable-name
export const AssetInfoTable: React.FC<Props> = ({asset}) => {
  const classes = useStyles();
  const tags = (asset.boundingBoxes ?? []).map((b) => b.tag);

  return (
    <TableContainer component={Paper}>
      <Table className={classes.table} size="small" aria-label="asset information table">
        <TableHead>
          <TableRow>
            <TableCell>Asset information</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          <TableRow>
            <TableCell component="th" scope="row">ID</TableCell>
            <TableCell>{asset.id}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell component="th" scope="row">Path</TableCell>
            <TableCell>{asset.path}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell component="th" scope="row">Tags</TableCell>
            <TableCell>{tags}</TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </TableContainer>)
}
