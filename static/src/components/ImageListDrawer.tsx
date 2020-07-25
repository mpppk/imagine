import Drawer from '@material-ui/core/Drawer';
import {makeStyles} from "@material-ui/core/styles";
import Toolbar from "@material-ui/core/Toolbar";
import React from 'react';
import {ImageGridList} from "./ImageGrid";

const useStyles = makeStyles(() => {
  return {
    drawer: {
      flexShrink: 0,
      width: 240,
    },
    drawerContainer: {
      overflow: 'auto',
    },
  }
});

interface ImageListDrawerProps {
  imagePaths: string[]
}

// tslint:disable-next-line variable-name
export const ImageListDrawer: React.FunctionComponent<ImageListDrawerProps> = props => {
  const classes = useStyles();
  return (
    <Drawer open={true} variant="persistent" anchor="left" className={classes.drawer}>
      <Toolbar />
      <div className={classes.drawerContainer}>
        <ImageGridList paths={props.imagePaths}/>
      </div>
    </Drawer>
  );
};

