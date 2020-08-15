import {reducerWithInitialState} from 'typescript-fsa-reducers';
import {assetActionCreators} from "../actions/asset";
import {globalActionCreators} from '../actions/global';
import {serverActionCreators} from "../actions/server";
import {Asset, Tag, WorkSpace} from '../models/models';
import {indexActionCreators} from "../actions";
import {immutableSplice} from "../util";

// fake data generator
const generateTags = (count: number) =>
  Array.from({length: count}, (_, k) => k).map(k => ({
    id: k,
    name: `item-${k}`,
  } as Tag));

export const globalInitialState = {
  assets: [] as Asset[],
  tags: generateTags(5) as Tag[], // FIXME
  // maxTagId: 10,
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
  .case(serverActionCreators.scanningAssets, (state, payload) => {
    return {...state, isScanningAssets: false, assets: [...state.assets, ...payload.assets]}
  })
  .case(serverActionCreators.finishAssetsScanning, (state) => {
    return {...state, isScanningAssets: false, hasMoreAssets: false}
  })
  .case(globalActionCreators.selectNewWorkSpace, (state, workspace) => {
    return {...state, currentWorkSpace: workspace};
  })
  .case(globalActionCreators.selectNewWorkSpace, (state, workspace) => {
    return {...state, currentWorkSpace: workspace};
  })
  .case(indexActionCreators.clickAddTagButton, (state, tag) => {
    return {...state, tags: [tag, ...state.tags]};
  })
  .case(indexActionCreators.renameTag, (state, tag) => {
    const targetTagIndex = state.tags.findIndex((t) => t.id === tag.id);
    if (targetTagIndex === -1) {
      // tslint:disable-next-line:no-console
      console.warn('unknown tag ID is provided', tag);
    }
    return {...state, tags: immutableSplice(state.tags, targetTagIndex, 1, tag)};
  })
  .case(indexActionCreators.updateTags, (state, tags) => {
    return {...state, tags};
  })
