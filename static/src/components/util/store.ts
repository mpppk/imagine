export const saveBasePath = (wsName: string, basePath: string) => {
  localStorage.setItem(`${wsName}/basePath`, basePath);
}

export const loadBasePath = (wsName: string): string | null => {
  return localStorage.getItem(`${wsName}/basePath`);
}