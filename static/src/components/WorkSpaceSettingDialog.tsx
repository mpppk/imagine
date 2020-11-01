import {Button} from "@material-ui/core";
import Dialog from "@material-ui/core/Dialog";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import DialogTitle from "@material-ui/core/DialogTitle";
import React, {useState} from "react";
import {WorkSpace} from "../models/models";
import TextField from "@material-ui/core/TextField";

interface WorkSpaceSettingDialogProps {
  open: boolean
  workspace: WorkSpace | null
  onClose: () => void
  onApply: (ws: WorkSpace) => void;
}

// const useStyles = makeStyles((theme) => ({
// }));

// tslint:disable-next-line:variable-name
export const WorkSpaceSettingDialog: React.FC<WorkSpaceSettingDialogProps> = (props) => {
  // const classes = useStyles();
  if (props.workspace === null) {
    return (
      <Dialog open={props.open} onClose={props.onClose} aria-labelledby="form-dialog-title">
        <DialogTitle>Workspace is not loaded</DialogTitle>
        <DialogActions>
          <Button onClick={props.onClose} color="primary">
            OK
          </Button>
        </DialogActions>
      </Dialog>
    )
  }

  const workspace = props.workspace as WorkSpace;
  const [basePath, setBasePath] = useState(workspace.basePath);

  const handleClickApplyButton = () => {
    props.onApply({...workspace, basePath});
  }

  const handleUpdateBasePath = (e: React.ChangeEvent<HTMLInputElement>) => {
    setBasePath(e.target.value);
  };

  return (
    <Dialog open={props.open} onClose={props.onClose}>
      <DialogTitle>workspace settings</DialogTitle>
      <DialogContent>
        <TextField label="base path" value={basePath} onChange={handleUpdateBasePath}/>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}>
          Cancel
        </Button>
        <Button onClick={handleClickApplyButton} color="primary">
          Apply
        </Button>
      </DialogActions>
    </Dialog>
  )
}