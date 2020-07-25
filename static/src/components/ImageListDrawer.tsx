import Drawer from '@material-ui/core/Drawer';
import {makeStyles} from "@material-ui/core/styles";
import Toolbar from "@material-ui/core/Toolbar";
import React from 'react';
import {ImageGridList} from "./ImageGrid";

const drawerWidth = 240;

const useStyles = makeStyles(() => {
  return {
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
  }
});

interface ImageListDrawerProps {
  imagePaths: string[]
  onClickImage: (path: string) => void
}

// tslint:disable-next-line variable-name
export const ImageListDrawer: React.FunctionComponent<ImageListDrawerProps> = props => {
  const classes = useStyles();
  return (
    <Drawer
      open={true}
      variant="persistent"
      anchor="left"
      className={classes.drawer}
      classes={{paper: classes.drawerPaper}}
    >
      <Toolbar />
      <div className={classes.drawerContainer}>
        <ImageGridList
          paths={props.imagePaths}
          onClickImage={props.onClickImage}
        />
      </div>
    </Drawer>
  );
};

