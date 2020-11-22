import { Draggable, DraggableHandlers } from './draggable';
import { Pixel } from '../../svg/svg';

export class TouchDraggable<E extends SVGElement> implements Draggable {
  constructor(private element: E, private handlers: DraggableHandlers) {
    const passive = { passive: true };
    this.element.addEventListener('touchstart', this._onTouchStart, passive);
    this.element.addEventListener('touchmove', this._onTouchMove, passive);
    this.element.addEventListener('touchend', this._onTouchEnd, passive);
  }

  destroy() {
    this.element.removeEventListener('touchstart', this._onTouchStart);
    this.element.removeEventListener('touchmove', this._onTouchMove);
    this.element.removeEventListener('touchend', this._onTouchEnd);
  }

  private initialTouch?: { x: Pixel; y: Pixel };

  private readonly _onTouchStart = (e: TouchEvent) => {
    e.stopPropagation();

    // 通常ありえない
    if (!e.currentTarget || !e.target) {
      return;
    }

    // ピンチズーム とかで誤動作させない
    if (e.changedTouches.length !== 1) {
      return;
    }

    const touch = e.changedTouches[0];
    const x = touch.clientX;
    const y = touch.clientY;

    this.initialTouch = { x, y };
    this.handlers.onDragStart?.(x, y, e);
  };

  private readonly _onTouchMove = (e: TouchEvent) => {
    e.stopPropagation();

    // 通常ありえない
    if (!e.currentTarget || !e.target) {
      return;
    }

    // ピンチズーム とかで誤動作させない
    if (e.changedTouches.length !== 1) {
      return;
    }

    if (this.initialTouch === undefined) {
      return;
    }

    const { x, y } = this.initialTouch;
    const { clientX, clientY } = e.changedTouches[0];

    this.handlers.onMove?.(clientX - x, clientY - y);
  };

  private readonly _onTouchEnd = (e: TouchEvent) => {
    e.stopPropagation();
    // 通常ありえない
    if (!e.currentTarget || !e.target) {
      return;
    }

    // ピンチズーム とかで誤動作させない
    if (e.changedTouches.length !== 1) {
      return;
    }

    const touch = e.changedTouches[0];
    const x = touch.clientX;
    const y = touch.clientY;
    this.handlers.onDragEnd?.(x, y, e);
    this.initialTouch = undefined;
  };
}
