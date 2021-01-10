// tslint:disable-next-line:variable-name
import React from 'react';
import { Tooltip } from '@material-ui/core';
import IconButton from '@material-ui/core/IconButton/IconButton';
import SettingsIcon from '@material-ui/icons/Settings';
import { makeStyles } from '@material-ui/styles';
import { Theme } from '@material-ui/core/styles';

const useStyles = makeStyles((theme: Theme) => ({
  button: {
    marginLeft: theme.spacing(-1), // FIXME
  },
}));

interface Props {
  onClick?: () => void;
}

// tslint:disable-next-line:variable-name
export const WorkSpaceSettingsButton: React.FC<Props> = (props) => {
  const classes = useStyles();
  return (
    <Tooltip
      title="Edit WorkSpace settings"
      aria-label="edit-workspace-settings"
      className={classes.button}
    >
      <IconButton
        data-cy="workspace-settings-button"
        edge="start"
        color="inherit"
        aria-label="workspace-setting"
        onClick={props.onClick}
      >
        <SettingsIcon />
      </IconButton>
    </Tooltip>
  );
};
