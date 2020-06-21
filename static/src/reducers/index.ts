import { reducerWithInitialState } from 'typescript-fsa-reducers';
import {serverActionCreators} from "../actions/server";

export const indexInitialState = {
  scanning: false
};

export type IndexState = typeof indexInitialState;

const startScan = (state: IndexState) => {
  return { ...state, scanning: true};
};

// const finishScan = (state: IndexState) => {
//   return { ...state, scanning: false};
// };

export const indexPage = reducerWithInitialState(indexInitialState)
  .case(serverActionCreators.startDirectoryScanning, (state) => {
    return  startScan(state);
  });
