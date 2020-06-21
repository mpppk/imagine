import actionCreatorFactory from 'typescript-fsa';

const wsActionCreatorFactory = actionCreatorFactory('WS');

export const wsActionCreators = {
  open: wsActionCreatorFactory<Event>('OPEN'),
  close: wsActionCreatorFactory<void>('CLOSE'),
  error: wsActionCreatorFactory<Event>('ERROR'),
  message: wsActionCreatorFactory<MessageEvent>('MESSAGE'),
};

