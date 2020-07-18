import * as React from 'react';
// import MaterialTable from "material-table";
import dynamic from "next/dynamic";
import {Button, FormControl} from "@material-ui/core";
import MenuItem from "@material-ui/core/MenuItem";
import Select from "@material-ui/core/Select";
import InputLabel from "@material-ui/core/InputLabel";
import {makeStyles} from "@material-ui/core/styles";

const MaterialTable = dynamic(() => import("material-table"), { ssr: false });

const useStyles = makeStyles((theme) => ({
  formControl: {
    margin: theme.spacing(1),
    minWidth: 120,
  },
  selectEmpty: {
    marginTop: theme.spacing(2),
  },
}));

export default () => {
  const classes = useStyles();
  const [age] = React.useState('');

  return (
    <div style={{ maxWidth: '100%' }}>
      <FormControl className={classes.formControl}>
        <InputLabel id="demo-simple-select-label">Age</InputLabel>
        <Select
          labelId="demo-simple-select-label"
          id="demo-simple-select"
          value={age}
        >
          <MenuItem value={10}>Ten</MenuItem>
          <MenuItem value={20}>Twenty</MenuItem>
          <MenuItem value={30}>Thirty</MenuItem>
        </Select>
      </FormControl>
      <MaterialTable
        columns={[
          { title: 'path', field: 'name' },
        ]}
        data={[{ path: 'xxx'}]}
        title="Assets"
      />
      <Button variant="outlined" color="primary">
        Load DB
      </Button>
    </div>
  );
}
