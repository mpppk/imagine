import { SagaIterator } from 'redux-saga';
import { call, delay, takeEvery } from 'redux-saga/effects';
import { ActionCreator } from 'typescript-fsa';
import { bindAsyncAction } from 'typescript-fsa-redux-saga';
import {
  globalActionCreators,
  globalAsyncActionCreators,
} from '../actions/global';
import { requestSignIn } from '../services/session';

export const signInWorker = bindAsyncAction(globalAsyncActionCreators.signIn)(
  function* (payload) {
    return yield call(requestSignIn, payload.email, payload.password);
  }
);

export const signOutWorker = bindAsyncAction(
  globalAsyncActionCreators.signOut,
  {
    skipStartedAction: true,
  }
)(function* (_payload) {
  yield delay(1000);
});

export function* watchClickSignInSubmitButton() {
  yield takeEvery(globalActionCreators.clickSignInSubmitButton, function* (
    action
  ) {
    yield call(signInWorker, action.payload);
  });
}

export type SagaWorker = <T>(params: T, ...args: any[]) => SagaIterator;
export const takeEveryAction = <T>(
  ac: ActionCreator<T>,
  worker: SagaWorker
) => {
  return function* () {
    yield takeEvery(ac, function* (action) {
      yield call(worker, action.payload);
    });
  };
};
