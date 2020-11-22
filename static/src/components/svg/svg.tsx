import { DraggableHandlers } from '../util/draggable/draggable';

export type Pixel = number;

export interface Layer extends DraggableHandlers {
  id: number;
  width: Pixel;
  height: Pixel;
  x: Pixel;
  y: Pixel;
}
