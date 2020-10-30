import {Button} from "@material-ui/core";
import Dialog from "@material-ui/core/Dialog";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import DialogTitle from "@material-ui/core/DialogTitle";
import React, {useState} from "react";
import FormControl from "@material-ui/core/FormControl";
import Select from "@material-ui/core/Select";
import MenuItem from "@material-ui/core/MenuItem";
import {makeStyles} from "@material-ui/core/styles";
import TextField from "@material-ui/core/TextField";
import Grid from "@material-ui/core/Grid";
import {QueryInput as QueryIn, QueryOp} from "../models/models";
import InputLabel from "@material-ui/core/InputLabel";
import {immutableSplice} from "../util";
import Switch from "@material-ui/core/Switch";

const useStyles = makeStyles((theme) => ({
  formControl: {
    margin: theme.spacing(1),
    minWidth: 120,
  },
  queryInputGrid: {
    marginBottom: theme.spacing(2),
  },
}));

interface QueryInputProps extends QueryIn {
  onUpdate: (q: QueryIn) => void;
}

// tslint:disable-next-line:variable-name
export const QueryInput: React.FC<QueryInputProps> = (props) => {
  const classes = useStyles();
  const handleChangeSelect = (e: React.ChangeEvent<{ name?: string; value: unknown }>) => {
    props.onUpdate({op: e.target.value as QueryOp, tagName: props.tagName});
  };

  const handleChangeTagName = (e: React.ChangeEvent<HTMLInputElement>) => {
    props.onUpdate({op: props.op, tagName: e.target.value});
  };

  return (
    <>
      <Grid item={true} xs={6} className={classes.queryInputGrid}>
        <FormControl className={classes.formControl}>
          <InputLabel>op</InputLabel>
          <Select
            label={'op'}
            value={props.op}
            onChange={handleChangeSelect}
          >
            <MenuItem value={'equals'}>Equals</MenuItem>
            <MenuItem value={'not-equals'}>Not Equals</MenuItem>
          </Select>
        </FormControl>
      </Grid>
      <Grid item={true} xs={6}>
        <TextField id="tag" label="tag" value={props.tagName} onChange={handleChangeTagName}/>
      </Grid>
    </>
  );
}

interface QueryFormProps {
  inputs: QueryInWithID[];
  onUpdate: (inputs: QueryInWithID[]) => void;
}

// tslint:disable-next-line:variable-name
export const QueryForm: React.FC<QueryFormProps> = (props) => {
  const handleClickAddButton = () => {
    props.onUpdate([...props.inputs, {op: 'equals', tagName: '', id: -1}]);
  };

  const handleUpdateQueryInput = (index: number, input: QueryIn) => {
    props.onUpdate(immutableSplice(props.inputs, index, 1, {...input, id: props.inputs[index].id}));
  };

  return (
    <>
      <Grid container={true} spacing={1}>
        {props.inputs.map((input, i) =>
          <QueryInput
            key={input.id}
            op={input.op}
            tagName={input.tagName}
            onUpdate={handleUpdateQueryInput.bind(null, i)}
          />)}
        <Button variant="outlined" color="primary" onClick={handleClickAddButton}>
          Add
        </Button>
      </Grid>
    </>
  );
};

interface FilterDialogProps {
  enabled: boolean
  open: boolean
  onClose: () => void
  inputs: QueryIn[]
  onClickApplyButton: (enabled: boolean, inputs: QueryIn[]) => void
  onChangeEnableSwitch: (enabled: boolean) => void;
}

interface QueryInWithID extends QueryIn {
  id: number;
}

const toQueryInWithIDList = (qi: QueryIn[]): QueryInWithID[] => {
  return qi.map((i, ind): QueryInWithID => ({...i, id: ind}));
}

// tslint:disable-next-line:variable-name
export const FilterDialog: React.FC<FilterDialogProps> = (props) => {
  const [ops, setOps] = useState(toQueryInWithIDList(props.inputs));
  const [opCnt, setOpCnt] = useState(props.inputs.length);
  // const classes = useStyles();

  const handleUpdateQueryForm = (inputs: QueryInWithID[]) => {
    let cnt = 0;
    setOps(inputs.map((i) => {
      return i.id === -1 ? {...i, id: opCnt + ++cnt} : i;
    }));
    if (cnt > 0) {
      setOpCnt(opCnt + cnt);
    }
  };

  const handleClickCancelButton = () => {
    setOps(toQueryInWithIDList(props.inputs));
    props.onClose();
  }

  const handleClickApplyButton = () => {
    props.onClickApplyButton(props.enabled, ops);
  };

  const handleChangeEnableSwitch = (e: React.ChangeEvent<HTMLInputElement>) => {
    props.onChangeEnableSwitch(e.target.checked);
  }

  return (
    <Dialog open={props.open} onClose={props.onClose} aria-labelledby="form-dialog-title">
      <DialogTitle>
        Filter settings
        <Switch
          checked={props.enabled}
          onChange={handleChangeEnableSwitch}
          color="primary"
        />
      </DialogTitle>
      <DialogContent>
        <QueryForm inputs={ops} onUpdate={handleUpdateQueryForm}/>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClickCancelButton}>
          Cancel
        </Button>
        <Button onClick={handleClickApplyButton} color="primary">
          Apply
        </Button>
      </DialogActions>
    </Dialog>
  )
}
