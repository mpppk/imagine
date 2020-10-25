import {Layer, Pixel} from "./svg";
import React from "react";
import {useDrag} from "../util/draggable/draggable";
import {ResizeHandler} from "./ResizeHandler";

interface Props extends Layer {
  className?: string;
  onScaleStart?: (width: Pixel, height: Pixel) => void;
  onScale?: (width: Pixel, height: Pixel) => void;
  onScaleEnd?: (width: Pixel, height: Pixel) => void;
}

export function RectLayer(props: Props) {
  const ref = useDrag("ontouchstart" in window, {
    onMove: props.onMove,
    onDragStart: props.onDragStart,
    onDragEnd: props.onDragEnd,
  });

  return (
    <>
      <rect
        className={props.className}
        fill="orange"
        width={props.width}
        height={props.height}
        x={props.x}
        y={props.y}
        ref={ref}
      />
      <ResizeHandler
        width={props.width}
        height={props.height}
        onScaleStart={props.onScaleStart}
        onScale={props.onScale}
        onScaleEnd={props.onScaleEnd}
      />
    </>
  );
}
