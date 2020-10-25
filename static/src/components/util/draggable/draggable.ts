import {useEffect, useRef} from "react";
import {MouseDraggable} from "./MouseDraggable";
import {TouchDraggable} from "./TouchDraggable";
import {Pixel} from "../../svg/svg";

export interface Draggable {
  destroy(): void;
}

export type DraggableMoveHandler = (dx: Pixel, dy: Pixel) => void;
export type DraggableDragStartHandler = (x: Pixel, y: Pixel, e: Event) => void;
export type DraggableDragEndHandler = (dx: Pixel, dy:Pixel, e: Event) => void;

export interface DraggableHandlers {
  onMove?: DraggableMoveHandler
  onDragStart?: DraggableDragStartHandler
  onDragEnd?: DraggableDragEndHandler
}

export const makeDraggable = (el: SVGElement, isTouchDevice: boolean, handlers: DraggableHandlers) => {
  if (isTouchDevice) {
    return new TouchDraggable(el, handlers);
  } else {
    return new MouseDraggable(el, handlers);
  }
}

export function useDrag(isTouchDevice: boolean, handlers: DraggableHandlers) {
  const ref = useRef<SVGRectElement | null>(null);

  useEffect(() => {
    const draggable = makeDraggable(ref.current!, isTouchDevice, handlers);

    return () => {
      draggable.destroy();
    };
  }, []);

  return ref;
}
