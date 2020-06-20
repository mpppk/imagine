import { reducerWithInitialState } from 'typescript-fsa-reducers';
import { globalAsyncActionCreators } from '../actions/global';
import { User } from '../models/models';

export const globalInitialState = {
  jwt: null as string | null, // FIXME
  user: null as User | null,
  waitingSignIn: false,
};

export type GlobalState = typeof globalInitialState;
export const global = reducerWithInitialState(globalInitialState)
  .case(globalAsyncActionCreators.signIn.started, (state) => {
    return { ...state, waitingSignIn: true };
  })
  .case(globalAsyncActionCreators.signIn.done, (state, payload) => {
    return {
      ...state,
      jwt: payload.result.jwt,
      user: payload.result.user,
      waitingSignIn: false,
    };
  })
  .case(globalAsyncActionCreators.signOut.done, (state) => {
    return {
      ...state,
      jwt: null,
      user: null,
      waitingSignIn: false,
    };
  });
