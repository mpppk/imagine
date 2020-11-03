import {reducerWithInitialState} from 'typescript-fsa-reducers';
import {fsActionCreators} from "../actions/fs";

export const indexInitialState = {
  scanning: false,
  scanCount: 0,
};

export type IndexState = typeof indexInitialState;

const startScan = (state: IndexState) => {
  return {...state, scanning: true, scanCount: 0};
};

const finishOrCancelScan = (state: IndexState) => {
  return {...state, scanning: false};
};

export const indexPage = reducerWithInitialState(indexInitialState)
  .case(fsActionCreators.scanStart, (state) => {
    return startScan(state);
  })
  .case(fsActionCreators.scanRunning, (state, payload) => {
    return {...state, scanCount: state.scanCount+payload.foundedAssetsNum};
  })
  .case(fsActionCreators.scanCancel, (state) => {
    return finishOrCancelScan(state);
  })
  .case(fsActionCreators.scanFinish, (state) => {
    return finishOrCancelScan(state);
  });
