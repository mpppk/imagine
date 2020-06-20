import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import * as React from 'react';

interface IProfileListProps {
  anchorEl: null | HTMLElement;
  onClickLogout: () => void;
  onClose: () => void;
}

// tslint:disable-next-line variable-name
const ProfileMenu: React.FunctionComponent<IProfileListProps> = props => {
  return (
    <Menu
      id="profile-menu"
      anchorEl={props.anchorEl}
      keepMounted={true}
      open={Boolean(props.anchorEl)}
      onClose={props.onClose}
    >
      <MenuItem onClick={props.onClickLogout}>Logout</MenuItem>
    </Menu>
  );
};

export default ProfileMenu;
