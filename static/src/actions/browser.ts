import actionCreatorFactory from 'typescript-fsa';

const browserActionCreatorFactory = actionCreatorFactory('BROWSER');

export interface ResizePayload {
  width: number;
  height: number;
}

export const browserActionCreators = {
  resize: browserActionCreatorFactory<ResizePayload>('RESIZE'),
};
