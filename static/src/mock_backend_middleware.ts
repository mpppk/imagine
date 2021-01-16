import { Middleware, MiddlewareAPI } from 'redux';
import { workspaceActionCreators } from './actions/workspace';
import { Action } from 'typescript-fsa';
import { assetActionCreators } from './actions/asset';
import { tagActionCreators } from './actions/tag';
import {
  boundingBoxActionCreators,
  BoundingBoxAssignRequestPayload,
  BoundingBoxUnAssignRequestPayload,
} from './actions/box';
import { Asset } from './models/models';
import { fsActionCreators } from './actions/fs';
import {
  ClickChangeBaseButtonPathPayload,
  indexActionCreators,
} from './actions';

const handle = (
  store: MiddlewareAPI,
  action: Action<any>
): Action<any> | void => {
  const boundingBoxHandler = new BoundingBoxActionHandler(store);
  const indexActionHandler = new IndexActionHandler(store);
  const workSpaceActionHandler = new WorkSpaceActionHandler(store);
  switch (action.type) {
    case indexActionCreators.clickChangeBasePathButton.type:
      indexActionHandler.clickChangeBasePathButton(action);
      break;
    case workspaceActionCreators.scanRequest.type:
      workSpaceActionHandler.scanRequest();
      break;
    case assetActionCreators.scanRequest.type:
      const a1 = assetActionCreators.scanRunning({
        count: 3,
        assets: [
          {
            id: 1,
            name: 'path1',
            path: 'path1',
            boundingBoxes: [
              { id: 1, tagID: 1, x: 0, y: 0, width: 0, height: 0 },
            ],
          },
          {
            id: 2,
            name: 'path2',
            path: 'path2',
            boundingBoxes: [
              { id: 1, tagID: 2, x: 0, y: 0, width: 0, height: 0 },
            ],
          },
          {
            id: 3,
            name: 'path3',
            path: 'path3',
            boundingBoxes: [
              { id: 1, tagID: 3, x: 0, y: 0, width: 0, height: 0 },
            ],
          },
        ],
        workSpaceName: 'default-workspace',
      });
      store.dispatch(a1);
      const a2 = assetActionCreators.scanFinish({
        workSpaceName: 'default-workspace',
        count: 3,
      });
      store.dispatch(a2);
      break;
    case workspaceActionCreators.select.type:
      const a = tagActionCreators.scanResult({
        tags: [
          { id: 1, name: 'tag1' },
          { id: 2, name: 'tag2' },
          { id: 3, name: 'tag3' },
        ],
        workSpaceName: 'default-workspace',
      });
      store.dispatch(a);
      break;
    case boundingBoxActionCreators.assignRequest.type:
      boundingBoxHandler.assignRequest(action);
      break;
    case boundingBoxActionCreators.unAssignRequest.type:
      boundingBoxHandler.unassignRequest(action);
      break;
  }
};

class IndexActionHandler {
  constructor(private store: MiddlewareAPI) {}

  clickChangeBasePathButton(_action: Action<ClickChangeBaseButtonPathPayload>) {
    const a = fsActionCreators.baseDirSelect({
      basePath: 'new-base-path',
      workSpaceName: 'default-workspace',
    });
    this.store.dispatch(a);
  }
}

// tslint:disable-next-line:max-classes-per-file
class WorkSpaceActionHandler {
  constructor(private store: MiddlewareAPI) {}

  scanRequest() {
    const newAction = workspaceActionCreators.scanResult({
      basePath: 'default/base/path',
      workspaces: [{ id: 1, name: 'default-workspace', basePath: '.' }],
    });
    this.store.dispatch(newAction);
  }
}

// tslint:disable-next-line:max-classes-per-file
class BoundingBoxActionHandler {
  constructor(private store: MiddlewareAPI) {}

  assignRequest(action: Action<BoundingBoxAssignRequestPayload>) {
    const boundingBoxes = action.payload.asset.boundingBoxes ?? [];
    boundingBoxes.push({ ...action.payload.box, id: boundingBoxes.length + 1 });
    const asset: Asset = { ...action.payload.asset, boundingBoxes };
    const a = boundingBoxActionCreators.assign({
      asset,
      box: action.payload.box,
      workSpaceName: 'default-workspace',
    });
    this.store.dispatch(a);
  }

  unassignRequest(action: Action<BoundingBoxUnAssignRequestPayload>) {
    let boundingBoxes = action.payload.asset.boundingBoxes ?? [];
    boundingBoxes = boundingBoxes.filter((b) => b.id !== action.payload.boxID);
    const asset: Asset = { ...action.payload.asset, boundingBoxes };
    const a = boundingBoxActionCreators.unAssign({
      asset,
      boxID: action.payload.boxID,
      workSpaceName: 'default-workspace',
    });
    this.store.dispatch(a);
  }
}

export const makeMockLorcaMiddleware = (): Middleware => (store) => (next) => (
  action
) => {
  next(action);
  handle(store, action);
};
