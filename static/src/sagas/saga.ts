import {all, call, fork, put, select, take, takeEvery, takeLatest} from '@redux-saga/core/effects';
import {ActionCreator} from 'typescript-fsa';
import {eventChannel, SagaIterator} from 'redux-saga';
import {workspaceActionCreators} from '../actions/workspace';
import {AssetWithIndex, newEmptyBoundingBox, Query, Tag, WorkSpace} from '../models/models';
import {ClickFilterApplyButtonPayload, indexActionCreators} from '../actions';
import {
  boundingBoxActionCreators,
  BoundingBoxMovePayload,
  BoundingBoxScalePayload,
  BoundingBoxUnAssignRequestPayload,
} from '../actions/box';
import {State} from '../reducers/reducer';
import {findAssetIndexById, findBoxIndexById, isDefaultBox} from '../util';
import {browserActionCreators} from "../actions/browser";
import debounce from "lodash/debounce";
import {assetActionCreators} from "../actions/asset";
import {toQuery} from "../reducers/global";
import {fsActionCreators} from "../actions/fs";

const scanWorkSpacesWorker = function* (workspaces: WorkSpace[]) {
  return yield put(workspaceActionCreators.select(workspaces[0]));
};

const fsScanStartWorkSpacesWorker = function* () {
  const p = yield select((s: State) => ({
    queries: s.global.queries,
    workSpaceName: s.global.currentWorkSpace?.name,
  }));

  if (p.workSpaceName === undefined) {
    return;
  }

  return yield put(assetActionCreators.scanRequest({
    ...p,
    requestNum: 10, // FIXME
    reset: true,
  }));
};

const boxMoveWorker = function* (payload: BoundingBoxMovePayload) {
  const state = yield select((s: State) => s.global);
  // FIXME: O(n)
  const index = findAssetIndexById(state.assets, payload.assetID);
  const asset = state.assets[index];
  if (asset.boundingBoxes == null) {
    return state;
  }
  const boxIndex = findBoxIndexById(asset.boundingBoxes, payload.boxID);
  const box = asset.boundingBoxes[boxIndex];
  const newBox = {
    ...box,
    x: payload.dx,
    y: payload.dy,
  }
  return yield put(boundingBoxActionCreators.modifyRequest({
    workSpaceName: state.currentWorkSpace.name,
    asset,
    box: newBox,
  }));
};

const boxScaleWorker = function* (payload: BoundingBoxScalePayload) {
  const state = yield select((s: State) => s.global);
  // FIXME: O(n)
  const index = findAssetIndexById(state.assets, payload.assetID);
  const asset = state.assets[index];
  if (asset.boundingBoxes == null) {
    return state;
  }
  const boxIndex = findBoxIndexById(asset.boundingBoxes, payload.boxID);
  const box = asset.boundingBoxes[boxIndex];
  const newBox = {
    ...box,
    width: payload.dx,
    height: payload.dy,
  }
  return yield put(boundingBoxActionCreators.modifyRequest({
    workSpaceName: state.currentWorkSpace.name,
    asset,
    box: newBox,
  }));
};

const clickFilterApplyButtonWorker = function* (payload: ClickFilterApplyButtonPayload) {
  const state = yield select((s: State) => s.global);
  if (!payload.changed && !payload.enabled) {
    return;
  }

  let queries = [] as Query[];
  if (payload.enabled) {
    queries = payload.queryInputs
      .map(toQuery.bind(null, state.tags))
      .filter((q): q is Query => q !== null);
  }
  return yield put(assetActionCreators.scanRequest({
    queries,
    requestNum: 10, // FIXME
    workSpaceName: state.currentWorkSpace.name,
    reset: true,
  }));
};

const selectTagWorker = function* (tag: Tag): any { // FIXME any
  const [asset, workSpaceName]: [AssetWithIndex | null, string?] = yield select<(s: State) => [AssetWithIndex | null, string?]>((s) => [s.global.selectedAsset, s.global.currentWorkSpace?.name]);
  if (workSpaceName === undefined) {
    // tslint:disable-next-line:no-console
    console.warn('failed to request to assign tag because workspace name is undefined');
    return;
  }

  if (asset === null) {
    // tslint:disable-next-line:no-console
    console.warn('failed to request to assign tag because asset is not selected');
    return;
  }

  let boxes = asset.boundingBoxes;
  if (boxes === null) {
    return yield put(boundingBoxActionCreators.assignRequest({
      asset, box: newEmptyBoundingBox(tag), workSpaceName,
    }));
  }
  boxes = boxes.filter(isDefaultBox).filter((box) => box.tag.id === tag.id);
  if (boxes.length === 0) {
    return yield put(boundingBoxActionCreators.assignRequest({
      asset, box: newEmptyBoundingBox(tag), workSpaceName,
    }));
  }

  for (const box of boxes) {
    const payload: BoundingBoxUnAssignRequestPayload = {
      asset,
      boxID: box.id,
      workSpaceName,
    };
    yield put(boundingBoxActionCreators.unAssignRequest(payload));
  }
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
  yield put(indexActionCreators.selectTag(tag));
};

export default function* rootSaga() {
  yield fork(resizeSaga);
  yield all([
    takeEveryAction(workspaceActionCreators.scanResult, scanWorkSpacesWorker)(),
    takeEveryAction(fsActionCreators.scanStart, fsScanStartWorkSpacesWorker)(),
    takeEveryAction(indexActionCreators.downNumberKey, downNumberKeyWorker)(),
    takeEveryAction(indexActionCreators.selectTag, selectTagWorker)(),
    takeEveryAction(indexActionCreators.clickFilterApplyButton, clickFilterApplyButtonWorker)(),
    takeEveryAction(boundingBoxActionCreators.move, boxMoveWorker)(),
    takeEveryAction(boundingBoxActionCreators.scale, boxScaleWorker)(),
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

export const takeLatestAction = <T>(
  ac: ActionCreator<T>,
  worker: (params: T, ...args: any[]) => SagaIterator
) => {
  return function* () {
    yield takeLatest(ac, function* (action) {
      yield call(worker, action.payload);
    });
  };
};

function resize() {
  return eventChannel(emitter => {
      if (process.browser) {
        const resizeEventHandler = debounce(() => {
          const width = window.innerWidth;
          const height = window.innerHeight;
          emitter({width, height});
        }, 200);
        window.addEventListener('resize', resizeEventHandler);
      }
      // tslint:disable-next-line:no-empty
      return () => {
      };
    }
  )
}

export function* resizeSaga() {
  const chan = yield call(resize);
  try {
    while (true) {
      const payload = yield take(chan)
      yield put(browserActionCreators.resize(payload));
    }
    // tslint:disable-next-line:no-empty
  } finally {
  }
}