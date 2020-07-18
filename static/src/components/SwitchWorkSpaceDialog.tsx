import {Button} from "@material-ui/core";
import Dialog from "@material-ui/core/Dialog";
import React, {useState} from "react";
import DialogTitle from "@material-ui/core/DialogTitle";
import DialogContent from "@material-ui/core/DialogContent";
import DialogActions from "@material-ui/core/DialogActions";
import Select from "@material-ui/core/Select";
import InputLabel from "@material-ui/core/InputLabel";
import {makeStyles} from "@material-ui/core/styles";
import FormControl from "@material-ui/core/FormControl";
import MenuItem from "@material-ui/core/MenuItem";

interface SwitchWorkSpaceDialogProps {
  open: boolean
  workspaces: string[]
  currentWorkSpace: string
  onClose: () => void
  onSelectWorkSpace: (ws: string) => void
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

export const SwitchWorkSpaceDialog: React.FC<SwitchWorkSpaceDialogProps> = (props) => {
  const classes = useStyles();
  const [workspace, setWorkSpace] = useState(props.currentWorkSpace)
  const handleChangeSelect = (event: React.ChangeEvent<{ value: unknown }>) => {
    setWorkSpace(event.target.value as string)
  }
  const handleClickSwitchButton = () => {
    if (workspace !== props.currentWorkSpace) {
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
          value={workspace}
        >
          {props.workspaces.map(ws => <MenuItem value={ws} key={ws}>{ws}</MenuItem>)}
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