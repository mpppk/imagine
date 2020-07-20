import actionCreatorFactory from 'typescript-fsa';
import {WorkSpace} from '../models/models';

const globalActionCreatorFactory = actionCreatorFactory('GLOBAL');

export const globalActionCreators = {
  requestWorkSpaces: globalActionCreatorFactory<void>('REQUEST_WORKSPACES'),
  selectNewWorkSpace: globalActionCreatorFactory<WorkSpace>('SELECT_NEW_WORKSPACE'),
};

