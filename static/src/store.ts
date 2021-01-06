import { Context, createWrapper, MakeStore } from 'next-redux-wrapper';
import { applyMiddleware, createStore, Middleware } from 'redux';
import createSagaMiddleware from 'redux-saga';
import { makeLorcaMiddleware, setupServerActionHandler } from './lib';
import { initialState, reducer, State } from './reducers/reducer';
import rootSaga from './sagas/saga';
import { makeMockLorcaMiddleware } from './mock_backend_middleware';

const sagaMiddleware = createSagaMiddleware();

const bindMiddleware = (middlewareList: Middleware[]) => {
  if (process.env.NODE_ENV !== 'production') {
    const { composeWithDevTools } = require('redux-devtools-extension');
    return composeWithDevTools(applyMiddleware(...middlewareList));
  }
  return applyMiddleware(...middlewareList);
};

const makeStore: MakeStore<State> = (_context: Context) => {
  const middlewares: Middleware[] = [sagaMiddleware];
  if (useMockBackEnd()) {
    middlewares.push(makeMockLorcaMiddleware());
  } else {
    middlewares.push(makeLorcaMiddleware());
  }

  const store = createStore(reducer, initialState, bindMiddleware(middlewares));

  (store as any).runSagaTask = () => {
    (store as any).sagaTask = sagaMiddleware.run(rootSaga); // FIXME Add type
  };

  (store as any).runSagaTask(); // FIXME Add type
  setupServerActionHandler(store);
  return store;
};

const isEnableDebugMode = (): boolean => {
  return (process.env.enableReduxWrapperDebugMode as any) as boolean;
};

const useMockBackEnd = (): boolean => {
  return (process.env.NEXT_PUBLIC_USE_MOCK_BACKEND as any) as boolean;
};

export const wrapper = createWrapper<State>(makeStore, {
  debug: isEnableDebugMode(),
});
