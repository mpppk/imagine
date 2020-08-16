import {reducerWithInitialState} from 'typescript-fsa-reducers';
import {assetActionCreators} from "../actions/asset";
import {workspaceActionCreators} from '../actions/workspace';
import {serverActionCreators} from "../actions/server";
import {Asset, Tag, WorkSpace} from '../models/models';
import {indexActionCreators} from "../actions";
import {immutableSplice} from "../util";

export const globalInitialState = {
  assets: [] as Asset[],
  tags: [] as Tag[],
  currentWorkSpace: null as WorkSpace | null,
  hasMoreAssets :true,
  isLoadingWorkSpaces: true,
  isScanningAssets: false,
  workspaces: null as WorkSpace[] | null,
};

export type GlobalState = typeof globalInitialState;
export const global = reducerWithInitialState(globalInitialState)
  .case(serverActionCreators.scanWorkSpaces, (state, workspaces) => {
    return {...state, workspaces, isLoadingWorkSpaces: false};
  })
  .case(assetActionCreators.requestAssets, (state) => {
    return {...state, isScanningAssets: true}
  })
  .case(assetActionCreators.scanRunning, (state, payload) => {
    return {...state, isScanningAssets: false, assets: [...state.assets, ...payload.assets]}
  })
  .case(assetActionCreators.scanFinish, (state) => {
    return {...state, isScanningAssets: false, hasMoreAssets: false}
  })
  .case(workspaceActionCreators.selectNewWorkSpace, (state, workspace) => {
    return {...state, currentWorkSpace: workspace};
  })
  .case(indexActionCreators.clickAddTagButton, (state, tag) => {
    return {...state, tags: [tag, ...state.tags]};
  })
  .case(serverActionCreators.tagScan, (state, payload) => {
    return {...state, tags: payload.tags};
  })
  .case(indexActionCreators.renameTag, (state, payload) => {
    const targetTagIndex = state.tags.findIndex((t) => t.id === payload.tag.id);
    if (targetTagIndex === -1) {
      // tslint:disable-next-line:no-console
      console.warn('unknown tag ID is provided', payload.tag);
    }
    return {...state, tags: immutableSplice(state.tags, targetTagIndex, 1, payload.tag)};
  })
  .case(indexActionCreators.updateTags, (state, payload) => {
    return {...state, tags: payload.tags};
  })
