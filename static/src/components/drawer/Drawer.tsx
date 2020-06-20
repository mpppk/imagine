import Drawer from '@material-ui/core/Drawer';
import { withStyles } from '@material-ui/core/styles';
import React from 'react';
import SideList from './SideList';

const styles = {};

export interface IMyDrawerProps {
  open: boolean;
  onClose: () => void;
  onClickSideList: () => void;
}

// tslint:disable-next-line variable-name
const MyDrawer: React.FunctionComponent<IMyDrawerProps> = props => {
  const { open, onClose, onClickSideList } = props;

  return (
    <div>
      <Drawer open={open} onClose={onClose}>
        <div
          tabIndex={0}
          role="button"
          onClick={onClickSideList}
          onKeyDown={onClickSideList}
        >
          <SideList />
        </div>
      </Drawer>
    </div>
  );
};

export default withStyles(styles)(MyDrawer);
