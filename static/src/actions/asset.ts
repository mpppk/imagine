import actionCreatorFactory from 'typescript-fsa';
import {WSPayload} from "./server";
import {Asset} from "../models/models";

const assetActionCreatorFactory = actionCreatorFactory('ASSET');

interface RequestAssetsPayload extends WSPayload {
  requestNum: number
}

interface ScanRunningPayload extends WSPayload{
  assets: Asset[]
}

export const assetActionCreators = {
  requestAssets: assetActionCreatorFactory<RequestAssetsPayload>('REQUEST_ASSETS'),
  scanRunning: assetActionCreatorFactory<ScanRunningPayload>('SCAN/RUNNING'),
  scanFinish: assetActionCreatorFactory<WSPayload>('SCAN/FINISH'),
};

