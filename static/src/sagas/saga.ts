import {
  all,
  call,
  fork,
  put,
  select,
  take,
  takeEvery,
  takeLatest,
} from '@redux-saga/core/effects';
import { ActionCreator } from 'typescript-fsa';
import { eventChannel, SagaIterator } from 'redux-saga';
import {
  workspaceActionCreators,
  WorkSpaceScanResultPayload,
} from '../actions/workspace';
import { Asset, newEmptyBoundingBox, Tag, WorkSpace } from '../models/models';
import { ClickFilterApplyButtonPayload, indexActionCreators } from '../actions';
import {
  boundingBoxActionCreators,
  BoundingBoxMovePayload,
  BoundingBoxScalePayload,
  BoundingBoxUnAssignRequestPayload,
} from '../actions/box';
import { State } from '../reducers/reducer';
import { findAssetIndexById, findBoxIndexById, isDefaultBox } from '../util';
import { browserActionCreators } from '../actions/browser';
import debounce from 'lodash/debounce';
import { assetActionCreators, AssetScanResultPayload } from '../actions/asset';
import { loadBasePath, saveBasePath } from '../components/util/store';
import { BaseDirSelectPayload, fsActionCreators } from '../actions/fs';

const selectWorkSpaceWorker = function* (workspace: WorkSpace) {
  const basePath = loadBasePath(workspace.name) ?? workspace.basePath;
  const action = fsActionCreators.baseDirSelect({
    basePath,
    workSpaceName: workspace.name,
  });
  return yield put(action);
};

const scanWorkSpacesWorker = function* (payload: WorkSpaceScanResultPayload) {
  // FIXME
  const workspace = payload.workspaces[0];
  const hasBasePath = ![undefined, null, ''].includes(payload.basePath);
  if (hasBasePath) {
    saveBasePath(workspace.name, payload.basePath);
  }

  return yield put(
    workspaceActionCreators.select({
      ...workspace,
      basePath: hasBasePath ? payload.basePath : workspace.basePath,
    })
  );
};

const baseDirSelectWorker = function* (payload: BaseDirSelectPayload) {
  saveBasePath(payload.workSpaceName, payload.basePath);
  const needToLoadAssets = yield select(
    (state: State) => state.global.needToLoadAssets
  );
  if (needToLoadAssets) {
    yield put(
      fsActionCreators.scanRequest({
        workSpaceName: payload.workSpaceName,
        basePath: payload.basePath,
      })
    );
  }
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
  };
  return yield put(
    boundingBoxActionCreators.modifyRequest({
      workSpaceName: state.currentWorkSpace.name,
      asset,
      box: newBox,
    })
  );
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
  };
  return yield put(
    boundingBoxActionCreators.modifyRequest({
      workSpaceName: state.currentWorkSpace.name,
      asset,
      box: newBox,
    })
  );
};

const clickFilterApplyButtonWorker = function* (
  payload: ClickFilterApplyButtonPayload
) {
  const state = yield select((s: State) => s.global);
  if (!payload.changed && !payload.enabled) {
    return;
  }

  return yield put(
    assetActionCreators.scanRequest({
      queries: payload.queries,
      requestNum: 50, // FIXME
      workSpaceName: state.currentWorkSpace.name,
      reset: true,
    })
  );
};

const fsScanRunningWorker = function* () {
  const state = yield select((s: State) => s.global);
  if (state.assets.length !== 0) {
    return;
  }
  return yield put(
    assetActionCreators.scanRequest({
      queries: state.queries,
      requestNum: 50, // FIXME
      workSpaceName: state.currentWorkSpace.name,
      reset: true,
    })
  );
};

const assetScanResultWorker = function* (payload: AssetScanResultPayload) {
  const hasQueries = yield select((s: State) => s.global.queries.length > 0);
  const basePath = yield select(
    (s: State) => s.global.currentWorkSpace?.basePath
  );
  if (payload.count === 0 && !hasQueries && basePath !== undefined) {
    yield put(
      fsActionCreators.scanRequest({
        workSpaceName: payload.workSpaceName,
        basePath,
      })
    );
  }
};

const selectTagWorker = function* (tag: Tag): any {
  // FIXME any
  const [asset, workSpaceName]: [Asset | null, string?] = yield select<
    (s: State) => [Asset | null, string?]
  >((s) => [s.global.selectedAsset, s.global.currentWorkSpace?.name]);
  if (workSpaceName === undefined) {
    // tslint:disable-next-line:no-console
    console.warn(
      'failed to request to assign tag because workspace name is undefined'
    );
    return;
  }

  if (asset === null) {
    // tslint:disable-next-line:no-console
    console.warn(
      'failed to request to assign tag because asset is not selected'
    );
    return;
  }

  let boxes = asset.boundingBoxes;
  if (boxes === null) {
    return yield put(
      boundingBoxActionCreators.assignRequest({
        asset,
        box: newEmptyBoundingBox(tag.id),
        workSpaceName,
      })
    );
  }
  boxes = boxes.filter(isDefaultBox).filter((box) => box.tagID === tag.id);
  if (boxes.length === 0) {
    return yield put(
      boundingBoxActionCreators.assignRequest({
        asset,
        box: newEmptyBoundingBox(tag.id),
        workSpaceName,
      })
    );
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

const downAlphabetKeyWorker = function* (key: string) {
  const state: State = yield select();
  if (!state.global.selectedAsset) {
    return;
  }

  if (!state.global.currentWorkSpace) {
    // tslint:disable-next-line:no-console
    console.info('workspace is not selected');
    return;
  }

  let index: number | undefined;
  switch (key) {
    case 'Q':
      index = 4;
      break;
    case 'W':
      index = 5;
      break;
    case 'E':
      index = 6;
      break;
    case 'R':
      index = 7;
      break;
    case 'A':
      index = 8;
      break;
    case 'S':
      index = 9;
      break;
    case 'D':
      index = 10;
      break;
    case 'F':
      index = 11;
      break;
    case 'Z':
      index = 12;
      break;
    case 'X':
      index = 13;
      break;
    case 'C':
      index = 14;
      break;
    case 'V':
      index = 15;
      break;
  }

  // tag list is 0-indexed, but number key is 1-indexed
  if (index !== undefined && index < state.global.tags.length) {
    const tag = state.global.tags[index];
    yield put(indexActionCreators.selectTag(tag));
  }
};

const downNumberKeyWorker = function* (key: number) {
  const state: State = yield select();
  if (key > state.global.tags.length || !state.global.selectedAsset) {
    return;
  }

  if (!state.global.currentWorkSpace) {
    // tslint:disable-next-line:no-console
    console.info('workspace is not selected');
    return;
  }

  // tag list is 0-indexed, but number key is 1-indexed
  const index = key === 0 ? 9 : key - 1;
  const tag = state.global.tags[index];
  yield put(indexActionCreators.selectTag(tag));
};

const downSymbolKeyWorker = function* (code: number) {
  const state: State = yield select();
  if (!state.global.selectedAsset) {
    return;
  }

  if (!state.global.currentWorkSpace) {
    // tslint:disable-next-line:no-console
    console.info(
      'bounding box assign/unassign request is not sent because workspace is not selected'
    );
    return;
  }

  switch (code) {
    case 189: // -
      if (state.global.tags.length > 10) {
        yield put(indexActionCreators.selectTag(state.global.tags[10]));
      }
      return;
    case 187: // ^
      if (state.global.tags.length > 11) {
        yield put(indexActionCreators.selectTag(state.global.tags[11]));
      }
      return;
    case 0: // ¥
      if (state.global.tags.length > 12) {
        yield put(indexActionCreators.selectTag(state.global.tags[12]));
      }
      return;
  }
};

export default function* rootSaga() {
  yield fork(resizeSaga);
  yield all([
    takeEveryAction(workspaceActionCreators.scanResult, scanWorkSpacesWorker)(),
    takeEveryAction(workspaceActionCreators.select, selectWorkSpaceWorker)(),
    takeEveryAction(fsActionCreators.baseDirSelect, baseDirSelectWorker)(),
    takeEveryAction(fsActionCreators.scanRunning, fsScanRunningWorker)(),
    takeEveryAction(assetActionCreators.scanFinish, assetScanResultWorker)(),
    takeEveryAction(
      indexActionCreators.downAlphabetKey,
      downAlphabetKeyWorker
    )(),
    takeEveryAction(indexActionCreators.downNumberKey, downNumberKeyWorker)(),
    takeEveryAction(indexActionCreators.downSymbolKey, downSymbolKeyWorker)(),
    takeEveryAction(indexActionCreators.selectTag, selectTagWorker)(),
    takeEveryAction(
      indexActionCreators.clickFilterApplyButton,
      clickFilterApplyButtonWorker
    )(),
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
  return eventChannel((emitter) => {
    if (process.browser) {
      const resizeEventHandler = debounce(() => {
        const width = window.innerWidth;
        const height = window.innerHeight;
        emitter({ width, height });
      }, 200);
      window.addEventListener('resize', resizeEventHandler);
    }
    // tslint:disable-next-line:no-empty
    return () => {};
  });
}

export function* resizeSaga() {
  const chan = yield call(resize);
  try {
    while (true) {
      const payload = yield take(chan);
      yield put(browserActionCreators.resize(payload));
    }
    // tslint:disable-next-line:no-empty
  } finally {
  }
}
