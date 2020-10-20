import {Layer, Pixel} from "./svg";
import React from "react";
import {useDrag} from "../util/draggable/draggable";
import {ResizeHandler} from "./ResizeHandler";

interface Props {
  src: Layer;
  className?: string;
  onScale: (width: Pixel, height: Pixel) => void;
}

// const useResizeObserver = (el: Element | null, callback: ResizeObserverCallback) => {
//   // const [width, setWidth] = useState(0);
//   // const [height, setHeight] = useState(0);
//
//   useEffect(() => {
//     const resizeObserver = new ResizeObserver(callback);
//     if (el !== null) {
//       resizeObserver.observe(el);
//     }
//   }, [el]);
// };

export function RectLayer({src, className, onScale}: Props) {
  const ref = useDrag("ontouchstart" in window, {
    onMove: src.onMove,
    onDragStart: src.onDragStart,
    onDragEnd: src.onDragEnd,
  });

  const {width, height} = ref.current === null ?
    {width: 0, height: 0} :
    ref.current.getBoundingClientRect();

  const handleScale = (dx: Pixel, dy: Pixel) => {
    onScale(Math.max(width + dx, 0), Math.max(height + dy, 0))
  };

  return (
    <>
      <rect
        className={className}
        fill="orange"
        width={src.width}
        height={src.height}
        x={src.x}
        y={src.y}
        ref={ref}
      />
      <ResizeHandler width={width} height={height} onScale={handleScale}/>
    </>
  );
}
