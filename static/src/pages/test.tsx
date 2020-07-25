import {Button, Typography} from "@material-ui/core";
import { createStyles, makeStyles, Theme } from '@material-ui/core/styles';
import Toolbar from '@material-ui/core/Toolbar';
import React from 'react';
import {useSelector} from "react-redux";
import {indexActionCreators} from "../actions";
import {ImageListDrawer} from "../components/ImageListDrawer";
import {useActions} from "../hooks";
import {WorkSpace} from "../models/models";
import {State} from "../reducers/reducer";

const drawerWidth = 240;

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    content: {
      flexGrow: 1,
      padding: theme.spacing(3),
    },
    drawer: {
      flexShrink: 0,
      width: drawerWidth,
    },
    drawerContainer: {
      overflow: 'auto',
    },
    drawerPaper: {
      width: drawerWidth,
    },
    root: {
      display: 'flex',
    },
  }),
);

const useHandlers = () => {
  const actionCreators = useActions(indexActionCreators);
  return {
    handleAddDirectoryButton: (ws: WorkSpace) => {
      actionCreators.clickAddDirectoryButton({workSpaceName: ws.name});
    }
  };
}

interface GlobalState {
  currentWorkSpace: WorkSpace | null
  imagePaths: string[]
  isLoadingWorkSpace: boolean
  isScanningDirectories: boolean
}

const selector = (state: State): GlobalState => ({
  currentWorkSpace: state.global.currentWorkSpace,
  imagePaths: state.indexPage.imagePaths,
  isLoadingWorkSpace: state.global.isLoadingWorkSpaces,
  isScanningDirectories: state.indexPage.scanning,
})

export default function ClippedDrawer() {
  const classes = useStyles();

  const handlers = useHandlers();
  const globalState = useSelector(selector)
  // const state = useSelector((s: State) => s.indexPage);
  // const workspace = useSelector((s: State) => s.global.currentWorkSpace);
  const handleClickAddDirectoryButton = () => {
    if (globalState.currentWorkSpace === null) {
      // tslint:disable-next-line:no-console
      console.warn('workspace is not selected, but AddDirectoryButton is clicked')
      return
    }
    handlers.handleAddDirectoryButton(globalState.currentWorkSpace)
  }
  return (
    <div className={classes.root}>
      <ImageListDrawer imagePaths={globalState.imagePaths}/>
      <main className={classes.content}>
        <Toolbar />
        <Typography paragraph={true}>
          Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt
          ut labore et dolore magna aliqua. Rhoncus dolor purus non enim praesent elementum
          facilisis leo vel. Risus at ultrices mi tempus imperdiet. Semper risus in hendrerit
          gravida rutrum quisque non tellus. Convallis convallis tellus id interdum velit laoreet id
          donec ultrices. Odio morbi quis commodo odio aenean sed adipiscing. Amet nisl suscipit
          adipiscing bibendum est ultricies integer quis. Cursus euismod quis viverra nibh cras.
          Metus vulputate eu scelerisque felis imperdiet proin fermentum leo. Mauris commodo quis
          imperdiet massa tincidunt. Cras tincidunt lobortis feugiat vivamus at augue. At augue eget
          arcu dictum varius duis at consectetur lorem. Velit sed ullamcorper morbi tincidunt. Lorem
          donec massa sapien faucibus et molestie ac.
        </Typography>
        <Typography paragraph={true}>
          Consequat mauris nunc congue nisi vitae suscipit. Fringilla est ullamcorper eget nulla
          facilisi etiam dignissim diam. Pulvinar elementum integer enim neque volutpat ac
          tincidunt. Ornare suspendisse sed nisi lacus sed viverra tellus. Purus sit amet volutpat
          consequat mauris. Elementum eu facilisis sed odio morbi. Euismod lacinia at quis risus sed
          vulputate odio. Morbi tincidunt ornare massa eget egestas purus viverra accumsan in. In
          hendrerit gravida rutrum quisque non tellus orci ac. Pellentesque nec nam aliquam sem et
          tortor. Habitant morbi tristique senectus et. Adipiscing elit duis tristique sollicitudin
          nibh sit. Ornare aenean euismod elementum nisi quis eleifend. Commodo viverra maecenas
          accumsan lacus vel facilisis. Nulla posuere sollicitudin aliquam ultrices sagittis orci a.
        </Typography>
        <Button variant="outlined" color="primary"
                disabled={globalState.isScanningDirectories || globalState.isLoadingWorkSpace}
                onClick={handleClickAddDirectoryButton}>
          {globalState.isScanningDirectories ? 'Scanning...' : 'Add Directory'}
        </Button>
      </main>
    </div>
  );
}