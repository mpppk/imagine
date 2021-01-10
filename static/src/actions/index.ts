import actionCreatorFactory from 'typescript-fsa';
import { Asset, Direction, Query, Tag, WorkSpace } from '../models/models';
import { WSPayload } from './workspace';

const indexActionCreatorFactory = actionCreatorFactory('INDEX');

export interface DragResizeHandlerPayload {
  dx: number;
  dy: number;
}

export interface ClickFilterApplyButtonPayload {
  enabled: boolean;
  changed: boolean;
  queries: Query[];
}

export interface ClickChangeBaseButtonPathPayload extends WSPayload {
  needToLoadAssets: boolean;
}

export const indexActionCreators = {
  selectTag: indexActionCreatorFactory<Tag>('TAG/SELECT'),
  clickChangeBasePathButton: indexActionCreatorFactory<ClickChangeBaseButtonPathPayload>(
    'CHANGE_BASE_PATH_BUTTON/CLICK'
  ),
  clickAddTagButton: indexActionCreatorFactory<void>('ADD_TAG_BUTTON/CLICK'),
  clickEditTagButton: indexActionCreatorFactory<Tag>('EDIT_TAG_BUTTON/CLICK'),
  clickFilterApplyButton: indexActionCreatorFactory<ClickFilterApplyButtonPayload>(
    'FILTER_APPLY_BUTTON/CLICK'
  ),
  clickWorkspaceName: indexActionCreatorFactory<WorkSpace>(
    'WORKSPACE_NAME/CLICK'
  ),
  changeFilterMode: indexActionCreatorFactory<boolean>('FILTER_MODE/CHANGE'),
  downAlphabetKey: indexActionCreatorFactory<string>('ALPHABET_KEY/DOWN'),
  downNumberKey: indexActionCreatorFactory<number>('NUMBER_KEY/DOWN'),
  downArrowKey: indexActionCreatorFactory<Direction>('ARROW_KEY/DOWN'),
  downSymbolKey: indexActionCreatorFactory<number>('SYMBOL_KEY/DOWN'),
  assetSelect: indexActionCreatorFactory<Asset>('ASSET/SELECT'),
  dragResizeHandler: indexActionCreatorFactory<DragResizeHandlerPayload>(
    'RESIZE_HANDLER/DRAG'
  ),
};
