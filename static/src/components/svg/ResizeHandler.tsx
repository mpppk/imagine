import { Pixel } from './svg';
import { useDrag } from '../util/draggable/draggable';
import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { Theme } from '@material-ui/core';

interface Props {
  x: Pixel;
  y: Pixel;
  width: Pixel;
  height: Pixel;
  onScaleStart?(x: Pixel, y: Pixel): void;
  onScale?(x: Pixel, y: Pixel): void;
  onScaleEnd?(x: Pixel, y: Pixel): void;
}

const HANDLE_SIZE = 10 as Pixel;

/**
 * 実際のリサイズハンドラよりもどのくらい当たり判定を大きくするか
 */
const TOLERANCE = 4 as Pixel;

const useStyles = makeStyles((theme: Theme) => {
  return {
    resizeHandler: {
      stroke: theme.palette.primary.light,
      fill: theme.palette.primary.light,
    },
  };
});

// tslint:disable-next-line:variable-name
export const ResizeHandler: React.FC<Props> = (props) => {
  const classes = useStyles();
  const x = props.x + props.width - HANDLE_SIZE / 2;
  const y = props.y + props.height - HANDLE_SIZE / 2;

  const ref = useDrag('ontouchstart' in window, {
    onDragStart: props.onScaleStart,
    onMove: props.onScale,
    onDragEnd: props.onScaleEnd,
  });

  return (
    <g>
      <rect
        className={classes.resizeHandler}
        width={HANDLE_SIZE}
        height={HANDLE_SIZE}
        x={x}
        y={y}
      />

      {/** 上に透明な当たり判定を大きめにかぶせる */}
      <rect
        ref={ref}
        fillOpacity="0"
        width={HANDLE_SIZE + TOLERANCE * 2}
        height={HANDLE_SIZE + TOLERANCE * 2}
        x={x - TOLERANCE}
        y={y - TOLERANCE}
        style={{ cursor: 'pointer' }}
      />
    </g>
  );
};
