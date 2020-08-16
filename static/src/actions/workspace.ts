import actionCreatorFactory from 'typescript-fsa';
import {WorkSpace} from "../models/models";

const workspaceActionCreatorFactory = actionCreatorFactory('WORKSPACE');

export interface WSPayload {
  workSpaceName: string
}

export const workspaceActionCreators = {
  scanRequest: workspaceActionCreatorFactory<void>('SCAN/REQUEST'),
  select: workspaceActionCreatorFactory<WorkSpace>('SELECT'),
  scanResult: workspaceActionCreatorFactory<WorkSpace[]>('SCAN/RESULT'),
};

