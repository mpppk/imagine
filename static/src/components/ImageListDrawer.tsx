import Drawer from '@material-ui/core/Drawer';
import {makeStyles} from "@material-ui/core/styles";
import Toolbar from "@material-ui/core/Toolbar";
import React from 'react';
import {ImageGridList} from "./ImageGrid";
import {VirtualizedAssetProps} from "../services/virtualizedAsset";

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

interface ImageListDrawerProps extends VirtualizedAssetProps {
  imagePaths: string[]
  onClickImage: (path: string, index: number) => void
  selectedIndex: number
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
          {...props}
          paths={props.imagePaths}
          onClickImage={props.onClickImage}
          cellHeight={200}
          selectedIndex={props.selectedIndex}
          width={drawerWidth}
        />
      </div>
    </Drawer>
  );
};

