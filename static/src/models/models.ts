export interface WorkSpace {
  name: string
}

export interface Tag {
  id: number
  name: string
}

export interface Asset {
  id: number
  name: string
  path: string
  tags: Tag[]
}