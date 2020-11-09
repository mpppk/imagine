import {Theme} from "@material-ui/core";
import IconButton from "@material-ui/core/IconButton";
import Paper from "@material-ui/core/Paper";
import {makeStyles} from "@material-ui/core/styles";
import TextField from "@material-ui/core/TextField";
import CheckCircleIcon from "@material-ui/icons/CheckCircle";
import React, {useMemo, useState} from "react";
import {Controller, useForm} from "react-hook-form";
import {FieldErrors} from "react-hook-form/dist/types/form";
import {Tag} from "../../models/models";

const useStyles = makeStyles((theme: Theme) => {
  return {
    checkCircleButton: {
      bottom: theme.spacing(1),
      float: "right",
    },
    item: {
      margin: `0 0 ${theme.spacing(1)}px 0`,
      padding: theme.spacing(2),
      position: "relative",
      userSelect: "none",
    },
    tagNameTextField: {
      width: 100,
    },
  }
});

interface Props {
  tag: Tag
  index: number
  onFinishEdit: (tag: Tag) => void
  errorMessage?: string
}

const useHandlers = (props: Props, localState: LocalState) => {
  return useMemo(() => {
    return {
      changeTagName: (e: React.ChangeEvent<HTMLInputElement>) => {
        const tagName = e.target.value;
        localState.setCurrentTagName(tagName);
      },
      submitTagName: (data: any) => {
        props.onFinishEdit({...props.tag, name: data.tagName})
      },
      keyDown: (e: React.KeyboardEvent<HTMLFormElement>) => e.stopPropagation(),
    };
  }, [props, localState])
}

type LocalState = ReturnType<typeof useLocalState>;

const useViewState = (props: Props, state: LocalState, errors: FieldErrors) => {
  const classes = useStyles()
  return useMemo(() => {
    return {
      paper: {
        className: classes.item,
      },
      textField: {
        defaultValue: state.currentTagName,
        error: !!errors.value || !!props.errorMessage,
        helperText: errors.value?.type ?? props.errorMessage,
        value: state.currentTagName,
      }
    };
  }, [props, state, errors])
};

const useLocalState = (props: Props) => {
  const [currentTagName, setCurrentTagName] = useState(props.tag.name);
  return {
    currentTagName,
    setCurrentTagName,
  }
}

const useViewStateAndHandlers = (props: Props, errors: FieldErrors) => {
  const localState = useLocalState(props);
  return {
    handlers: useHandlers(props, localState),
    viewState: useViewState(props,localState, errors ),
  }
}

// tslint:disable-next-line:variable-name
export const EditingTagListItem: React.FC<Props> = (props) => {
  const classes = useStyles();
  const {handleSubmit, control, errors} = useForm();
  const {handlers, viewState} = useViewStateAndHandlers(props, errors);

  return (
      <Paper
        className={viewState.paper.className}
      >
        <form onSubmit={handleSubmit(handlers.submitTagName)} onKeyDown={handlers.keyDown}>
          <Controller
            as={TextField}
            name="tagName"
            rules={{required: true}}
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
            type="submit"
            aria-label="update-tag"
            className={classes.checkCircleButton}
          >
            <CheckCircleIcon/>
          </IconButton>
        </form>
      </Paper>
    );
}
