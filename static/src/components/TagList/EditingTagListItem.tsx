import { Theme } from '@material-ui/core';
import IconButton from '@material-ui/core/IconButton';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import CheckCircleIcon from '@material-ui/icons/CheckCircle';
import React, { useMemo, useState } from 'react';
import { Controller, useForm } from 'react-hook-form';
import { FieldErrors } from 'react-hook-form/dist/types/form';
import { Tag } from '../../models/models';

const useStyles = makeStyles((theme: Theme) => {
  // FIXME: duplicated code
  const baseItemStyles = {
    marginRight: theme.spacing(1),
    marginBottom: theme.spacing(1),
    marginTop: theme.spacing(1),
    paddingLeft: theme.spacing(1),
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    alignContent: 'center',
  };
  return {
    item: {
      ...baseItemStyles,
      position: 'relative',
      userSelect: 'none',
    },
    tagNameTextField: {
      width: 140,
    },
  };
});

interface Props {
  tag: Tag;
  index: number;
  onFinishEdit: (tag: Tag) => void;
  errorMessage?: string;
}

const useHandlers = (props: Props, localState: LocalState) => {
  return useMemo(() => {
    return {
      changeTagName: (e: React.ChangeEvent<HTMLInputElement>) => {
        const tagName = e.target.value;
        localState.setCurrentTagName(tagName);
      },
      submitTagName: (data: any) => {
        props.onFinishEdit({ ...props.tag, name: data.tagName });
      },
      keyDown: (e: React.KeyboardEvent<HTMLFormElement>) => e.stopPropagation(),
    };
  }, [props, localState]);
};

type LocalState = ReturnType<typeof useLocalState>;

const useViewState = (props: Props, state: LocalState, errors: FieldErrors) => {
  return useMemo(() => {
    return {
      textField: {
        defaultValue: state.currentTagName,
        error: !!errors.value || !!props.errorMessage,
        helperText: errors.value?.type ?? props.errorMessage,
        value: state.currentTagName,
      },
    };
  }, [props, state, errors]);
};

const useLocalState = (props: Props) => {
  const [currentTagName, setCurrentTagName] = useState(props.tag.name);
  return {
    currentTagName,
    setCurrentTagName,
  };
};

const useViewStateAndHandlers = (props: Props, errors: FieldErrors) => {
  const localState = useLocalState(props);
  return {
    handlers: useHandlers(props, localState),
    viewState: useViewState(props, localState, errors),
  };
};

// tslint:disable-next-line:variable-name
export const EditingTagListItem: React.FC<Props> = (props) => {
  const classes = useStyles();
  const { handleSubmit, control, errors } = useForm();
  const { handlers, viewState } = useViewStateAndHandlers(props, errors);

  return (
    <Paper>
      <form
        onSubmit={handleSubmit(handlers.submitTagName)}
        onKeyDown={handlers.keyDown}
      >
        <div className={classes.item}>
          <Controller
            data-cy="tag-name-form"
            as={TextField}
            name="tagName"
            rules={{ required: true }}
            control={control}
            className={classes.tagNameTextField}
            value={viewState.textField.value}
            defaultValue={viewState.textField.defaultValue}
            autoFocus={true}
            error={viewState.textField.error}
            helperText={viewState.textField.helperText}
            onChange={handlers.changeTagName}
          />
          <IconButton
            data-cy="save-tag-name-button"
            type="submit"
            aria-label="update-tag"
          >
            <CheckCircleIcon />
          </IconButton>
        </div>
      </form>
    </Paper>
  );
};
