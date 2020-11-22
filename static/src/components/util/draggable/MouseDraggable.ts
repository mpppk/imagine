import { Draggable, DraggableHandlers } from './draggable';
import { Pixel } from '../../svg/svg';

export class MouseDraggable<E extends SVGElement> implements Draggable {
  constructor(private element: E, private handlers: DraggableHandlers) {
    const passive = { passive: true };
    this.element.addEventListener('mousedown', this._onDragStart, passive);
    window.addEventListener('mousemove', this._onDrag, passive);
    window.addEventListener('mouseup', this._onDragEnd, passive);
  }

  destroy() {
    this.element.removeEventListener('mousedown', this._onDragStart);
    window.removeEventListener('mousemove', this._onDrag);
    window.removeEventListener('mouseup', this._onDragEnd);
  }

  private initialDrag?: { x: Pixel; y: Pixel };

  private readonly _onDragStart = (e: MouseEvent) => {
    e.stopPropagation();

    // 通常ありえない
    if (!e.currentTarget || !e.target) {
      return;
    }

    this.initialDrag = { x: e.clientX, y: e.clientY };
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

    const { x, y } = this.initialDrag;
    this.handlers.onMove?.(e.clientX - x, e.clientY - y);
  };

  private readonly _onDragEnd = (e: MouseEvent) => {
    e.stopPropagation();
    // 通常ありえない
    if (!e.currentTarget || !e.target) {
      return;
    }

    if (this.initialDrag === undefined) {
      return;
    }
    const { x, y } = this.initialDrag;
    this.handlers.onDragEnd?.(e.clientX - x, e.clientY - y, e);
    this.initialDrag = undefined;
  };
}
