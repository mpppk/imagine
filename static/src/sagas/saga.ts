import { all, call, put, select, takeEvery } from '@redux-saga/core/effects';
import { ActionCreator } from 'typescript-fsa';
import { SagaIterator } from 'redux-saga';
import { workspaceActionCreators } from '../actions/workspace';
import { BoundingBoxRequest, WorkSpace } from '../models/models';
import { indexActionCreators } from '../actions';
import {
  boundingBoxActionCreators,
  BoundingBoxAssignRequestPayload,
  BoundingBoxUnAssignRequestPayload,
} from '../actions/box';
import { State } from '../reducers/reducer';
import { isDefaultBox } from '../util';

const scanWorkSpacesWorker = function* (workspaces: WorkSpace[]) {
  return yield put(workspaceActionCreators.select(workspaces[0]));
};

const downNumberKeyWorker = function* (key: number) {
  const state: State = yield select();
  if (key > state.global.tags.length || !state.global.selectedAsset) {
    return;
  }

  if (!state.global.currentWorkSpace) {
    // tslint:disable-next-line:no-console
    console.info(
      'bounding box assign/unassign request is not sent because workspace is not selected'
    );
    return;
  }

  // tag list is 0-indexed, but number key is 1-indexed
  const tag = state.global.tags[key - 1];

  // 初期状態のboxが存在する場合はunassign
  let boxes = state.global.selectedAsset.boundingBoxes ?? [];
  boxes = boxes.filter(isDefaultBox).filter((box) => box.tag.id === tag.id);
  if (boxes.length > 0) {
    // unassign
    for (const box of boxes) {
      const payload: BoundingBoxUnAssignRequestPayload = {
        asset: state.global.selectedAsset,
        boxID: box.id,
        workSpaceName: state.global.currentWorkSpace!.name,
      };
      yield put(boundingBoxActionCreators.unAssignRequest(payload));
    }
    return;
  } else {
    // assign
    const box: BoundingBoxRequest = {
      // FIXME
      tag,
      x: 0,
      y: 0,
      width: 0,
      height: 0,
    };

    const payload: BoundingBoxAssignRequestPayload = {
      asset: state.global.selectedAsset,
      box,
      workSpaceName: state.global.currentWorkSpace!.name,
    };

    return yield put(boundingBoxActionCreators.assignRequest(payload));
  }
};

export default function* rootSaga() {
  yield all([
    takeEveryAction(workspaceActionCreators.scanResult, scanWorkSpacesWorker)(),
    takeEveryAction(indexActionCreators.downNumberKey, downNumberKeyWorker)(),
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
