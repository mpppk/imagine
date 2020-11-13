import {Button, Tooltip} from "@material-ui/core";
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
import {Query, QueryOp} from "../models/models";
import InputLabel from "@material-ui/core/InputLabel";
import {immutableSplice} from "../util";
import Switch from "@material-ui/core/Switch";
import IconButton from "@material-ui/core/IconButton";
import DeleteIcon from "@material-ui/icons/Delete";

const useStyles = makeStyles((theme) => ({
  formControl: {
    margin: theme.spacing(1),
    minWidth: 120,
  },
  queryInputGrid: {
    marginBottom: theme.spacing(2),
  },
}));

interface PathFormProps {
  value: string
  onChangePath: (path: string) => void;
}

// tslint:disable-next-line:variable-name
const PathForm: React.FC<PathFormProps> = (props) => {
  const handleChangePath = (e: React.ChangeEvent<HTMLInputElement>) => {
    props.onChangePath(e.target.value);
  };
  return (
    <TextField label="path" value={props.value} onChange={handleChangePath}/>
  )
};


interface TagNameFormProps {
  value: string
  onChangeTagName: (tagName: string) => void;
}

// tslint:disable-next-line:variable-name
const TagNameForm: React.FC<TagNameFormProps> = (props) => {
  const handleChangeTagName = (e: React.ChangeEvent<HTMLInputElement>) => {
    props.onChangeTagName(e.target.value);
  };
  return (
    <TextField id="tag" label="tag" value={props.value} onChange={handleChangeTagName}/>
  )
};

interface StartWithFormProps {
  prefix: string
  onChangePrefix: (tagName: string) => void;
}

// tslint:disable-next-line:variable-name
const StartWithForm: React.FC<StartWithFormProps> = (props) => {
  const handleChangePrefix = (e: React.ChangeEvent<HTMLInputElement>) => {
    props.onChangePrefix(e.target.value);
  };
  return (
    <TextField label="prefix" value={props.prefix} onChange={handleChangePrefix}/>
  )
};

interface QueryInputFormProps {
  op: QueryOp
  value: string
  onChangeValue: (value: string) => void;
}

interface QueryInputProps {
  op: QueryOp;
  value: string
  onUpdate: (q: Query) => void;
  onDelete: () => void;
}

// tslint:disable-next-line:variable-name
const QueryInputForm: React.FC<QueryInputFormProps> = (props) => {
  switch (props.op) {
    case 'equals':
    case 'not-equals':
      return <TagNameForm value={props.value} onChangeTagName={props.onChangeValue}/>
    case 'start-with':
      return <StartWithForm prefix={props.value} onChangePrefix={props.onChangeValue}/>
    case 'no-tags':
      return <div/>;
    case 'path-equals':
      return <PathForm value={props.value} onChangePath={props.onChangeValue}/>
    default:
      throw new Error("unknown query op is provided. " + props.op)
  }
}

// tslint:disable-next-line:variable-name
export const QueryInput: React.FC<QueryInputProps> = (props) => {
  const classes = useStyles();
  const handleChangeSelect = (e: React.ChangeEvent<{ name?: string; value: unknown }>) => {
    props.onUpdate({op: e.target.value as QueryOp, value: props.value});
  };

  const handleChangeFormValue = (v: string)=> {
    props.onUpdate({op: props.op, value: v});
  };

  return (
    <>
      <Grid item={true} xs={5} className={classes.queryInputGrid} >
        <FormControl className={classes.formControl}>
          <InputLabel>op</InputLabel>
          <Select
            label={'op'}
            value={props.op}
            onChange={handleChangeSelect}
          >
            <MenuItem value={'equals'}>Equals</MenuItem>
            <MenuItem value={'not-equals'}>Not Equals</MenuItem>
            <MenuItem value={'start-with'}>Start With</MenuItem>
            <MenuItem value={'no-tags'}>No Tags</MenuItem>
            <MenuItem value={'path-equals'}>Path Equals</MenuItem>
          </Select>
        </FormControl>
      </Grid>
      <Grid item={true} xs={5}>
        <QueryInputForm op={props.op} value={props.value} onChangeValue={handleChangeFormValue}/>
      </Grid>
      <Grid item={true} xs={2}>
        <Tooltip title="delete query" aria-label="delete-query">
          <IconButton
            aria-label="delete-query"
            onClick={props.onDelete}
          >
            <DeleteIcon/>
          </IconButton>
        </Tooltip>
      </Grid>
    </>
  );
}

interface QueryFormProps {
  inputs: QueryInWithID[];
  onUpdate: (inputs: QueryInWithID[]) => void;
  onDelete: (id: number) => void;
}

// tslint:disable-next-line:variable-name
export const QueryForm: React.FC<QueryFormProps> = (props) => {
  const handleClickAddButton = () => {
    props.onUpdate([...props.inputs, {op: 'equals', value: '', id: -1}]);
  };

  const handleUpdateQueryInput = (index: number, input: Query) => {
    props.onUpdate(immutableSplice(props.inputs, index, 1, {...input, id: props.inputs[index].id}));
  };

  return (
    <>
      <Grid container={true} spacing={1} alignItems={'flex-start'} justify={'center'}>
        {props.inputs.map((input, i) =>
          <QueryInput
            key={input.id}
            op={input.op}
            value={input.value}
            onUpdate={handleUpdateQueryInput.bind(null, i)}
            onDelete={props.onDelete.bind(null, input.id)}
          />)}
        <Grid container={true} spacing={1}>
          <Grid item={true} xs={12}>
            <Button variant="outlined" color="primary" onClick={handleClickAddButton}>
              Add
            </Button>
          </Grid>
        </Grid>
      </Grid>
    </>
  );
};

interface FilterDialogProps {
  open: boolean
  onClose: () => void
  inputs: Query[]
  onClickApplyButton: (enabled: boolean, changed: boolean, inputs: Query[]) => void
}

interface QueryInWithID {
  op: QueryOp;
  value: string;
  id: number;
}

const toQueryInWithIDList = (qi: Query[]): QueryInWithID[] => {
  return qi.map((i, ind): QueryInWithID => ({...i, id: ind}));
}

// tslint:disable-next-line:variable-name
export const FilterDialog: React.FC<FilterDialogProps> = (props) => {
  const [ops, setOps] = useState(toQueryInWithIDList(props.inputs));
  const [opCnt, setOpCnt] = useState(props.inputs.length);
  const [enable, setEnable] = useState(false);
  const [lastEnableSwitchMode, setLastEnableSwitchMode] = useState(false);
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
    props.onClickApplyButton(enable, lastEnableSwitchMode !== enable, ops);
    setLastEnableSwitchMode(enable);
  };

  const handleChangeEnableSwitch = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEnable(e.target.checked);
  }

  const handleDeleteQuery = (id: number) => {
    setOps(ops.filter((o) => o.id !== id));
  };

  return (
    <Dialog open={props.open} onClose={props.onClose} aria-labelledby="form-dialog-title">
      <DialogTitle>
        Filter settings
        <Switch
          checked={enable}
          onChange={handleChangeEnableSwitch}
          color="primary"
        />
      </DialogTitle>
      <DialogContent>
        <QueryForm
          inputs={ops}
          onUpdate={handleUpdateQueryForm}
          onDelete={handleDeleteQuery}
        />
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
