import AppBar from '@material-ui/core/AppBar/AppBar';
import Button from '@material-ui/core/Button';
import IconButton from '@material-ui/core/IconButton/IconButton';
import { Theme } from '@material-ui/core/styles';
import Toolbar from '@material-ui/core/Toolbar/Toolbar';
import Typography from '@material-ui/core/Typography/Typography';
import MenuIcon from '@material-ui/icons/Menu';
import { makeStyles } from '@material-ui/styles';
import * as React from 'react';
import { useState } from 'react';
import { useSelector } from 'react-redux';
import { workspaceActionCreators } from '../actions/workspace';
import { Query } from '../models/models';
import { State } from '../reducers/reducer';
import MyDrawer from './drawer/Drawer';
import { SwitchWorkSpaceDialog } from './SwitchWorkSpaceDialog';
import { FilterButton } from './FilterButton';
import { FilterDialog } from './FilterDialog';
import { indexActionCreators } from '../actions';
import { useActions } from '../hooks';
import { WorkSpaceSettingDialog } from './WorkSpaceSettingDialog';
import { Tooltip } from '@material-ui/core';
import SettingsIcon from '@material-ui/icons/Settings';
import { fsActionCreators } from '../actions/fs';

const useStyles = makeStyles((theme: Theme) => ({
  appBar: {
    zIndex: theme.zIndex.drawer + 1,
  },
  menuButton: {
    marginRight: theme.spacing(2),
  },
  title: {
    flexGrow: 1,
  },
  workspaceSettingButton: {
    marginLeft: theme.spacing(-1), // FIXME
  },
  filterButton: {
    marginLeft: theme.spacing(0), // FIXME
  },
}));

const useLocalState = () => {
  const [isDrawerOpen, setDrawerOpen] = useState(false);
  const [isSwitchWorkSpaceDialogOpen, setSwitchWorkSpaceDialogOpen] = useState(
    false
  );
  const [openFilterDialog, setOpenFilterDialog] = useState(false);
  const [openWorkSpaceSettingDialog, setWorkSpaceSettingDialog] = useState(
    false
  );
  return {
    isDrawerOpen,
    setDrawerOpen,
    isSwitchWorkSpaceDialogOpen,
    setSwitchWorkSpaceDialogOpen,
    openFilterDialog,
    setOpenFilterDialog,
    openWorkSpaceSettingDialog,
    setWorkSpaceSettingDialog,
  };
};

type LocalState = ReturnType<typeof useLocalState>;

const useViewState = (localState: LocalState) => {
  const viewState = {
    ...localState,
    disableChangeBasePathButton: useSelector(
      (s: State) => s.indexPage.scanning || s.global.isLoadingWorkSpaces
    ),
    currentWorkSpace: useSelector((s: State) => s.global.currentWorkSpace),
    workspaces: useSelector((s: State) => s.global.workspaces),
    isLoadingWorkSpaces: useSelector(
      (s: State) => s.global.isLoadingWorkSpaces
    ),
    isFiltered: useSelector(
      (s: State) => s.global.queries.length > 0 && s.global.filterEnabled
    ),
    scanStatus: useSelector((s: State) => ({
      running: s.indexPage.scanning,
      count: s.indexPage.scanCount,
    })),
  };

  let scanMsg =
    viewState.scanStatus.count === 0
      ? null
      : `🔎 ${viewState.scanStatus.count}`;
  if (viewState.scanStatus.count !== 0 && !viewState.scanStatus.running) {
    scanMsg = `✔ ${viewState.scanStatus.count}`;
  }

  return {
    ...viewState,
    scanMsg,
  };
};

type ViewState = ReturnType<typeof useViewState>;

const useHandlers = (localState: LocalState, viewState: ViewState) => {
  const indexActionDispatcher = useActions(indexActionCreators);
  const workspaceActionDispatcher = useActions(workspaceActionCreators);
  const fsActionDispatcher = useActions(fsActionCreators);

  return {
    genDrawerHandler: (open: boolean) => () => localState.setDrawerOpen(open),
    clickOpenSwitchWorkSpaceDialogButton: localState.setSwitchWorkSpaceDialogOpen.bind(
      null,
      true
    ),
    closeSwitchWorkSpaceDialog: localState.setSwitchWorkSpaceDialogOpen.bind(
      null,
      false
    ),
    selectWorkSpace: workspaceActionDispatcher.select,
    clickFilterButton: localState.setOpenFilterDialog.bind(null, true),
    closeFilter: localState.setOpenFilterDialog.bind(null, false),

    clickFilterApplyButton: (
      enabled: boolean,
      changed: boolean,
      queries: Query[]
    ) => {
      localState.setOpenFilterDialog(false);
      indexActionDispatcher.clickFilterApplyButton({
        enabled,
        changed,
        queries: enabled ? queries : [],
      });
    },

    closeWorkSpaceSetting: localState.setWorkSpaceSettingDialog.bind(
      null,
      false
    ),
    clickWorkSpaceSettingButton: localState.setWorkSpaceSettingDialog.bind(
      null,
      true
    ),
    clickChangeBasePathButton: (needToLoadAssets: boolean) => {
      if (viewState.currentWorkSpace === null) {
        // tslint:disable-next-line:no-console
        console.warn(
          'workspace is not selected, but AddDirectoryButton is clicked'
        );
        return;
      }
      indexActionDispatcher.clickChangeBasePathButton({
        needToLoadAssets,
        workSpaceName: viewState.currentWorkSpace.name,
      });
    },
  };
};

// tslint:disable-next-line variable-name
export function MyAppBar() {
  const classes = useStyles(undefined);
  const localState = useLocalState();
  const viewState = useViewState(localState);
  const handlers = useHandlers(localState, viewState);

  return (
    <>
      <AppBar position="fixed" className={classes.appBar}>
        <MyDrawer
          open={viewState.isDrawerOpen}
          onClose={handlers.genDrawerHandler(false)}
          onClickSideList={handlers.genDrawerHandler(false)}
        />
        <Toolbar>
          <IconButton
            edge="start"
            className={classes.menuButton}
            color="inherit"
            aria-label="menu"
            onClick={handlers.genDrawerHandler(true)}
          >
            <MenuIcon />
          </IconButton>
          <Typography variant="h6" className={classes.title}>
            <span>
              {viewState.currentWorkSpace === null
                ? 'loading workspace...'
                : viewState.currentWorkSpace.name}
            </span>
            <FilterButton
              className={classes.filterButton}
              onClick={handlers.clickFilterButton}
              dot={viewState.isFiltered}
            />
            <Tooltip
              title="Edit WorkSpace settings"
              aria-label="edit-workspace-settings"
              className={classes.workspaceSettingButton}
            >
              <IconButton
                edge="start"
                color="inherit"
                aria-label="workspace-setting"
                onClick={handlers.clickWorkSpaceSettingButton}
              >
                <SettingsIcon />
              </IconButton>
            </Tooltip>
          </Typography>
          {viewState.scanMsg}
          <Button
            color="inherit"
            disabled={viewState.isLoadingWorkSpaces}
            onClick={handlers.clickOpenSwitchWorkSpaceDialogButton}
          >
            Switch WorkSpace
          </Button>
        </Toolbar>
      </AppBar>
      <SwitchWorkSpaceDialog
        workspaces={viewState.workspaces === null ? [] : viewState.workspaces}
        open={viewState.isSwitchWorkSpaceDialogOpen}
        onClose={handlers.closeSwitchWorkSpaceDialog}
        currentWorkSpace={viewState.currentWorkSpace}
        onSelectWorkSpace={handlers.selectWorkSpace}
      />
      <FilterDialog
        onClickApplyButton={handlers.clickFilterApplyButton}
        open={viewState.openFilterDialog}
        onClose={handlers.closeFilter}
        inputs={[]}
      />
      <WorkSpaceSettingDialog
        onClose={handlers.closeWorkSpaceSetting}
        open={viewState.openWorkSpaceSettingDialog}
        workspace={viewState.currentWorkSpace}
        disableChangeBasePathButton={viewState.disableChangeBasePathButton}
        onClickChangeBaseDirButton={handlers.clickChangeBasePathButton}
      />
    </>
  );
}
