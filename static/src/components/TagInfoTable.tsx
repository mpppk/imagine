import { TableHead, Theme } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';
import TableContainer from '@material-ui/core/TableContainer';
import Paper from '@material-ui/core/Paper';
import Table from '@material-ui/core/Table';
import TableRow from '@material-ui/core/TableRow';
import TableCell from '@material-ui/core/TableCell';
import TableBody from '@material-ui/core/TableBody';

const useStyles = makeStyles((_theme: Theme) => {
  return {
    table: {},
  };
});

interface Props {
  className: string;
  tagID: number;
  tagName: string;
}

// tslint:disable-next-line:variable-name
export const TagInfoTable: React.FC<Props> = (props) => {
  const classes = useStyles();

  return (
    <TableContainer className={props.className} component={Paper}>
      <Table
        className={classes.table}
        size="small"
        aria-label="tag information table"
        data-cy="tag-information-table"
      >
        <TableHead>
          <TableRow>
            <TableCell>Tag</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          <TableRow>
            <TableCell component="th" scope="row">
              ID
            </TableCell>
            <TableCell data-cy="tag-id">{props.tagID}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell component="th" scope="row">
              Name
            </TableCell>
            <TableCell data-cy="tag-name">{props.tagName}</TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </TableContainer>
  );
};
