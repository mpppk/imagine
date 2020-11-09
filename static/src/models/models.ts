export interface WorkSpace {
  id: number
  name: string
  basePath: string
}

export interface Tag {
  id: number
  name: string
}

export interface BoundingBox {
  id: number
  tag: Tag
  x: number
  y: number
  width: number
  height: number
}

export type BoundingBoxRequest = Omit<BoundingBox, 'id'>;

export const newEmptyBoundingBox = (tag: Tag): BoundingBoxRequest => ({
  tag,
  x: 0, y: 0, width: 0, height: 0
});

export interface Asset {
  id: number
  name: string
  path: string
  boundingBoxes: BoundingBox[] | null
}

// 何番目のAssetか
export interface AssetWithIndex extends Asset {
  index: number
}

export type Direction = 'LEFT' | 'RIGHT' | 'UP' | 'DOWN';

export type Query = EqualsQuery | NotEqualsQuery | StartWithQuery | NoTagsQuery;

export type QueryOp = 'equals' | 'not-equals' | 'start-with' | 'no-tags';
export interface EqualsQuery {
  op: 'equals'
  value: string
}

export interface NotEqualsQuery {
  op: 'not-equals'
  value: string
}

export interface StartWithQuery {
  op: 'start-with'
  value: string
}

export interface NoTagsQuery {
  op: 'no-tags'
  value: string
}
