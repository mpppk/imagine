import { Middleware, MiddlewareAPI } from 'redux';
import { workspaceActionCreators } from './actions/workspace';
import { Action } from 'typescript-fsa';
import { assetActionCreators } from './actions/asset';
import { tagActionCreators } from './actions/tag';

const handle = (
  store: MiddlewareAPI,
  action: Action<any>
): Action<any> | void => {
  switch (action.type) {
    case workspaceActionCreators.scanRequest.type:
      const newAction = workspaceActionCreators.scanResult([
        { id: 1, name: 'default-workspace', basePath: '.' },
      ]);
      store.dispatch(newAction);
      break;
    case assetActionCreators.scanRequest.type:
      const a1 = assetActionCreators.scanRunning({
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
  }
};

export const makeMockLorcaMiddleware = (): Middleware => (store) => (next) => (
  action
) => {
  next(action);
  handle(store, action);
};
