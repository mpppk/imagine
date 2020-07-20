import { reducerWithInitialState } from 'typescript-fsa-reducers';
import {globalActionCreators} from '../actions/global';
import {serverActionCreators} from "../actions/server";
import {WorkSpace} from '../models/models';

export const globalInitialState = {
  currentWorkSpace: null as WorkSpace | null,
  isLoadingWorkSpaces: true,
  workspaces: null as WorkSpace[] | null,
};

export type GlobalState = typeof globalInitialState;
export const global = reducerWithInitialState(globalInitialState)
  .case(serverActionCreators.scanWorkSpaces, (state,workspaces) => {
    const currentWorkSpace = workspaces.length > 0 ? workspaces[0] : null;
    return { ...state, workspaces, isLoadingWorkSpaces: false, currentWorkSpace};
  })
  .case(globalActionCreators.selectNewWorkSpace, (state,workspace) => {
    return { ...state, currentWorkSpace: workspace };
  })
