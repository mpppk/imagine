import {NextPage} from 'next';
import React from 'react';
import {ImageGridList} from "../components/ImageGrid";
import {Button} from "@material-ui/core";
import {useActions} from "../hooks";
import {indexActionCreators} from "../actions";

const useHandlers =  () => {
  const actionCreators = useActions(indexActionCreators);
  return {
    handleAddDirectoryButton: () => {
      actionCreators.clickAddDirectoryButton();
    }
  };
}

// tslint:disable-next-line variable-name
export const Index: NextPage = () => {
  const handlers = useHandlers();
  return (
    <div>
      <ImageGridList/>
      <Button variant="outlined" color="primary">
        Edit Query
      </Button>
      <Button variant="outlined" color="primary" onClick={handlers.handleAddDirectoryButton}>
        Add Directory
      </Button>
    </div>
  );
};

export default Index;
