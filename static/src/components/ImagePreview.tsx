import React, {Reducer, useEffect, useReducer, useRef} from "react";
import {AssetWithIndex, BoundingBox} from "../models/models";
import {RectLayer} from "./svg/RectLayer";
import {Pixel} from "./svg/svg";
import {Action} from "typescript-fsa";

interface Props {
  src: string
  asset: AssetWithIndex
  onMoveBoundingBox: (boxID: number, dx: Pixel, dy: Pixel) => void
  onScaleBoundingBox: (boxID: number, dx: Pixel, dy: Pixel) => void
  onDeleteBoundingBox: (boxID: number) => void;
}
interface BoxProps extends Omit<BoundingBox, 'tag'> {
  onScale: (width: Pixel, height: Pixel) => void;
  onMove: (x: Pixel, y: Pixel) => void;
  onDelete: () => void;
}

type BoxState = {
  x: Pixel,
  y: Pixel,
  w: Pixel,
  h: Pixel,
  initX: Pixel,
  initY: Pixel,
  initW: Pixel,
  initH: Pixel,
};

const newBoxState = (x: Pixel, y: Pixel, w: Pixel, h: Pixel): BoxState => ({
  x, y, w, h,
  initX: x,
  initY: y,
  initW: w,
  initH: h,
});

type BoxPayload = { dx: Pixel, dy: Pixel };

interface MoveAction extends Action<BoxPayload> {
  type: 'move'
}

interface ScaleAction extends Action<BoxPayload> {
  type: 'scale'
}

interface MoveEndAction extends Action<BoxPayload> {
  type: 'moveEnd'
}

interface ScaleEndAction extends Action<BoxPayload> {
  type: 'scaleEnd'
}

type BoxAction = MoveAction | ScaleAction | MoveEndAction | ScaleEndAction;
const boxReducer: Reducer<BoxState, BoxAction> = (state, action) => {
  switch (action.type) {
    case 'move':
      const x = Math.max(Math.min(state.initX + action.payload.dx, 500), 0);
      const y = Math.max(Math.min(state.initY + action.payload.dy, 500), 0);
      return {...state, x, y}
    case 'scale':
      const w = Math.max(Math.min(state.initW + action.payload.dx, 500), 0);
      const h = Math.max(Math.min(state.initH + action.payload.dy, 500), 0);
      return {...state, w, h}
    case 'moveEnd':
      const initX = Math.max(Math.min(state.initX + action.payload.dx, 500), 0);
      const initY = Math.max(Math.min(state.initY + action.payload.dy, 500), 0);
      return {...state, initX, initY}
    case 'scaleEnd':
      const initW = Math.max(Math.min(state.initW + action.payload.dx, 500), 0);
      const initH = Math.max(Math.min(state.initH + action.payload.dy, 500), 0);
      return {...state, initW, initH}
    default:
      const _: never = action;
      return _;
  }
}

// tslint:disable-next-line:variable-name
const Box: React.FC<BoxProps> = (props) => {
  const [state, dispatch] = useReducer(boxReducer, newBoxState(props.x, props.y, props.width, props.height));

  useEffect(props.onMove.bind(null, state.x, state.y), [state.x, state.y])
  useEffect(props.onScale.bind(null, state.w, state.h), [state.w, state.h])

  const handleDragEnd = (dx: Pixel, dy: Pixel) => {
    dispatch({type: 'moveEnd', payload: {dx, dy}})
  };

  const handleScaleEnd = (dx: Pixel, dy: Pixel) => {
    dispatch({type: 'scaleEnd', payload: {dx, dy}})
  };

  const handleScale = (dx: Pixel, dy: Pixel) => {
    dispatch({type: 'scale', payload: {dx, dy}})
  }
  const handleMove = (dx: Pixel, dy: Pixel) => {
    dispatch({type: 'move', payload: {dx, dy}})
  }
  return <RectLayer
    key={props.id}
    onScaleEnd={handleScaleEnd}
    onScale={handleScale}
    onDragEnd={handleDragEnd}
    onMove={handleMove}
    onDelete={props.onDelete}
    height={props.height}
    id={props.id}
    width={props.width}
    x={props.x}
    y={props.y}
  />
}

// tslint:disable-next-line:variable-name
export const ImagePreview: React.FC<Props> = (props) => {
  const imageRef = useRef<SVGImageElement>(null);
  const boxes = props.asset.boundingBoxes ?? [];

  return (
    <div>
      <svg id="canvas" viewBox="0 0 500 500" width="500" height="500">
        <image href={props.src} width={'100%'} height={'100%'} ref={imageRef}/>
        {boxes.map((box) => {
          return <Box
            key={box.id}
            onScale={props.onScaleBoundingBox.bind(null, box.id)}
            onMove={props.onMoveBoundingBox.bind(null, box.id)}
            onDelete={props.onDeleteBoundingBox.bind(null, box.id)}
            id={box.id}
            x={box.x}
            y={box.y}
            width={box.width}
            height={box.height}/>
        })}
      </svg>
    </div>
  )
}
