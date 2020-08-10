import {reducerWithInitialState} from 'typescript-fsa-reducers';
import {globalActionCreators} from '../actions/global';
import {serverActionCreators} from "../actions/server";
import {Asset, WorkSpace} from '../models/models';
import {assetActionCreators} from "../actions/asset";

export const globalInitialState = {
  assets: [] as Asset[],
  currentWorkSpace: null as WorkSpace | null,
  hasMoreAssets :true,
  isLoadingWorkSpaces: true,
  isScanningAssets: false,
  workspaces: null as WorkSpace[] | null,
};

export type GlobalState = typeof globalInitialState;
export const global = reducerWithInitialState(globalInitialState)
  .case(serverActionCreators.scanWorkSpaces, (state, workspaces) => {
    const currentWorkSpace = workspaces.length > 0 ? workspaces[0] : null;
    return {...state, workspaces, isLoadingWorkSpaces: false, currentWorkSpace};
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
