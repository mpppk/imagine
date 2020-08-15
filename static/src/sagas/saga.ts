import {all, call, put, takeEvery} from '@redux-saga/core/effects';
import {ActionCreator} from "typescript-fsa";
import { SagaIterator } from 'redux-saga';
import {serverActionCreators} from "../actions/server";
import {globalActionCreators} from "../actions/global";
import {WorkSpace} from "../models/models";

const scanWorkSpacesWorker = function*(workspaces: WorkSpace[]) {
  return yield put(globalActionCreators.selectNewWorkSpace(workspaces[0]));
}

export default function* rootSaga() {
  yield all([
    takeEveryAction(serverActionCreators.scanWorkSpaces, scanWorkSpacesWorker)()
  ]);
}

// export type SagaWorker = <T>(params: T, ...args: any[]) => SagaIterator;
export const takeEveryAction = <T>(
  ac: ActionCreator<T>,
  worker: (params: T, ...args: any[]) => SagaIterator
) => {
  return function* () {
    yield takeEvery(ac, function* (action) {
      yield call(worker, action.payload);
    });
  };
};

