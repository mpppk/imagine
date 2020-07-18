import {NextPage} from 'next';
import React from 'react';
import {ImageGridList} from "../components/ImageGrid";
import {Button, LinearProgressProps} from "@material-ui/core";
import {useActions} from "../hooks";
import {indexActionCreators} from "../actions";
import {useSelector} from "react-redux";
import {State} from "../reducers/reducer";
import Box from "@material-ui/core/Box";
import LinearProgress from "@material-ui/core/LinearProgress";
import Typography from "@material-ui/core/Typography";

const useHandlers =  () => {
  const actionCreators = useActions(indexActionCreators);
  return {
    handleAddDirectoryButton: (workspaceName: string) => {
      actionCreators.clickAddDirectoryButton({workSpaceName: workspaceName});
    }
  };
}

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
  const handlers = useHandlers();
  const state = useSelector((s: State) => s.indexPage);
  const workspace = useSelector((s: State) => s.global.currentWorkSpace);
  const handleClickAddDirectoryButton = () => {
    handlers.handleAddDirectoryButton(workspace)
  }
  return (
    <div>
      <ImageGridList paths={state.imagePaths}/>
      <Button variant="outlined" color="primary">
        Edit Query
      </Button>
      <Button variant="outlined" color="primary" disabled={state.scanning} onClick={handleClickAddDirectoryButton}>
        {state.scanning ? 'Scanning...' : 'Add Directory'}
      </Button>
      {state.scanning ? <LinearProgressWithLabel value={50} /> : null}
    </div>
  );
};

export default Index;
