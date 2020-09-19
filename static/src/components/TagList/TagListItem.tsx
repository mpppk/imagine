import {Theme} from "@material-ui/core";
import IconButton from "@material-ui/core/IconButton";
import Paper from "@material-ui/core/Paper";
import {makeStyles} from "@material-ui/core/styles";
import DeleteIcon from "@material-ui/icons/Delete";
import EditIcon from "@material-ui/icons/Edit";
import React, {useMemo} from "react";
import {Draggable} from "react-beautiful-dnd";
import {Tag} from "../../models/models";

const useStyles = makeStyles((theme: Theme) => {
  return {
    disabledItem: {
      margin: `0 0 ${theme.spacing(1)}px 0`,
      padding: theme.spacing(2),
      position: "relative",
      userSelect: "none",
    },
    draggingItem: {
      background: theme.palette.action.selected,
      margin: `0 0 ${theme.spacing(1)}px 0`,
      padding: theme.spacing(2),
      userSelect: "none",
    },
    item: {
      margin: `0 0 ${theme.spacing(1)}px 0`,
      padding: theme.spacing(2),
      position: "relative",
      userSelect: "none",
    },
    itemButton: {
      bottom: theme.spacing(2),
      float: "right",
    },
    assignedItem: {
      borderLeft: `thick solid ${theme.palette.primary.light}`,
    },
    selectedItem: {
      fontWeight: theme.typography.fontWeightBold,
    }
  }
});

interface Props {
  tag: Tag
  index: number
  selected: boolean
  assigned: boolean
  disabled?: boolean
  onClick: (tag: Tag) => void
  onClickEditButton: (tag: Tag) => void
  onClickDeleteButton: (tag: Tag) => void
}

const useHandlers = (props: Props) => {
  return useMemo(() => {
    return {
      click: () => {
        props.onClick(props.tag)
      },
      clickDeleteButton: () => {
        props.onClickDeleteButton(props.tag)
      },
      clickEditButton: () => {
        props.onClickEditButton(props.tag)
      },
    };
  }, [props])
};

const useViewState = (props: Props) => {
  const classes = useStyles();
  return useMemo(() => {
    const paperClassName = props.disabled ? classes.disabledItem : classes.item;
    return {
      paper: {
        genClassNames: (isDragging: boolean, assigned: boolean, selected: boolean) => {
          const ret = [isDragging ? classes.draggingItem : paperClassName];
          if (assigned) {
            ret.push(classes.assignedItem);
          }
          if (selected) {
            ret.push(classes.selectedItem);
          }
          return ret;
        },
      },
    };
  }, [props, classes]);
}

// tslint:disable-next-line:variable-name
export const TagListItem: React.FC<Props> = (props) => {
  const classes = useStyles()
  const handlers = useHandlers(props);
  const viewState = useViewState(props);

  const genPaperClassName = (isDragging: boolean) => {
    return viewState.paper.genClassNames(isDragging, props.assigned, props.selected).join(' ');
  }

  return (<Draggable key={props.tag.name} draggableId={props.tag.name} index={props.index} isDragDisabled={props.disabled}>
    {(provided, snapshot) => (
      <Paper
        onClick={handlers.click}
        elevation={props.selected ? 4 : 1}
        ref={provided.innerRef}
        {...provided.draggableProps}
        {...provided.dragHandleProps}
        className={genPaperClassName(snapshot.isDragging)}
        style={{...provided.draggableProps.style}}
      >
        {props.tag.name}
        <IconButton
          disabled={props.disabled}
          aria-label="delete"
          className={classes.itemButton}
          onClick={handlers.clickDeleteButton}
        >
          <DeleteIcon/>
        </IconButton>
        <IconButton
          disabled={props.disabled}
          aria-label="edit"
          className={classes.itemButton}
          onClick={handlers.clickEditButton}
        >
          <EditIcon/>
        </IconButton>
      </Paper>
    )}
  </Draggable>)
}

