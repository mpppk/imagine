import { AnyAction, Middleware, Store } from 'redux';

export const makeLorcaMiddleware = (): Middleware => (_store) => (next) => (
  action
) => {
  // @ts-ignore
  if (process.browser && (window as any).dispatchToServer === undefined) {
    // tslint:disable-next-line:no-console
    console.error(
      'action is not dispatched to server because dispatchToServer function is undefined. action:',
      action
    );
  }
  // @ts-ignore
  dispatchToServer(action);
  next(action);
};

const makeServerActionHandler = (store: Store) => {
  return (action: AnyAction) => {
    store.dispatch(action);
  };
};

export const setupServerActionHandler = (store: Store) => {
  if (process.browser) {
    (window as any).handleServerAction = makeServerActionHandler(store);
  }
};
