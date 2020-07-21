export const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

export const immutableSplice = <T>(array: T[], start: number, deleteCount: number, ...item: T[]): T[] => {
  return [...array.slice(0, start), ...item, ...array.slice(start+deleteCount)]
}