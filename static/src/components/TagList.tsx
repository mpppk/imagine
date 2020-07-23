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
import {immutableSplice, reorder} from "../util";
import {EditingTagListItem} from "./TagList/EditingTagListItem";

const useStyles = makeStyles((theme: Theme) => {
  return {
    addButton: {
      width: 250 - theme.spacing(2)
    },
    checkCircleButton: {
      bottom: theme.spacing(1),
      float: "right",
    },
    disabledItem: {
      color: 'gray',
      margin: `0 0 ${theme.spacing(1)}px 0`,
      padding: theme.spacing(2),
      position: "relative",
      userSelect: "none",
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
    tagNameTextField: {
      width: 100,
    },
  }
});

interface TagListItemProps {
  tag: Tag
  index: number
  disabled?: boolean
  onClickEditButton: (tag: Tag) => void
  onClickDeleteButton: (tag: Tag) => void
}

// tslint:disable-next-line:variable-name
export const TagListItem: React.FC<TagListItemProps> = (props) => {
  const classes = useStyles()
  const tag = props.tag;

  const genClickEditButtonHandler = (t: Tag) => () => {
    props.onClickEditButton(t)
  }

  const genClickDeleteButtonHandler = (t: Tag) => () => {
    props.onClickDeleteButton(t)
  }

  const paperClassName = props.disabled ? classes.disabledItem : classes.item;

  return (<Draggable key={tag.name} draggableId={tag.name} index={props.index} isDragDisabled={props.disabled}>
    {(provided2, snapshot2) => (
      <Paper
        ref={provided2.innerRef}
        {...provided2.draggableProps}
        {...provided2.dragHandleProps}
        className={snapshot2.isDragging ? classes.draggingItem : paperClassName}
        style={{...provided2.draggableProps.style}}
      >
        {tag.name}
        <IconButton
          disabled={props.disabled}
          aria-label="delete"
          className={classes.itemButton}
          onClick={genClickDeleteButtonHandler(tag)}
        >
          <DeleteIcon/>
        </IconButton>
        <IconButton
          disabled={props.disabled}
          aria-label="edit"
          className={classes.itemButton}
          onClick={genClickEditButtonHandler(tag)}
        >
          <EditIcon/>
        </IconButton>
      </Paper>
    )}
  </Draggable>)
}

export interface TagListProps {
  tags: Tag[]
  editTagId?: number
  onClickAddButton: (tags: Tag[]) => void
  onClickEditButton: (tag: Tag) => void
  onClickDeleteButton?: (tag: Tag) => void
  onUpdate: (tags: Tag[]) => void
  onRename?: (tag: Tag) => void
}

// tslint:disable-next-line:variable-name
export const TagList: React.FC<TagListProps> = (props) => {
  const classes = useStyles();
  const [tagNameDuplicatedError, setTagNameDuplicatedError] = useState(false)

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

    props.onUpdate(newTags)
  }
  const handleClickAddButton = () => {
    props.onClickAddButton(props.tags);
  }

  const handleClickItemEditButton = (tag: Tag) => {
    props.onClickEditButton(tag);
  }

  const handleClickItemDeleteButton = (tag: Tag) => {
    const index = props.tags.findIndex((t) => t.id === tag.id);
    const newTags = immutableSplice(props.tags, index, 1);
    props.onUpdate(newTags);
    props.onClickDeleteButton?.(tag);
  }

  const genFinishItemEditHandler = (tag: Tag) => {
    const tagNameSet = props.tags.reduce((m, t) => {
      if (tag.id !== t.id) {
        m.add(t.name);
      }
      return m;
    }, new Set<string>())
    const isDupName = tagNameSet.has(tag.name)
    setTagNameDuplicatedError(isDupName);
    if (isDupName) {
      return
    }
    const index = props.tags.findIndex((t) => t.id === tag.id);
    const newTags = immutableSplice(props.tags, index, 1, tag);
    props.onUpdate(newTags);
    props.onRename?.(tag);
  }

  return (
    <div>
      <DragDropContext onDragEnd={onDragEnd}>
        <Droppable droppableId="droppable" isDropDisabled={!!props.editTagId}>
          {(provided, snapshot) => (
            <List
              {...provided.droppableProps}
              ref={provided.innerRef}
              component="nav"
              className={snapshot.isDraggingOver ? classes.draggingList : classes.list}
            >
              <Button
                variant="outlined"
                color="primary"
                disabled={!!props.editTagId}
                className={classes.addButton}
                onClick={handleClickAddButton}
              >
                <AddIcon/>
              </Button>
              {props.tags.map((tag, index) => (
                props.editTagId === tag.id ?
                  <EditingTagListItem
                    key={tag.id}
                    tag={tag}
                    errorMessage={tagNameDuplicatedError ? 'name duplicated' : undefined}
                    index={index}
                    onFinishEdit={genFinishItemEditHandler}
                  /> :
                  <TagListItem
                    disabled={!!props.editTagId}
                    key={tag.id}
                    tag={tag}
                    index={index}
                    onClickEditButton={handleClickItemEditButton}
                    onClickDeleteButton={handleClickItemDeleteButton}
                  />
              ))}
              {provided.placeholder}
            </List>
          )}
        </Droppable>
      </DragDropContext>
    </div>
  )
}
