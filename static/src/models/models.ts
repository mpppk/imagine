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

export interface Asset {
  id: number
  name: string
  path: string
  boundingBoxes: BoundingBox[]
}