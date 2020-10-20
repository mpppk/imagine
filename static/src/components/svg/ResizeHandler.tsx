import {Pixel} from "./svg";
import {useDrag} from "../util/draggable/draggable";

interface Props {
  width: Pixel
  height: Pixel
  onScale(x: Pixel, y: Pixel): void;
}

const HANDLE_SIZE = 10 as Pixel;

/**
 * 実際のリサイズハンドラよりもどのくらい当たり判定を大きくするか
 */
const TOLERANCE = 4 as Pixel;

export function ResizeHandler({
  width,
  height,
  onScale,
}: Props) {
  const ref = useDrag("ontouchstart" in window, {
    onMove: onScale,
  });

  const x = width - HANDLE_SIZE / 2;
  const y = height - HANDLE_SIZE / 2;

  return (
    <g>
      <rect
        fill="green"
        stroke="#666666"
        strokeWidth="1"
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
        style={{ cursor: "pointer" }}
      />
    </g>
  );
}
