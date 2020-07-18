import actionCreatorFactory from 'typescript-fsa';
import {WorkSpace} from "../models/models";

const serverActionCreatorFactory = actionCreatorFactory('SERVER');

export interface WSPayload {
  workSpaceName: string
}


interface ScanningImagesPayload extends WSPayload{
  paths: string[]
}

export const serverActionCreators = {
  cancelDirectoryScanning: serverActionCreatorFactory<WSPayload>(
    'CANCEL_DIRECTORY_SCANNING'
  ),
  finishDirectoryScanning: serverActionCreatorFactory<WSPayload>(
    'FINISH_DIRECTORY_SCANNING'
  ),
  scanWorkSpaces: serverActionCreatorFactory<WorkSpace[]>(
    'SCAN_WORKSPACES'
  ),
  scanningImages: serverActionCreatorFactory<ScanningImagesPayload>(
    'SCANNING_IMAGES'
  ),
  startDirectoryScanning: serverActionCreatorFactory<WSPayload>(
    'START_DIRECTORY_SCANNING'
  ),
};
