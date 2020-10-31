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
import {workspaceActionCreators} from "../actions/workspace";
import {QueryInput, WorkSpace} from "../models/models";
import {State} from "../reducers/reducer";
import MyDrawer from './drawer/Drawer';
import {SwitchWorkSpaceDialog} from "./SwitchWorkSpaceDialog";
import {FilterButton} from "./FilterButton";
import {FilterDialog} from "./FilterDialog";
import {indexActionCreators} from "../actions";
import {useActions} from "../hooks";
import {WorkSpaceSettingDialog} from "./WorkSpaceSettingDialog";

const useStyles = makeStyles((theme: Theme) => ({
  appBar: {
    zIndex: theme.zIndex.drawer + 1,
  },
  menuButton: {
    marginRight: theme.spacing(2)
  },
  title: {
    flexGrow: 1,
    cursor: "pointer",
  },
}));

// tslint:disable-next-line variable-name
export function MyAppBar() {
  const classes = useStyles(undefined);
  const [isDrawerOpen, setDrawerOpen] = useState(false);
  const handleDrawer = (open: boolean) => () => setDrawerOpen(open);
  const [isSwitchWorkSpaceDialogOpen, setSwitchWorkSpaceDialogOpen] = useState(false);
  const [openFilterDialog, setOpenFilterDialog] = useState(false);
  const [openWorkSpaceSettingDialog, setWorkSpaceSettingDialog] = useState(false);
  const currentWorkSpace = useSelector((s: State) => s.global.currentWorkSpace)
  const workspaces = useSelector((s: State) => s.global.workspaces)
  const isLoadingWorkSpaces = useSelector((s: State) => s.global.isLoadingWorkSpaces)
  const isFiltered = useSelector((s: State) => {
    return s.global.queries.length > 0 && s.global.filterEnabled;
  });
  const dispatch = useDispatch();
  const indexActionDispatcher = useActions(indexActionCreators);
  const workspaceActionDispatcher = useActions(workspaceActionCreators);

  const handleClickOpenSwitchWorkSpaceDialogButton = () => {
    setSwitchWorkSpaceDialogOpen(true)
  }

  const handleCloseSwitchWorkSpaceDialog = () => {
    setSwitchWorkSpaceDialogOpen(false)
  }

  const handleSelectWorkSpace = (ws: WorkSpace) => {
    dispatch(workspaceActionCreators.select(ws));
  }

  const handleClickFilterButton = () => {
    setOpenFilterDialog(true);
  };

  const handleCloseFilter = () => {
    setOpenFilterDialog(false);
  }

  const handleClickFilterApplyButton = (enabled: boolean, changed: boolean, queryInputs: QueryInput[]) => {
    setOpenFilterDialog(false);
    indexActionDispatcher.clickFilterApplyButton({enabled, changed, queryInputs});
  };

  const handleClickWorkSpaceName = () => {
    setWorkSpaceSettingDialog(true);
  };

  const handleApplyWorkSpaceSetting = (workspace: WorkSpace) => {
    workspaceActionDispatcher.updateRequest(workspace);
    setWorkSpaceSettingDialog(false);
  }

  const handleCloseWorkSpaceSetting = () => {
    setWorkSpaceSettingDialog(false);
  };

  return (
    <>
      <AppBar position="fixed" className={classes.appBar}>
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
          <FilterButton onClick={handleClickFilterButton} dot={isFiltered}/>
            <Typography
              variant="h6"
              className={classes.title}
              onClick={handleClickWorkSpaceName}
            >
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
        open={isSwitchWorkSpaceDialogOpen}
        onClose={handleCloseSwitchWorkSpaceDialog}
        currentWorkSpace={currentWorkSpace}
        onSelectWorkSpace={handleSelectWorkSpace}
      />
      <FilterDialog
        onClickApplyButton={handleClickFilterApplyButton}
        open={openFilterDialog}
        onClose={handleCloseFilter}
        inputs={[]}
      />
      <WorkSpaceSettingDialog
       onApply={handleApplyWorkSpaceSetting}
       onClose={handleCloseWorkSpaceSetting}
       open={openWorkSpaceSettingDialog}
       workspace={currentWorkSpace}/>
    </>
  );
}
