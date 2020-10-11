import {Layer} from "./svg";
import React from "react";
import {useDrag} from "../util/draggable/draggable";

interface Props {
  src: Layer;
  className?: string;
}

export function RectLayer({ src, className }: Props) {
  const ref = useDrag("ontouchstart" in window, {
    onMove: src.onMove,
    onDragStart: src.onDragStart,
    onDragEnd: src.onDragEnd
  });

  return (
    <rect
      className={className}
      fill="orange"
      width={src.width}
      height={src.height}
      x={src.x}
      y={src.y}
      ref={ref}
    />
  );
}
