import { reducerWithInitialState } from 'typescript-fsa-reducers';
import { assetActionCreators } from '../actions/asset';
import { workspaceActionCreators } from '../actions/workspace';
import { Asset, BoundingBox, Query, Tag, WorkSpace } from '../models/models';
import { indexActionCreators } from '../actions';
import {
  findAssetIndexById,
  immutableSplice,
  replaceBoxById,
  replaceBy,
} from '../util';
import { tagActionCreators } from '../actions/tag';
import { boundingBoxActionCreators } from '../actions/box';
import { browserActionCreators } from '../actions/browser';
import { fsActionCreators } from '../actions/fs';

export const globalInitialState = {
  filterEnabled: false,
  queries: [] as Query[],
  assets: [] as Asset[],
  selectedAsset: null as Asset | null,
  tags: [] as Tag[],
  currentWorkSpace: null as WorkSpace | null,
  initialBoundingBox: null as BoundingBox | null,
  hasMoreAssets: true,
  isLoadingWorkSpaces: true,
  isScanningAssets: false,
  selectedTagId: undefined as number | undefined,
  workspaces: null as WorkSpace[] | null,
  windowHeight: 720,
};
export type GlobalState = typeof globalInitialState;
export const global = reducerWithInitialState(globalInitialState)
  .case(browserActionCreators.resize, (state, payload) => {
    return { ...state, windowHeight: payload.height };
  })
  .case(fsActionCreators.baseDirSelect, (state, payload) => {
    if (state.currentWorkSpace === null || state.workspaces === null) {
      return { ...state };
    }
    const currentWorkSpace: WorkSpace = {
      ...state.currentWorkSpace,
      basePath: payload.basePath,
    };
    const workspaces = replaceBy(
      state.workspaces,
      currentWorkSpace,
      (w) => w.name === currentWorkSpace.name
    );
    return { ...state, assets: [], currentWorkSpace, workspaces };
  })
  .case(workspaceActionCreators.scanResult, (state, workspaces) => {
    return { ...state, workspaces, isLoadingWorkSpaces: false };
  })
  .case(assetActionCreators.scanRequest, (state) => {
    return { ...state, isScanningAssets: true, hasMoreAssets: true };
  })
  .case(assetActionCreators.scanRunning, (state, payload) => {
    const assets = [...state.assets, ...payload.assets];
    return { ...state, isScanningAssets: false, assets };
  })
  .case(assetActionCreators.scanFinish, (state) => {
    return { ...state, isScanningAssets: false, hasMoreAssets: false };
  })
  .case(workspaceActionCreators.select, (state, workspace) => {
    return { ...state, currentWorkSpace: { ...workspace } };
  })
  .case(fsActionCreators.baseDirSelect, (state, payload) => {
    if (state.currentWorkSpace === null) {
      return state;
    }
    return {
      ...state,
      currentWorkSpace: {
        ...state.currentWorkSpace,
        basePath: payload.basePath,
      },
    };
  })
  .case(workspaceActionCreators.updateRequest, (state, workspace) => {
    if (state.workspaces === null) {
      return { ...state, workspaces: [workspace] };
    }
    const workspaces = replaceBy(
      state.workspaces,
      workspace,
      (w) => w.name === workspace.name
    );
    const currentWorkSpace =
      state.currentWorkSpace === null ||
      state.currentWorkSpace.name !== workspace.name
        ? state.currentWorkSpace
        : workspace;
    return { ...state, currentWorkSpace, workspaces };
  })
  .case(indexActionCreators.selectTag, (state, tag) => {
    return { ...state, selectedTagId: tag.id };
  })
  .case(indexActionCreators.clickAddTagButton, (state, tag) => {
    return { ...state, tags: [tag, ...state.tags] };
  })
  .case(indexActionCreators.clickFilterApplyButton, (state, payload) => {
    if (!state.filterEnabled && !payload.enabled) {
      return { ...state };
    }
    return {
      ...resetAssets(state),
      queries: payload.queries,
      filterEnabled: payload.enabled,
      hasMoreAssets: true,
    };
  })
  .case(tagActionCreators.scanResult, (state, payload) => {
    return { ...state, tags: payload.tags };
  })
  .case(tagActionCreators.rename, (state, payload) => {
    const targetTagIndex = state.tags.findIndex((t) => t.id === payload.tag.id);
    if (targetTagIndex === -1) {
      // tslint:disable-next-line:no-console
      console.warn('unknown tag ID is provided', payload.tag);
    }
    return {
      ...state,
      tags: immutableSplice(state.tags, targetTagIndex, 1, payload.tag),
    };
  })
  .case(tagActionCreators.update, (state, payload) => {
    return { ...state, tags: payload.tags };
  })
  .case(indexActionCreators.assetSelect, (state, asset) => {
    return { ...state, selectedAsset: { ...asset } };
  })
  .case(boundingBoxActionCreators.assign, (state, payload) => {
    return { ...state, ...updateAssets(state, { ...payload.asset }) };
  })
  .case(boundingBoxActionCreators.unAssign, (state, payload) => {
    return { ...state, ...updateAssets(state, { ...payload.asset }) };
  })
  .case(boundingBoxActionCreators.modifyRequest, (state, payload) => {
    if (payload.asset.boundingBoxes == null) {
      return state;
    }
    const newBoxes = replaceBoxById(payload.asset.boundingBoxes, payload.box);
    return {
      ...state,
      ...updateAssets(state, { ...payload.asset, boundingBoxes: newBoxes }),
    };
  })
  .case(boundingBoxActionCreators.deleteRequest, (state, payload) => {
    if (
      state.selectedAsset === null ||
      state.selectedAsset.boundingBoxes == null
    ) {
      return state;
    }
    const boundingBoxes = state.selectedAsset.boundingBoxes.filter(
      (b) => b.id !== payload.boxID
    );
    return {
      ...state,
      ...updateAssets(state, { ...state.selectedAsset, boundingBoxes }),
    };
  })
  .case(indexActionCreators.downArrowKey, (state, payload) => {
    if (!state.selectedAsset) {
      return { ...state };
    }
    const index = findAssetIndexById(state.assets, state.selectedAsset.id);
    switch (payload) {
      case 'UP':
        if (index === 0) {
          return { ...state };
        }
        return {
          ...state,
          selectedAsset: { ...state.assets[index - 1], index: index - 1 },
        };
      case 'DOWN':
        if (index === state.assets.length - 1) {
          return { ...state };
        }
        return {
          ...state,
          selectedAsset: { ...state.assets[index + 1], index: index + 1 },
        };
      default:
        return { ...state };
    }
  });

const resetAssets = (state: GlobalState): GlobalState => {
  return {
    ...state,
    assets: [],
    selectedAsset: null,
    selectedTagId: undefined,
  };
};

const updateAssets = (state: GlobalState, asset: Asset) => {
  const assets = replaceBy(state.assets, asset, (a) => a.id === asset.id);
  const selectedAsset =
    state.selectedAsset?.id === asset.id ? asset : state.selectedAsset;
  return { assets, selectedAsset };
};
