import { delay, takeEvery } from 'redux-saga/effects';
import { bindAsyncAction } from 'typescript-fsa-redux-saga';
import {
  counterActionCreators,
  counterAsyncActionCreators,
  IRequestAmountChangingWithSleepPayload,
} from '../actions/counter';

export const counterIncrementWorker = bindAsyncAction(
  counterAsyncActionCreators.changeAmountWithSleep
)(function* (payload: IRequestAmountChangingWithSleepPayload) {
  yield delay(payload.sleep);
  return { amount: payload.amount };
});

export const counterIncrementWorkerWrapper = () =>
  counterIncrementWorker({ amount: 1, sleep: 1000 });

export function* watchIncrementAsync() {
  yield takeEvery(
    counterActionCreators.clickAsyncIncrementButton.type,
    counterIncrementWorkerWrapper
  );
}
