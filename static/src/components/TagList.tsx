import { Theme } from '@material-ui/core';
import Button from '@material-ui/core/Button';
import List from '@material-ui/core/List';
import { makeStyles } from '@material-ui/core/styles';
import AddIcon from '@material-ui/icons/Add';
import React, { useMemo, useState } from 'react';
import { DragDropContext, Droppable } from 'react-beautiful-dnd';
import { Tag } from '../models/models';
import { immutableSplice, isDupNamedTag, reorder } from '../util';
import { EditingTagListItem } from './TagList/EditingTagListItem';
import { TagListItem } from './TagList/TagListItem';

const useStyles = makeStyles((theme: Theme) => {
  return {
    addButton: {
      width: 250 - theme.spacing(3),
    },
    draggingList: {
      padding: theme.spacing(1),
      width: 250,
    },
    list: {
      padding: theme.spacing(1),
      width: 250,
    },
  };
});

interface Props {
  tags: Tag[];
  editTagId?: number;
  selectedTagId?: number;
  assignedTagIds: number[];
  onClick: (tag: Tag) => void;
  onClickEditButton: (tag: Tag) => void;
  onClickDeleteButton?: (tag: Tag) => void;
  onUpdate?: (tags: Tag[]) => void;
  onRename?: (tag: Tag) => void;
}

const useLocalState = () => {
  const [tagNameDuplicatedError, setTagNameDuplicatedError] = useState(false);
  const [showNewTagForm, setShowNewTagForm] = useState(false);
  return {
    // tslint:disable-next-line:object-literal-sort-keys
    showNewTagForm,
    tagNameDuplicatedError,
    setShowNewTagForm,
    setTagNameDuplicatedError,
  };
};

type LocalState = ReturnType<typeof useLocalState>;

const useHandlers = (props: Props, localState: LocalState) => {
  return useMemo(() => {
    return {
      clickAddButton: () => {
        localState.setShowNewTagForm(true);
      },

      clickItem: (tag: Tag) => {
        props.onClick(tag);
      },

      clickItemEditButton: (tag: Tag) => {
        props.onClickEditButton(tag);
      },

      clickItemDeleteButton: (tag: Tag) => {
        const index = props.tags.findIndex((t) => t.id === tag.id);
        const newTags = immutableSplice(props.tags, index, 1);
        props.onUpdate?.(newTags);
        props.onClickDeleteButton?.(tag);
      },

      dragEnd: (result: any) => {
        // dropped outside the list
        if (!result.destination) {
          return;
        }

        const newTags = reorder(
          props.tags,
          result.source.index,
          result.destination.index
        );

        props.onUpdate?.(newTags);
      },

      finishItemEdit: (tag: Tag) => {
        const isDupName = isDupNamedTag(props.tags, tag);
        localState.setTagNameDuplicatedError(isDupName);
        if (isDupName) {
          return;
        }
        if (localState.showNewTagForm) {
          localState.setShowNewTagForm(false);
          props.onUpdate?.([tag, ...props.tags]);
          props.onRename?.(tag);
          return;
        }
        const index = props.tags.findIndex((t) => t.id === tag.id);
        const newTags = immutableSplice(props.tags, index, 1, tag);
        props.onUpdate?.(newTags);
        props.onRename?.(tag);
      },
    };
  }, [props, localState]);
};

const useViewState = (localState: LocalState) => {
  return useMemo(() => {
    return {
      editingTagErrorMessage: localState.tagNameDuplicatedError
        ? 'name duplicated'
        : undefined,
    };
  }, [localState]);
};

// tslint:disable-next-line:variable-name
export const TagList: React.FC<Props> = (props) => {
  const classes = useStyles();
  const localState = useLocalState();
  const viewState = useViewState(localState);
  const handlers = useHandlers(props, localState);

  const baseIndex = localState.showNewTagForm ? 1 : 0;

  return (
    <div data-cy="tag-list">
      <DragDropContext onDragEnd={handlers.dragEnd}>
        <Droppable droppableId="droppable" isDropDisabled={!!props.editTagId}>
          {(provided, snapshot) => (
            <List
              {...provided.droppableProps}
              ref={provided.innerRef}
              component="nav"
              className={
                snapshot.isDraggingOver ? classes.draggingList : classes.list
              }
            >
              <Button
                data-cy="add-new-tag-button"
                variant="outlined"
                color="primary"
                disabled={!!props.editTagId}
                className={classes.addButton}
                onClick={handlers.clickAddButton}
              >
                <AddIcon />
              </Button>
              {localState.showNewTagForm ? (
                <EditingTagListItem
                  key={'new-tag'}
                  tag={{ id: 0, name: '' }}
                  errorMessage={viewState.editingTagErrorMessage}
                  index={0}
                  onFinishEdit={handlers.finishItemEdit}
                />
              ) : null}
              {props.tags.map((tag, index) =>
                props.editTagId === tag.id ? (
                  <EditingTagListItem
                    key={tag.id}
                    tag={tag}
                    errorMessage={viewState.editingTagErrorMessage}
                    index={baseIndex + index}
                    onFinishEdit={handlers.finishItemEdit}
                  />
                ) : (
                  <TagListItem
                    disabled={!!props.editTagId}
                    selected={tag.id === props.selectedTagId}
                    assigned={props.assignedTagIds.includes(tag.id)}
                    key={tag.id}
                    tag={tag}
                    index={index}
                    onClick={handlers.clickItem}
                    onClickEditButton={handlers.clickItemEditButton}
                    onClickDeleteButton={handlers.clickItemDeleteButton}
                  />
                )
              )}
              {provided.placeholder}
            </List>
          )}
        </Droppable>
      </DragDropContext>
    </div>
  );
};
