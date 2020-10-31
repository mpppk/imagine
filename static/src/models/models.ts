export interface WorkSpace {
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

export type Query = EqualsQuery | NotEqualsQuery;

export type QueryOp = 'equals' | 'not-equals';
export interface EqualsQuery {
  op: 'equals'
  tag: Tag
}

export interface NotEqualsQuery {
  op: 'not-equals'
  tag: Tag
}

export interface QueryInput {
  op: QueryOp
  tagName: string
}

