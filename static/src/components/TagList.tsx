import {Theme} from "@material-ui/core";
import Button from "@material-ui/core/Button";
import IconButton from "@material-ui/core/IconButton";
import List from "@material-ui/core/List";
import Paper from "@material-ui/core/Paper";
import {makeStyles} from "@material-ui/core/styles";
import TextField from "@material-ui/core/TextField";
import AddIcon from '@material-ui/icons/Add';
import CheckCircleIcon from '@material-ui/icons/CheckCircle';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import React, {useState} from "react";
import {DragDropContext, Draggable, Droppable} from "react-beautiful-dnd";
import {Controller, useForm} from "react-hook-form";
import {Tag} from "../models/models";
import {immutableSplice} from "../util";

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
    checkCircleButton: {
      bottom: theme.spacing(1),
      float: "right",
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

export interface TagListProps {
  tags: Tag[]
  editTagId?: number
  onClickAddButton: (tags: Tag[]) => void
  onClickEditButton: (tag: Tag) => void
  onClickDeleteButton?: (tag: Tag) => void
  onUpdate: (tags: Tag[]) => void
  onRename?: (tag: Tag) => void
}

interface TagListItemProps {
  tag: Tag
  index: number
  onClickEditButton: (tag: Tag) => void
  onClickDeleteButton: (tag: Tag) => void
}

interface EditingTagListItemProps {
  tag: Tag
  index: number
  onFinishEdit: (tag: Tag) => void
  errorMessage?: string
}

// tslint:disable-next-line:variable-name
export const EditingTagListItem: React.FC<EditingTagListItemProps> = (props) => {
  const classes = useStyles()
  const tag = props.tag;
  const [currentTagName, setCurrentTagName] = useState(tag.name);
  const {handleSubmit, control, errors} = useForm();

  const handleSubmitTagName = (data: any) => {
    props.onFinishEdit({...tag, name: data.tagName})
  }

  const handleChangeTagName = (e: React.ChangeEvent<HTMLInputElement>) => {
    const tagName = e.target.value;
    setCurrentTagName(tagName);
  }

  return (<Draggable key={tag.id} draggableId={tag.id.toString()} index={props.index}>
    {(provided2, snapshot2) => (
      <Paper
        ref={provided2.innerRef}
        {...provided2.draggableProps}
        {...provided2.dragHandleProps}
        className={snapshot2.isDragging ? classes.draggingItem : classes.item}
        style={{...provided2.draggableProps.style}}
      >
        <form onSubmit={handleSubmit(handleSubmitTagName)}>
          <Controller
            as={TextField}
            name="tagName"
            rules={{required: true}}
            control={control}

            className={classes.tagNameTextField}
            value={currentTagName}
            defaultValue={currentTagName}
            autoFocus={true}
            error={errors.tagName!! || props.errorMessage !== undefined}
            helperText={errors.tagName?.type ?? props.errorMessage}
            onChange={handleChangeTagName}
          />
          <IconButton
            type="submit"
            aria-label="update-tag"
            className={classes.checkCircleButton}
          >
            <CheckCircleIcon/>
          </IconButton>
        </form>
      </Paper>
    )}
  </Draggable>)
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
        <IconButton
          aria-label="delete"
          className={classes.itemButton}
          onClick={genClickDeleteButtonHandler(tag)}
        >
          <DeleteIcon/>
        </IconButton>
        <IconButton
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
        <Droppable droppableId="droppable">
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
