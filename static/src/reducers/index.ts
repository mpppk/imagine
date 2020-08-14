import {reducerWithInitialState} from 'typescript-fsa-reducers';
import {serverActionCreators} from "../actions/server";

export const indexInitialState = {
  imagePaths: [] as string[],
  scanning: false,
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
  });
