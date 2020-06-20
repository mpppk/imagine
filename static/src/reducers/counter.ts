import { reducerWithInitialState } from 'typescript-fsa-reducers';
import {
  counterActionCreators,
  counterAsyncActionCreators,
} from '../actions/counter';

export const counterInitialState = {
  count: 0,
};

export type CounterState = typeof counterInitialState;

const addCount = (state: CounterState, amount: number) => {
  return { ...state, count: state.count + amount };
};

export const counter = reducerWithInitialState(counterInitialState)
  .case(counterActionCreators.clickIncrementButton, (state) => {
    return addCount(state, 1);
  })
  .case(counterActionCreators.clickDecrementButton, (state) => {
    return addCount(state, -1);
  })
  .case(
    counterAsyncActionCreators.changeAmountWithSleep.done,
    (state, payload) => {
      return addCount(state, payload.result.amount);
    }
  );
