import React from 'react';
import { Tooltip } from '@material-ui/core';
import IconButton from '@material-ui/core/IconButton/IconButton';
import Badge from '@material-ui/core/Badge';
import FilterListIcon from '@material-ui/icons/FilterList';

interface Props {
  className?: string;
  onClick?: () => void;
  dot: boolean;
}

// tslint:disable-next-line:variable-name
const FilterListIconWithBadge: React.FC = () => {
  return (
    <Badge
      variant="dot"
      overlap="circle"
      color="error"
      anchorOrigin={{
        vertical: 'bottom',
        horizontal: 'right',
      }}
    >
      <FilterListIcon />
    </Badge>
  );
};

// tslint:disable-next-line:variable-name
export const FilterButton: React.FC<Props> = (props) => {
  return (
    <Tooltip
      className={props.className}
      title="Filter images"
      aria-label="filter-images"
    >
      <IconButton
        edge="start"
        color="inherit"
        aria-label="filter"
        onClick={props.onClick}
      >
        {props.dot ? <FilterListIconWithBadge /> : <FilterListIcon />}
      </IconButton>
    </Tooltip>
  );
};
