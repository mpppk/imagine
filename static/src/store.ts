import { Context, createWrapper, MakeStore } from 'next-redux-wrapper';
import { applyMiddleware, createStore, Middleware } from 'redux';
import createSagaMiddleware from 'redux-saga';
import { initialState, reducer, State } from './reducers/reducer';
import rootSaga from './sagas/saga';

const sagaMiddleware = createSagaMiddleware();

const bindMiddleware = (middlewareList: Middleware[]) => {
  if (process.env.NODE_ENV !== 'production') {
    const { composeWithDevTools } = require('redux-devtools-extension');
    return composeWithDevTools(applyMiddleware(...middlewareList));
  }
  return applyMiddleware(...middlewareList);
};

const makeStore: MakeStore<State> = (_context: Context) => {
  const store = createStore(
    reducer,
    initialState,
    bindMiddleware([sagaMiddleware])
  );

  (store as any).runSagaTask = () => {
    (store as any).sagaTask = sagaMiddleware.run(rootSaga); // FIXME Add type
  };

  (store as any).runSagaTask(); // FIXME Add type
  return store;
}

const isEnableDebugMode = (): boolean => {
  return process.env.enableReduxWrapperDebugMode as any as boolean;
}

export const wrapper = createWrapper<State>(makeStore, {debug: isEnableDebugMode()})
