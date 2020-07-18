import { reducerWithInitialState } from 'typescript-fsa-reducers';
import {globalActionCreators} from '../actions/global';
import { User } from '../models/models';

export const globalInitialState = {
  jwt: null as string | null, // FIXME
  user: null as User | null,
  waitingSignIn: false,
  currentWorkSpace: 'a',
  // currentWorkSpace: null as string | null,
  workspaces: ['a', 'b', 'c'] as string[], // FIXME
  isLoadingWorkSpaces: false
};

export type GlobalState = typeof globalInitialState;
export const global = reducerWithInitialState(globalInitialState)
  .case(globalActionCreators.selectNewWorkSpace, (state,workspace) => {
    return { ...state, currentWorkSpace: workspace };
  })
