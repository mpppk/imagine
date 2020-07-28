import Paper from '@material-ui/core/Paper';
import {makeStyles} from '@material-ui/core/styles';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import * as React from 'react';
import AutoSizer from 'react-virtualized-auto-sizer';
import {FixedSizeList as List} from 'react-window';
import {Toolbar} from "@material-ui/core";

const useStyles = makeStyles({
  table: {
    height: 500,
    minWidth: 650,
  },
});

// tslint:disable-next-line:variable-name
const AssetTableRow = () => {
  return (
    <TableRow key={'name'}>
      <TableCell component="th" scope="row">
        name
      </TableCell>
      <TableCell align="right">{'calories'}</TableCell>
    </TableRow>
  );
};

// tslint:disable-next-line:variable-name
const AssetTable = () => {
  const classes = useStyles();
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
                <List
                  width={width}
                  height={height}
                  itemCount={1000}
                  itemSize={35}
                >
                  {AssetTableRow}
                </List>
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

