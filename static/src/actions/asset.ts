import actionCreatorFactory from 'typescript-fsa';
import { Asset, Query } from '../models/models';
import { WSPayload } from './workspace';

const assetActionCreatorFactory = actionCreatorFactory('ASSET');

export interface AssetScanRequestPayload extends WSPayload {
  requestNum: number;
  queries: Query[];
  reset: boolean;
}

export interface AssetScanRunningPayload extends WSPayload {
  assets: Asset[];
  count: number;
}

export interface AssetScanResultPayload extends WSPayload {
  count: number;
}

export const assetActionCreators = {
  scanRequest: assetActionCreatorFactory<AssetScanRequestPayload>(
    'SCAN/REQUEST'
  ),
  scanRunning: assetActionCreatorFactory<AssetScanRunningPayload>(
    'SCAN/RUNNING'
  ),
  scanFinish: assetActionCreatorFactory<AssetScanResultPayload>('SCAN/FINISH'),
};
