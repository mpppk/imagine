import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import GridList from '@material-ui/core/GridList';
import GridListTile from '@material-ui/core/GridListTile';

const useStyles = makeStyles((theme) => ({
  root: {
    display: 'flex',
    flexWrap: 'wrap',
    justifyContent: 'space-around',
    overflow: 'hidden',
    backgroundColor: theme.palette.background.paper,
  },
  gridList: {
    width: 500,
    height: 450,
  },
}));

/**
 * The example data is structured as follows:
 *
 * import image from 'path/to/image.jpg';
 * [etc...]
 *
 * const tileData = [
 *   {
 *     img: image,
 *     title: 'Image',
 *     author: 'author',
 *     cols: 2,
 *   },
 *   {
 *     [etc...]
 *   },
 * ];
 */

// const tileData = [
//   {
//     img: 'https://i.gyazo.com/f71cffd2e4f237030e7f6c745ce3eeeb.png',
//     title: 'everest',
//     author: 'mpppk',
//     key: 1,
//     cols: 1,
//   },
//   {
//     img: 'https://i.gyazo.com/f71cffd2e4f237030e7f6c745ce3eeeb.png',
//     title: 'everest',
//     author: 'mpppk',
//     key: 2,
//     cols: 2,
//   }
// ];

export interface ImageGridListProps {
  paths: string[]
}

interface GridData {
  imgPath: string
  cols: number
}

const toTileData = (path: string): GridData => {
  return {
    imgPath: path,
    cols: 1,
  }
}

export function ImageGridList(props: ImageGridListProps) {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <GridList cellHeight={160} className={classes.gridList} cols={3}>
        {props.paths.map(toTileData).map((tile) => (
          <GridListTile key={tile.imgPath} cols={tile.cols || 1}>
            <img src={tile.imgPath} />
          </GridListTile>
        ))}
      </GridList>
    </div>
  );
}