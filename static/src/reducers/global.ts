import { reducerWithInitialState } from 'typescript-fsa-reducers';
import {globalActionCreators} from '../actions/global';
import {serverActionCreators} from "../actions/server";
import {User, WorkSpace} from '../models/models';

export const globalInitialState = {
  currentWorkSpace: null as WorkSpace | null,
  isLoadingWorkSpaces: true,
  jwt: null as string | null, // FIXME
  user: null as User | null,
  waitingSignIn: false,
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
