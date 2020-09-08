import Drawer from '@material-ui/core/Drawer';
import {makeStyles} from "@material-ui/core/styles";
import Toolbar from "@material-ui/core/Toolbar";
import React from 'react';
import {Tag} from "../models/models";
import {TagList} from "./TagList";

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

interface TagListDrawerProps {
  tags: Tag[]
  editTagId?: number
  selectedTagId?: number
  assignedTagIds: number[]
  onClickAddButton: () => void
  onClickEditButton: (tag: Tag) => void
  onClickDeleteButton?: (tag: Tag) => void
  onUpdate?: (newTags: Tag[]) => void
  onRename?: (tag: Tag) => void
}

// tslint:disable-next-line variable-name
export const TagListDrawer: React.FunctionComponent<TagListDrawerProps> = props => {
  const classes = useStyles();
  return (
    <Drawer
      open={true}
      variant="persistent"
      anchor="right"
      className={classes.drawer}
      classes={{
        paper: classes.drawerPaper,
      }}
    >
      <Toolbar/>
      <div className={classes.drawerContainer}>
        <TagList
          tags={props.tags}
          editTagId={props.editTagId}
          selectedTagId={props.selectedTagId}
          assignedTagIds={props.assignedTagIds}
          onClickAddButton={props.onClickAddButton}
          onClickEditButton={props.onClickEditButton}
          onRename={props.onRename}
          onUpdate={props.onUpdate}
        />
      </div>
    </Drawer>
  );
};

