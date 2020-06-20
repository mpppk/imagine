import actionCreatorFactory  from 'typescript-fsa';

export interface IRequestAmountChangingPayload {
  amount: number;
}

const counterActionCreatorFactory = actionCreatorFactory('COUNTER');

export const counterActionCreators = {
  clickAsyncIncrementButton: counterActionCreatorFactory<void>(
    'CLICK_ASYNC_INCREMENT_BUTTON'
  ),
  clickDecrementButton: counterActionCreatorFactory<void>(
    'CLICK_DECREMENT_BUTTON'
  ),
  clickIncrementButton: counterActionCreatorFactory<void>(
    'CLICK_INCREMENT_BUTTON'
  ),
  requestAmountChanging: counterActionCreatorFactory<
    IRequestAmountChangingPayload
  >('REQUEST_AMOUNT_CHANGING')
};

export interface IRequestAmountChangingWithSleepPayload
  extends IRequestAmountChangingPayload {
  sleep: number;
}

export const counterAsyncActionCreators = {
  changeAmountWithSleep: counterActionCreatorFactory.async<
    IRequestAmountChangingWithSleepPayload,
    IRequestAmountChangingPayload,
    any
  >('REQUEST_AMOUNT_CHANGING_WITH_SLEEP')
};
