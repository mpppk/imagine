import GridList from '@material-ui/core/GridList';
import GridListTile from '@material-ui/core/GridListTile';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';

const useStyles = makeStyles((theme) => ({
  gridList: {
    height: '100%',
  },
  gridListTile: {
    cursor: 'pointer',
  },
  root: {
    backgroundColor: theme.palette.background.paper,
    display: 'flex',
    flexWrap: 'wrap',
    justifyContent: 'space-around',
    overflow: 'hidden',
  },
}));

interface Props {
  paths: string[]
  onClickImage: (path: string) => void
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

export function ImageGridList(props: Props) {
  const classes = useStyles();

  const genClickImageHandler = (imgPath: string) => () => {
    props.onClickImage(imgPath);
  };

  return (
    <div className={classes.root}>
      <GridList cellHeight={200} className={classes.gridList} cols={1}>
        {props.paths.map(toTileData).map((tile) => (
          <GridListTile
            key={tile.imgPath}
            cols={tile.cols || 1}
            onClick={genClickImageHandler(tile.imgPath)}
            className={classes.gridListTile}
          >
            <img src={tile.imgPath} />
          </GridListTile>
        ))}
      </GridList>
    </div>
  );
}