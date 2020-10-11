import {Draggable, DraggableHandlers} from "./draggable";
import {Pixel} from "../../svg/svg";

export class MouseDraggable<E extends SVGElement> implements Draggable {
  constructor(private element: E, private handlers: DraggableHandlers) {
    const passive = {passive: true};
    this.element.addEventListener("mousedown", this._onDragStart, passive);
    this.element.addEventListener("mousemove", this._onDrag, passive);
    this.element.addEventListener("mouseup", this._onDragEnd, passive);
  }

  destroy() {
    this.element.removeEventListener("dragstart", this._onDragStart);
    this.element.removeEventListener("drag", this._onDrag);
    this.element.removeEventListener("dragend", this._onDragEnd);
  }

  private initialDrag?: { x: Pixel, y: Pixel }

  private readonly _onDragStart = (e: MouseEvent) => {
    e.stopPropagation();

    // 通常ありえない
    if (!e.currentTarget || !e.target) {
      return;
    }

    this.initialDrag = {x: e.clientX, y: e.clientY};
    this.handlers.onDragStart?.(this.initialDrag.x, this.initialDrag.y, e);
  };

  private readonly _onDrag = (e: MouseEvent) => {
    e.stopPropagation();

    // 通常ありえない
    if (!e.currentTarget || !e.target) {
      return;
    }

    if (this.initialDrag === undefined) {
      return;
    }

    const {x, y} = this.initialDrag;
    this.handlers.onMove?.(e.clientX - x, e.clientY - y);
  };

  private readonly _onDragEnd = (e: MouseEvent) => {
    e.stopPropagation();
    this.handlers.onDragEnd?.(e);
    this.initialDrag = undefined;
  }
}
