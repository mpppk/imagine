import { Theme } from '@material-ui/core';
import IconButton from '@material-ui/core/IconButton';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import React, { useMemo } from 'react';
import { Draggable } from 'react-beautiful-dnd';
import { Tag } from '../../models/models';

const useStyles = makeStyles((theme: Theme) => {
  const baseItemStyles = {
    marginRight: theme.spacing(1),
    marginBottom: theme.spacing(1),
    marginTop: theme.spacing(1),
    padding: theme.spacing(1),
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    alignContent: 'center',
  };

  return {
    tagName: {
      overflowWrap: 'break-word',
      maxWidth: 140,
    },
    draggingItem: {
      ...baseItemStyles,
      background: theme.palette.action.selected,
    },
    item: { ...baseItemStyles },
    buttonContainer: {
      display: 'inline-flex',
    },
    itemButton: {
      padding: theme.spacing(0.5),
    },
    assignedItem: {
      borderLeft: `thick solid ${theme.palette.primary.light}`,
    },
    selectedItem: {
      fontWeight: theme.typography.fontWeightBold,
    },
  };
});

interface Props {
  tag: Tag;
  index: number;
  selected: boolean;
  assigned: boolean;
  disabled?: boolean;
  onClick: (tag: Tag) => void;
  onClickEditButton: (tag: Tag) => void;
  onClickDeleteButton: (tag: Tag) => void;
}

const useHandlers = (props: Props) => {
  return useMemo(() => {
    return {
      click: () => {
        props.onClick(props.tag);
      },
      clickDeleteButton: (e: React.MouseEvent) => {
        e.stopPropagation();
        props.onClickDeleteButton(props.tag);
      },
      clickEditButton: (e: React.MouseEvent) => {
        e.stopPropagation();
        props.onClickEditButton(props.tag);
      },
    };
  }, [props]);
};

const useViewState = (props: Props) => {
  const classes = useStyles();
  return useMemo(() => {
    return {
      paper: {
        genClassNames: (
          isDragging: boolean,
          assigned: boolean,
          selected: boolean
        ) => {
          const ret = [isDragging ? classes.draggingItem : classes.item];
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
};

// tslint:disable-next-line:variable-name
export const TagListItem: React.FC<Props> = (props) => {
  const classes = useStyles();
  const handlers = useHandlers(props);
  const viewState = useViewState(props);

  const genPaperClassName = (isDragging: boolean) => {
    return viewState.paper
      .genClassNames(isDragging, props.assigned, props.selected)
      .join(' ');
  };

  let tagPrefix = '';
  if (props.index < 9) {
    tagPrefix = (props.index + 1).toString() + ': ';
  } else if (props.index === 9) {
    tagPrefix = '0: ';
  } else if (props.index === 10) {
    tagPrefix = '-: ';
  } else if (props.index === 11) {
    tagPrefix = '^: ';
  } else if (props.index === 12) {
    tagPrefix = '¥: ';
  }

  return (
    <Draggable
      key={props.tag.name}
      draggableId={props.tag.name}
      index={props.index}
      isDragDisabled={props.disabled}
    >
      {(provided, snapshot) => (
        <Paper
          data-cy="tag-list-item"
          onClick={handlers.click}
          elevation={props.selected ? 4 : 1}
          ref={provided.innerRef}
          {...provided.draggableProps}
          {...provided.dragHandleProps}
          className={genPaperClassName(snapshot.isDragging)}
          style={{ ...provided.draggableProps.style }}
        >
          <div className={classes.tagName}>{tagPrefix + props.tag.name}</div>
          <div className={classes.buttonContainer}>
            <IconButton
              data-cy="delete-tag-button"
              disabled={props.disabled}
              aria-label="delete"
              className={classes.itemButton}
              onClick={handlers.clickDeleteButton}
            >
              <DeleteIcon />
            </IconButton>
            <IconButton
              data-cy="edit-tag-button"
              disabled={props.disabled}
              aria-label="edit"
              className={classes.itemButton}
              onClick={handlers.clickEditButton}
            >
              <EditIcon />
            </IconButton>
          </div>
        </Paper>
      )}
    </Draggable>
  );
};
