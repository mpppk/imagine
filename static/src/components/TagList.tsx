import {Theme} from "@material-ui/core";
import Button from "@material-ui/core/Button";
import IconButton from "@material-ui/core/IconButton";
import List from "@material-ui/core/List";
import Paper from "@material-ui/core/Paper";
import {makeStyles} from "@material-ui/core/styles";
import AddIcon from '@material-ui/icons/Add';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import React, {useState} from "react";
import {DragDropContext, Draggable, Droppable} from "react-beautiful-dnd";
import {Tag} from "../models/models";

// a little function to help us with reordering the result
const reorder = (list: Tag[], startIndex: number, endIndex: number) => {
  const result = Array.from(list);
  const [removed] = result.splice(startIndex, 1);
  result.splice(endIndex, 0, removed);

  return result;
};

const useStyles = makeStyles((theme: Theme) => {
  return {
    addButton: {
      width: 250 - theme.spacing(2)
    },
    draggingItem: {
      background: "lightgray",
      margin: `0 0 ${theme.spacing(1)}px 0`,
      padding: theme.spacing(2),
      userSelect: "none",
    },
    draggingList: {
      background: "lightgray",
      padding: theme.spacing(1),
      width: 250
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
    list: {
      padding: theme.spacing(1),
      width: 250
    },
  }
});

export interface TagListProps {
  tags: Tag[]
  onChangeOrder: (tags: Tag[]) => void
  onClickAddButton: (tags: Tag[]) => void
}

interface TagListItemProps {
  tag: Tag
  index: number
  onClickEditButton: (tag: Tag, index: number) => void
}

// tslint:disable-next-line:variable-name
export const TagListItem: React.FC<TagListItemProps> = (props) => {
  const classes = useStyles()
  const tag = props.tag;

  const generateHandleClickEditButton = (t: Tag, index: number) => () => {
    props.onClickEditButton(t, index)
  }

  return (<Draggable key={tag.name} draggableId={tag.name} index={props.index}>
    {(provided2, snapshot2) => (
      <Paper
        ref={provided2.innerRef}
        {...provided2.draggableProps}
        {...provided2.dragHandleProps}
        className={snapshot2.isDragging ? classes.draggingItem : classes.item}
        style={{...provided2.draggableProps.style}}
      >
        {tag.name}
        <IconButton aria-label="delete" className={classes.itemButton}>
          <DeleteIcon/>
        </IconButton>
        <IconButton
          aria-label="edit"
          className={classes.itemButton}
          onClick={generateHandleClickEditButton(tag, props.index)}
        >
          <EditIcon/>
        </IconButton>
      </Paper>
    )}
  </Draggable>)
}

// tslint:disable-next-line:variable-name
export const TagList: React.FC<TagListProps> = (props) => {
  const classes = useStyles();
  const [editItemIndex, setEditItemIndex] = useState(null as number | null);

  const onDragEnd = (result: any) => {
    // dropped outside the list
    if (!result.destination) {
      return;
    }

    const newTags = reorder(
      props.tags,
      result.source.index,
      result.destination.index
    );

    props.onChangeOrder(newTags)
  }
  const handleClickAddButton = () => {
    props.onClickAddButton(props.tags);
  }

  const handleClickItemEditButton = (tag: Tag, index: number) => {
    console.log('handle edit', tag, editItemIndex);
    setEditItemIndex(index);
  }

  return (
    <div>
      <DragDropContext onDragEnd={onDragEnd}>
        <Droppable droppableId="droppable">
          {(provided, snapshot) => (
            <List
              {...provided.droppableProps}
              ref={provided.innerRef}
              component="nav"
              className={snapshot.isDraggingOver ? classes.draggingList : classes.list}
            >
              {props.tags.map((tag, index) => (
                <TagListItem key={tag.name} tag={tag} index={index}
                             onClickEditButton={handleClickItemEditButton}/>
              ))}
              {provided.placeholder}
              <Button
                variant="outlined"
                color="primary"
                className={classes.addButton}
                onClick={handleClickAddButton}
              >
                <AddIcon/>
              </Button>
            </List>
          )}
        </Droppable>
      </DragDropContext>
    </div>
  )
}
