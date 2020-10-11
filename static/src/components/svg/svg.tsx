import {DraggableHandlers} from "../util/draggable/draggable";

export type Pixel = number;

export interface Layer extends DraggableHandlers {
  id: number;
  width: Pixel | string;
  height: Pixel | string;
  x: Pixel;
  y: Pixel;
}

