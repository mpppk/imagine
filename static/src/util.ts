import { Asset, BoundingBox, Direction, Tag } from './models/models';

export const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

export const immutableSplice = <T>(
  array: T[],
  start: number,
  deleteCount: number,
  ...item: T[]
): T[] => {
  return [
    ...array.slice(0, start),
    ...item,
    ...array.slice(start + deleteCount),
  ];
};

// a little function to help us with reordering the result
export const reorder = (list: Tag[], startIndex: number, endIndex: number) => {
  const result = Array.from(list);
  const [removed] = result.splice(startIndex, 1);
  result.splice(endIndex, 0, removed);

  return result;
};

export const isDupNamedTag = (tags: Tag[], newTag: Tag) => {
  const tagNameSet = tags.reduce((m, t) => {
    if (newTag.id !== t.id) {
      m.add(t.name);
    }
    return m;
  }, new Set<string>());
  return tagNameSet.has(newTag.name);
};

export const assetPathToUrl = (p: string) => `http://localhost:1323/static${p}`;

export const replaceBy = <T>(array: T[], newElm: T, f: (v: T) => boolean) => {
  const newArray = [] as T[];
  for (const v of array) {
    newArray.push(f(v) ? newElm : v);
  }
  return newArray;
};

export const isDefaultBox = (box: BoundingBox): boolean => {
  return box.height === 0 && box.width === 0 && box.x === 0 && box.y === 0;
};

export const isArrowKeyCode = (keyCode: number): boolean =>
  keyCode >= 37 && keyCode <= 40;

export const keyCodeToDirection = (keyCode: number): Direction => {
  if (!isArrowKeyCode(keyCode)) {
    throw new Error(
      'failed to convert keycode to direction. keycode: ' + keyCode
    );
  }
  switch (keyCode) {
    case 37:
      return 'LEFT';
    case 38:
      return 'UP';
    case 39:
      return 'RIGHT';
    default:
      return 'DOWN';
  }
};

export const findAssetIndexById = (assets: Asset[], id: number): number => {
  return assets.findIndex((a) => a.id === id);
};

export const findBoxIndexById = (boxes: BoundingBox[], id: number): number => {
  return boxes.findIndex((b) => b.id === id);
};

export const replaceBoxById = (boxes: BoundingBox[], newBox: BoundingBox) => {
  const index = boxes.findIndex((b) => b.id === newBox.id);
  if (index === -1) {
    // tslint:disable-next-line:no-console
    console.warn('box not found: ', newBox);
    return boxes;
  }
  const newBoxes = [...boxes];
  newBoxes[index] = newBox;
  return newBoxes;
}
