import {Tag} from "./models/models";

export const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

export const immutableSplice = <T>(array: T[], start: number, deleteCount: number, ...item: T[]): T[] => {
  return [...array.slice(0, start), ...item, ...array.slice(start+deleteCount)]
}

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
}

