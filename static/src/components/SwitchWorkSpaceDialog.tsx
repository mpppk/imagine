import {Button} from "@material-ui/core";
import Dialog from "@material-ui/core/Dialog";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import DialogTitle from "@material-ui/core/DialogTitle";
import FormControl from "@material-ui/core/FormControl";
import InputLabel from "@material-ui/core/InputLabel";
import MenuItem from "@material-ui/core/MenuItem";
import Select from "@material-ui/core/Select";
import {makeStyles} from "@material-ui/core/styles";
import React, {useEffect, useState} from "react";
import {WorkSpace} from "../models/models";

interface SwitchWorkSpaceDialogProps {
  open: boolean
  workspaces: WorkSpace[]
  currentWorkSpace: WorkSpace | null
  onClose: () => void
  onSelectWorkSpace: (ws: WorkSpace) => void
}

const useStyles = makeStyles((theme) => ({
  formControl: {
    margin: theme.spacing(1),
    minWidth: 120,
  },
  selectEmpty: {
    marginTop: theme.spacing(2),
  },
}));

// tslint:disable-next-line:variable-name
export const SwitchWorkSpaceDialog: React.FC<SwitchWorkSpaceDialogProps> = (props) => {
  const classes = useStyles();
  const [workspace, setWorkSpace] = useState(props.currentWorkSpace)

  useEffect(() => {
    if (props.currentWorkSpace !== undefined) {
      setWorkSpace(props.currentWorkSpace)
    }
  }, [props.currentWorkSpace])

  const handleChangeSelect = (event: React.ChangeEvent<{ value: unknown }>) => {
    const ws = props.workspaces.find(w => w.name === event.target.value)
    if (ws !== undefined) {
      setWorkSpace(ws)
    }
  }
  const handleClickSwitchButton = () => {
    if (workspace !== props.currentWorkSpace && workspace !== null) {
      props.onSelectWorkSpace(workspace)
    }
    props.onClose()
  }

  return (
    <Dialog open={props.open} onClose={props.onClose} aria-labelledby="form-dialog-title">
      <DialogTitle>Select workspace</DialogTitle>
      <DialogContent>
              <FormControl className={classes.formControl}>
        <InputLabel>work space</InputLabel>
        <Select
          onChange={handleChangeSelect}
          value={workspace ? workspace.name : ''}
        >
          {props.workspaces.map(ws => <MenuItem value={ws.name} key={ws.name}>{ws.name}</MenuItem>)}
        </Select>
      </FormControl>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} color="primary">
          Cancel
        </Button>
        <Button onClick={handleClickSwitchButton} color="primary">
          Switch
        </Button>
      </DialogActions>
    </Dialog>
  )
}