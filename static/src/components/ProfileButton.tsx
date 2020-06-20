import { Avatar, Theme } from '@material-ui/core';
import Button from '@material-ui/core/Button';
import { makeStyles } from '@material-ui/styles';
import React from 'react';
import { User } from '../models/models';
import ProfileMenu from './ProfileMenu';

export interface IProfileButtonProps {
  user: User;
  onClickLogout: () => void;
}

const useStyles = makeStyles((_: Theme) => ({
  avatar: {
    margin: 10
  }
}));

// tslint:disable-next-line variable-name
const ProfileButton: React.FunctionComponent<IProfileButtonProps> = props => {
  const classes = useStyles(undefined);
  const { user } = props;

  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
  const handleClickProfileButton = (e: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(e.currentTarget);
  };
  const setNullAnchorEl = () => setAnchorEl(null);
  const handleClickLogout = () => {
    setNullAnchorEl();
    props.onClickLogout();
  };

  return (
    <div>
      <Button
        aria-controls="profile-menu"
        aria-haspopup="true"
        color="inherit"
        onClick={handleClickProfileButton}
      >
        <Avatar
          aria-label='user profile avatar'
          alt="Avatar Icon"
          src={user.photoURL}
          className={classes.avatar}
        />
      </Button>
      <ProfileMenu
        anchorEl={anchorEl}
        onClickLogout={handleClickLogout}
        onClose={setNullAnchorEl}
      />
    </div>
  );
};

export default ProfileButton;
