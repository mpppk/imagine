import AppBar from '@material-ui/core/AppBar/AppBar';
import Button from "@material-ui/core/Button";
import IconButton from '@material-ui/core/IconButton/IconButton';
import {Theme} from '@material-ui/core/styles';
import Toolbar from '@material-ui/core/Toolbar/Toolbar';
import Typography from '@material-ui/core/Typography/Typography';
import MenuIcon from '@material-ui/icons/Menu';
import {makeStyles} from '@material-ui/styles';
import * as React from 'react';
import {useState} from 'react';
import {useDispatch, useSelector} from "react-redux";
import {globalActionCreators} from "../actions/global";
import {WorkSpace} from "../models/models";
import {State} from "../reducers/reducer";
import MyDrawer from './drawer/Drawer';
import {SwitchWorkSpaceDialog} from "./SwitchWorkSpaceDialog";

const useStyles = makeStyles((theme: Theme) => ({
  menuButton: {
    marginRight: theme.spacing(2)
  },
  root: {
    flexGrow: 1
  },
  title: {
    flexGrow: 1
  },
}));

// tslint:disable-next-line variable-name
export function MyAppBar() {
  const classes = useStyles(undefined);
  const [isDrawerOpen, setDrawerOpen] = useState(false);
  const handleDrawer = (open: boolean) => () => setDrawerOpen(open);
  const [isWorkSpaceDialogOpen, setWorkSpaceDialogOpen] = useState(false);
  const currentWorkSpace = useSelector((s: State) => s.global.currentWorkSpace)
  const workspaces = useSelector((s: State) => s.global.workspaces)
  const isLoadingWorkSpaces = useSelector((s: State) => s.global.isLoadingWorkSpaces)
  const dispatch = useDispatch();

  const handleClickOpenSwitchWorkSpaceDialogButton = () => {
    setWorkSpaceDialogOpen(true)
  }

  const handleCloseSwitchWorkSpaceDialog = () => {
    setWorkSpaceDialogOpen(false)
  }

  const handleSelectWorkSpace = (ws: WorkSpace) => {
    dispatch(globalActionCreators.selectNewWorkSpace(ws));
  }

  return (
    <div className={classes.root}>
      <AppBar position="static">
        <MyDrawer
          open={isDrawerOpen}
          onClose={handleDrawer(false)}
          onClickSideList={handleDrawer(false)}
        />
        <Toolbar>
          <IconButton
            edge="start"
            className={classes.menuButton}
            color="inherit"
            aria-label="menu"
            onClick={handleDrawer(true)}
          >
            <MenuIcon/>
          </IconButton>
          <Typography variant="h6" className={classes.title}>
            {currentWorkSpace === null ? 'loading workspace...' : currentWorkSpace.name}
          </Typography>
          <Button color="inherit" disabled={isLoadingWorkSpaces}
                  onClick={handleClickOpenSwitchWorkSpaceDialogButton}>
            Switch WorkSpace
          </Button>
        </Toolbar>
      </AppBar>
      <SwitchWorkSpaceDialog
        workspaces={workspaces === null ? [] : workspaces}
        open={isWorkSpaceDialogOpen}
        onClose={handleCloseSwitchWorkSpaceDialog}
        currentWorkSpace={currentWorkSpace}
        onSelectWorkSpace={handleSelectWorkSpace}
      />
    </div>
  );
}
