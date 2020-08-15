import actionCreatorFactory from 'typescript-fsa';
import {Asset, Tag, WorkSpace} from "../models/models";

const serverActionCreatorFactory = actionCreatorFactory('SERVER');

export interface WSPayload {
  workSpaceName: string
}

interface ScanningImagesPayload extends WSPayload{
  paths: string[]
}

interface ScanningAssetsPayload extends WSPayload{
  assets: Asset[]
}

interface TagScanPayload extends WSPayload {
  tags: Tag[]
}

export const serverActionCreators = {
  cancelDirectoryScanning: serverActionCreatorFactory<WSPayload>(
    'CANCEL_DIRECTORY_SCANNING'
  ),
  finishAssetsScanning: serverActionCreatorFactory<WSPayload>(
    'FINISH_ASSETS_SCANNING'
  ),
  finishDirectoryScanning: serverActionCreatorFactory<WSPayload>(
    'FINISH_DIRECTORY_SCANNING'
  ),
  scanWorkSpaces: serverActionCreatorFactory<WorkSpace[]>(
    'SCAN_WORKSPACES'
  ),
  scanningAssets: serverActionCreatorFactory<ScanningAssetsPayload>(
    'SCANNING_ASSETS'
  ),
  scanningImages: serverActionCreatorFactory<ScanningImagesPayload>(
    'SCANNING_IMAGES'
  ),
  startDirectoryScanning: serverActionCreatorFactory<WSPayload>(
    'START_DIRECTORY_SCANNING'
  ),
  tagScan: serverActionCreatorFactory<TagScanPayload>('TAG/SCAN')
};
