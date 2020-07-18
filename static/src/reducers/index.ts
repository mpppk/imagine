import {reducerWithInitialState} from 'typescript-fsa-reducers';
import {serverActionCreators} from "../actions/server";

export const indexInitialState = {
  scanning: false,
  imagePaths: [] as string[],
};

export type IndexState = typeof indexInitialState;

const startScan = (state: IndexState) => {
  return {...state, scanning: true};
};

const finishOrCancelScan = (state: IndexState) => {
  return {...state, scanning: false};
};

export const indexPage = reducerWithInitialState(indexInitialState)
  .case(serverActionCreators.startDirectoryScanning, (state) => {
    return startScan(state);
  })
  .case(serverActionCreators.cancelDirectoryScanning, (state) => {
    return finishOrCancelScan(state);
  })
  .case(serverActionCreators.finishDirectoryScanning, (state) => {
    return finishOrCancelScan(state);
  })
  .case(serverActionCreators.scanningImages, (state, payload) => {
    const newPaths = payload.map(p => `http://localhost:1323/static${p}`)
    return {...state, imagePaths: [...state.imagePaths, ...newPaths]};
  });
