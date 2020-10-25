import actionCreatorFactory from 'typescript-fsa';
import {Asset, Direction, Tag} from "../models/models";
import {WSPayload} from "./workspace";

const indexActionCreatorFactory = actionCreatorFactory('INDEX');

export interface DragResizeHandlerPayload {
  dx: number;
  dy: number;
}

export const indexActionCreators = {
  selectTag: indexActionCreatorFactory<Tag>('TAG/SELECT'),
  clickAddDirectoryButton: indexActionCreatorFactory<WSPayload>('ADD_DIRECTORY_BUTTON/CLICK'),
  clickAddTagButton: indexActionCreatorFactory<Tag>('ADD_TAG_BUTTON/CLICK'),
  clickEditTagButton: indexActionCreatorFactory<Tag>('EDIT_TAG_BUTTON/CLICK'),
  downNumberKey: indexActionCreatorFactory<number>('NUMBER_KEY/DOWN'),
  downArrowKey: indexActionCreatorFactory<Direction>('ARROW_KEY/DOWN'),
  assetSelect: indexActionCreatorFactory<Asset>('ASSET/SELECT'),
  dragResizeHandler: indexActionCreatorFactory<DragResizeHandlerPayload>('RESIZE_HANDLER/DRAG'),
};

