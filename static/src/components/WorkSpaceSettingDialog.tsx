import {Button, Typography} from "@material-ui/core";
import Dialog from "@material-ui/core/Dialog";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import DialogTitle from "@material-ui/core/DialogTitle";
import React from "react";
import {WorkSpace} from "../models/models";
import {makeStyles} from "@material-ui/core/styles";

interface WorkSpaceSettingDialogProps {
  open: boolean
  workspace: WorkSpace | null
  onClose: () => void
  onClickChangeBasePathButton: () => void;
  disableChangeBasePathButton: boolean;
}

const useStyles = makeStyles((theme) => ({
  basePathContainer: {
    marginBottom: theme.spacing(1)
  }
}));

// tslint:disable-next-line:variable-name
export const WorkSpaceSettingDialog: React.FC<WorkSpaceSettingDialogProps> = (props) => {
  const classes = useStyles();
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

  return (
    <Dialog open={props.open} onClose={props.onClose}>
      <DialogTitle>workspace settings</DialogTitle>
      <DialogContent>
        <div className={classes.basePathContainer}>
          <Typography variant={'subtitle2'}>
            Base Path
          </Typography>
          <Typography variant={'body2'}>
            {props.workspace.basePath}
          </Typography>
        </div>
        <Button variant="outlined" color="primary"
                disabled={props.disableChangeBasePathButton}
                onClick={props.onClickChangeBasePathButton}>
          Change
        </Button>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}>
          Close
        </Button>
      </DialogActions>
    </Dialog>
  )
}