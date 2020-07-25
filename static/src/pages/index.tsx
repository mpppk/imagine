import {Button, LinearProgressProps, Theme} from "@material-ui/core";
import Box from "@material-ui/core/Box";
import Grid from "@material-ui/core/Grid";
import LinearProgress from "@material-ui/core/LinearProgress";
import {makeStyles} from "@material-ui/core/styles";
import Typography from "@material-ui/core/Typography";
import {NextPage} from 'next';
import React from 'react';
import {useSelector} from "react-redux";
import {indexActionCreators} from "../actions";
import {ImageGridList} from "../components/ImageGrid";
import {useActions} from "../hooks";
import {WorkSpace} from "../models/models";
import {State} from "../reducers/reducer";

const useStyles = makeStyles((theme: Theme) => {
  return {
    content: {
      flexGrow: 1,
      padding: theme.spacing(3),
    },
    root: {
      display: 'flex',
    },
  }
});

const useHandlers = () => {
  const actionCreators = useActions(indexActionCreators);
  return {
    addDirectoryButton: (ws: WorkSpace) => {
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

function LinearProgressWithLabel(props: LinearProgressProps & { value: number }) {
  return (
    <Box display="flex" alignItems="center">
      <Box width="100%" mr={1}>
        <LinearProgress variant="determinate" {...props} />
      </Box>
      <Box minWidth={35}>
        <Typography variant="body2" color="textSecondary">{`${Math.round(
          props.value,
        )}%`}</Typography>
      </Box>
    </Box>
  );
}

// tslint:disable-next-line variable-name
export const Index: NextPage = () => {
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
    handlers.addDirectoryButton(globalState.currentWorkSpace)
  }
  return (
    <div className={classes.root}>
      <div className={classes.content}>
        <Grid container={true} spacing={3}>
          <Grid item={true} xs={3}>
            <ImageGridList paths={globalState.imagePaths}/>
          </Grid>
        </Grid>
        <Button variant="outlined" color="primary" disabled={globalState.isLoadingWorkSpace}>
          Edit Query
        </Button>
        <Button variant="outlined" color="primary"
                disabled={globalState.isScanningDirectories || globalState.isLoadingWorkSpace}
                onClick={handleClickAddDirectoryButton}>
          {globalState.isScanningDirectories ? 'Scanning...' : 'Add Directory'}
        </Button>
        {globalState.isScanningDirectories ? <LinearProgressWithLabel value={50}/> : null}
      </div>
    </div>
  );
};

export default Index;
