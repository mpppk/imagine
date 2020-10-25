import {reducerWithInitialState} from 'typescript-fsa-reducers';
import {assetActionCreators} from "../actions/asset";
import {workspaceActionCreators} from '../actions/workspace';
import {Asset, AssetWithIndex, BoundingBox, Tag, WorkSpace} from '../models/models';
import {indexActionCreators} from "../actions";
import {findAssetIndexById, findBoxIndexById, immutableSplice, replaceBoxById, replaceBy} from "../util";
import {tagActionCreators} from "../actions/tag";
import {boundingBoxActionCreators} from "../actions/box";
import {browserActionCreators} from "../actions/browser";

export const globalInitialState = {
  assets: [] as Asset[],
  selectedAsset: null as AssetWithIndex | null,
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
    return {...state, windowHeight: payload.height};
  })
  .case(workspaceActionCreators.scanResult, (state, workspaces) => {
    return {...state, workspaces, isLoadingWorkSpaces: false};
  })
  .case(assetActionCreators.scanRequest, (state) => {
    return {...state, isScanningAssets: true}
  })
  .case(assetActionCreators.scanRunning, (state, payload) => {
    const assets = {...state.assets};
    let assetsNum = Object.keys(assets).length;
    payload.assets.reduce((prev, cur) => {
      assetsNum++;
      return {...prev, [assetsNum]: cur}
    }, {...state.assets});
    return {...state, isScanningAssets: false, assets: [...state.assets, ...payload.assets]}
  })
  .case(assetActionCreators.scanFinish, (state) => {
    return {...state, isScanningAssets: false, hasMoreAssets: false}
  })
  .case(workspaceActionCreators.select, (state, workspace) => {
    return {...state, currentWorkSpace: workspace};
  })
  .case(indexActionCreators.selectTag, (state, tag) => {
    return {...state, selectedTagId: tag.id};
  })
  .case(indexActionCreators.clickAddTagButton, (state, tag) => {
    return {...state, tags: [tag, ...state.tags]};
  })
  .case(tagActionCreators.scanResult, (state, payload) => {
    return {...state, tags: payload.tags};
  })
  .case(tagActionCreators.rename, (state, payload) => {
    const targetTagIndex = state.tags.findIndex((t) => t.id === payload.tag.id);
    if (targetTagIndex === -1) {
      // tslint:disable-next-line:no-console
      console.warn('unknown tag ID is provided', payload.tag);
    }
    return {...state, tags: immutableSplice(state.tags, targetTagIndex, 1, payload.tag)};
  })
  .case(tagActionCreators.update, (state, payload) => {
    return {...state, tags: payload.tags};
  })
  .case(indexActionCreators.assetSelect, (state, asset) => {
    const index = findAssetIndexById(state.assets, asset.id);
    return {...state, selectedAsset: {...asset, index}};
  })
  .case(boundingBoxActionCreators.assign, (state, payload) => {
    // FIXME: O(n)
    const index = findAssetIndexById(state.assets, payload.asset.id);
    return {...state, ...updateAssets(state, {...payload.asset, index})};
  })
  .case(boundingBoxActionCreators.unAssign, (state, payload) => {
    // FIXME: O(n)
    const index = findAssetIndexById(state.assets, payload.asset.id);
    return {...state, ...updateAssets(state, {...payload.asset, index})};
  })
  .case(boundingBoxActionCreators.modifyRequest, (state, payload) => {
    // FIXME: O(n)
    const index = findAssetIndexById(state.assets, payload.asset.id);
    if (payload.asset.boundingBoxes == null) {
      return state;
    }
    const newBoxes = replaceBoxById(payload.asset.boundingBoxes, payload.box)
    return {...state, ...updateAssets(state, {...payload.asset, index, boundingBoxes: newBoxes})};
  })
  .case(boundingBoxActionCreators.move, (state, payload) => {
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
    const newBoxes = replaceBoxById(asset.boundingBoxes, newBox);
    return {...state, ...updateAssets(state, {...asset, index, boundingBoxes: newBoxes})};
  })
  .case(boundingBoxActionCreators.scale, (state, payload) => {
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
    const newBoxes = replaceBoxById(asset.boundingBoxes, newBox);
    return {...state, ...updateAssets(state, {...asset, index, boundingBoxes: newBoxes})};
  })
  .case(indexActionCreators.downArrowKey, (state, payload) => {
    if (!state.selectedAsset) {
      return {...state};
    }
    const index = findAssetIndexById(state.assets, state.selectedAsset.id);
    switch (payload) {
      case 'UP':
        if (index === 0) {
          return {...state};
        }
        return {...state, selectedAsset: {...state.assets[index - 1], index: index - 1}};
      case 'DOWN':
        if (index === state.assets.length - 1) {
          return {...state};
        }
        return {...state, selectedAsset: {...state.assets[index + 1], index: index + 1}};
      default:
        return {...state};
    }
  })

const updateAssets = (state: GlobalState, asset: AssetWithIndex) => {
  const assets = replaceBy(state.assets, asset, (a) => a.id === asset.id);
  const selectedAsset = state.selectedAsset?.id === asset.id ? asset : state.selectedAsset;
  return {assets, selectedAsset};
}
