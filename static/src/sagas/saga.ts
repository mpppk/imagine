import {all, call, put, select, takeEvery} from '@redux-saga/core/effects';
import {ActionCreator} from "typescript-fsa";
import { SagaIterator } from 'redux-saga';
import {workspaceActionCreators} from '../actions/workspace';
import {BoundingBoxRequest, WorkSpace} from "../models/models";
import {indexActionCreators} from "../actions";
import {boundingBoxActionCreators, BoundingBoxAssignRequestPayload} from "../actions/box";
import {State} from "../reducers/reducer";

const scanWorkSpacesWorker = function*(workspaces: WorkSpace[]) {
  return yield put(workspaceActionCreators.select(workspaces[0]));
}

const downNumberKeyWorker = function*(key: number) {
  const state: State = yield select();
  if (key > state.global.tags.length || !state.global.selectedAsset) {
    return;
  }
  const tag = state.global.tags[key];
  const box: BoundingBoxRequest = { // FIXME
    tag,
    x: 0,
    y: 0,
    width: 0,
    height: 0,
  };

  const payload: BoundingBoxAssignRequestPayload = {
    asset: state.global.selectedAsset,
    box,
  };

  return yield put(boundingBoxActionCreators.assignRequest(payload));
}

export default function* rootSaga() {
  yield all([
    takeEveryAction(workspaceActionCreators.scanResult, scanWorkSpacesWorker)(),
    takeEveryAction(indexActionCreators.downNumberKey, downNumberKeyWorker)()
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

