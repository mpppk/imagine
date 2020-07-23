import {Theme} from "@material-ui/core";
import Button from "@material-ui/core/Button";
import List from "@material-ui/core/List";
import {makeStyles} from "@material-ui/core/styles";
import AddIcon from '@material-ui/icons/Add';
import React, {useState} from "react";
import {DragDropContext, Droppable} from "react-beautiful-dnd";
import {Tag} from "../models/models";
import {immutableSplice, reorder} from "../util";
import {EditingTagListItem} from "./TagList/EditingTagListItem";
import {TagListItem} from "./TagList/TagListItem";

const useStyles = makeStyles((theme: Theme) => {
  return {
    addButton: {
      width: 250 - theme.spacing(2)
    },
    draggingList: {
      background: "lightgray",
      padding: theme.spacing(1),
      width: 250
    },
    list: {
      padding: theme.spacing(1),
      width: 250
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
