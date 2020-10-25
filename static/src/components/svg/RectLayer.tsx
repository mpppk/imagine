import {Layer, Pixel} from "./svg";
import React from "react";
import {useDrag} from "../util/draggable/draggable";
import {ResizeHandler} from "./ResizeHandler";

interface Props extends Layer {
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

export function RectLayer(props: Props) {
  const ref = useDrag("ontouchstart" in window, {
    onMove: props.onMove,
    onDragStart: props.onDragStart,
    onDragEnd: props.onDragEnd,
  });

  const {width, height} = ref.current === null ?
    {width: 0, height: 0} :
    ref.current.getBoundingClientRect();

  // const handleScale = (dx: Pixel, dy: Pixel) => {
  //   props.onScale(Math.max(props.width + dx, 0), Math.max(props.height + dy, 0))
    // onScale(Math.max(width + dx, 0), Math.max(height + dy, 0))
  // };

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
      <ResizeHandler width={width} height={height} onScale={props.onScale}/>
    </>
  );
}
