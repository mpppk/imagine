export interface WorkSpace {
  name: string
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

export interface Asset {
  id: number
  name: string
  path: string
  boundingBoxes: BoundingBox[]
}

// 何番目のAssetか
export interface AssetWithIndex extends Asset {
  index: number
}

export type Direction = 'LEFT' | 'RIGHT' | 'UP' | 'DOWN';