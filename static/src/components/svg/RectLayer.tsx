import {Layer, Pixel} from "./svg";
import React from "react";
import {useDrag} from "../util/draggable/draggable";
import {ResizeHandler} from "./ResizeHandler";
import {makeStyles} from "@material-ui/core/styles";
import {Theme} from "@material-ui/core";

const useStyles = makeStyles((theme: Theme) => {
  return {
    rect: {
      fill: 'transparent',
      stroke: theme.palette.primary.light,
      strokeWidth: 4,
      cursor: 'move'
    },
    resizeHandler: {
      stroke: theme.palette.primary.light,
    },
    deleteCircle: {
      cursor: "pointer",
    }
  }
});

interface Props extends Layer {
  onScaleStart?: (width: Pixel, height: Pixel) => void;
  onScale?: (width: Pixel, height: Pixel) => void;
  onScaleEnd?: (width: Pixel, height: Pixel) => void;
  onDelete?: () => void;
}

export function RectLayer(props: Props) {
  const classes = useStyles();
  const ref = useDrag("ontouchstart" in window, {
    onMove: props.onMove,
    onDragStart: props.onDragStart,
    onDragEnd: props.onDragEnd,
  });

  return (
    <>
      <rect
        className={classes.rect}
        fill="orange"
        width={props.width}
        height={props.height}
        x={props.x}
        y={props.y}
        ref={ref}
      />
      <ResizeHandler
        x={props.x}
        y={props.y}
        width={props.width}
        height={props.height}
        onScaleStart={props.onScaleStart}
        onScale={props.onScale}
        onScaleEnd={props.onScaleEnd}
      />
      <circle
        className={classes.deleteCircle}
        cx={props.x+props.width+20}
        cy={props.y+props.height}
        fill={'orange'}
        r={10}
        onClick={props.onDelete}
      />
    </>
  );
}
