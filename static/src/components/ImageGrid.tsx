import GridList from '@material-ui/core/GridList';
import GridListTile from '@material-ui/core/GridListTile';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';

const useStyles = makeStyles((theme) => ({
  gridList: {
    height: '100%',
    width: 240,
  },
  root: {
    backgroundColor: theme.palette.background.paper,
    display: 'flex',
    flexWrap: 'wrap',
    justifyContent: 'space-around',
    overflow: 'hidden',
  },
}));

export interface ImageGridListProps {
  paths: string[]
}

interface GridData {
  imgPath: string
  cols: number
}

const toTileData = (path: string): GridData => {
  return {
    cols: 1,
    imgPath: path,
  }
}

export function ImageGridList(props: ImageGridListProps) {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <GridList cellHeight={160} className={classes.gridList} cols={1}>
        {props.paths.map(toTileData).map((tile) => (
          <GridListTile key={tile.imgPath} cols={tile.cols || 1}>
            <img src={tile.imgPath} />
          </GridListTile>
        ))}
      </GridList>
    </div>
  );
}