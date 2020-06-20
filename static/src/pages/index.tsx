import {NextPage} from 'next';
import React from 'react';
import {ImageGridList} from "../components/ImageGrid";
import {Button} from "@material-ui/core";

// tslint:disable-next-line variable-name
export const Index: NextPage = () => {
  return (
    <div>
      <ImageGridList/>
      <Button variant="outlined" color="primary">
        Edit Query
      </Button>
    </div>
  );
};

export default Index;
