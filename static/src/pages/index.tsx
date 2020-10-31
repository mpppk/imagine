import {Button, Typography} from "@material-ui/core";
import {createStyles, makeStyles, Theme} from '@material-ui/core/styles';
import Toolbar from '@material-ui/core/Toolbar';
import React, {useState} from 'react';
import {resetServerContext} from "react-beautiful-dnd";
import {useSelector} from "react-redux";
import {indexActionCreators} from "../actions";
import {ImageListDrawer} from "../components/ImageListDrawer";
import {ImagePreview} from "../components/ImagePreview";
import {TagListDrawer} from "../components/TagListDrawer";
import {useActions, useVirtualizedAsset} from "../hooks";
import {Asset, AssetWithIndex, Tag, WorkSpace} from "../models/models";
import {State} from "../reducers/reducer";
import {assetPathToUrl, findAssetIndexById, isArrowKeyCode, keyCodeToDirection} from "../util";
import {tagActionCreators} from "../actions/tag";
import uniq from "lodash/uniq";
import _ from "lodash";
import {boundingBoxActionCreators} from "../actions/box";
import {Pixel} from "../components/svg/svg";

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    content: {
      flexGrow: 1,
      padding: theme.spacing(3),
    },
    root: {
      display: 'flex',
    },
  }),
);

const useHandlers = (localState: LocalState, setLocalState: (s: LocalState) => void, globalState: GlobalState) => {
  const indexActionDispatcher = useActions(indexActionCreators);
  const tagActionDispatcher = useActions(tagActionCreators);
  const boxActionDispatcher = useActions(boundingBoxActionCreators);

  return {
    ...indexActionDispatcher,
    ...tagActionDispatcher,
    clickAddDirectoryButton: (ws: WorkSpace) => {
      indexActionDispatcher.clickAddDirectoryButton({workSpaceName: ws.name});
    },
    clickTag: (tag: Tag) => {
      indexActionDispatcher.selectTag(tag);
    },
    clickAddTagButton: () => {
      const tag: Tag = {id: globalState.tags.length + 1, name: ''};
      indexActionDispatcher.clickAddTagButton(tag);
      setLocalState({...localState, editTagId: tag.id});
    },
    renameTag: (tag: Tag) => {
      const workSpaceName = globalState.currentWorkSpace?.name!;
      if (workSpaceName === undefined) {
        // tslint:disable-next-line:no-console
        console.warn('workspace is null but tag is renamed', tag);
      }
      tagActionDispatcher.rename({workSpaceName, tag});
      setLocalState({...localState, editTagId: null});
    },
    clickEditTagButton: (tag: Tag) => {
      setLocalState({...localState, editTagId: tag.id});
    },
    clickImage: (__: string, index: number) => {
      indexActionDispatcher.assetSelect(globalState.assets[index])
    },
    updateTags: (tags: Tag[]) => {
      const workSpaceName = globalState.currentWorkSpace?.name!;
      if (workSpaceName === undefined) {
        // tslint:disable-next-line:no-console
        console.warn('workspace is null but tags are updated', tags);
      }
      const newTags: Tag[] = tags.map((t, i) => ({...t, index: i}))
      tagActionDispatcher.update({workSpaceName, tags: newTags});
    },
    keyDown: (e: any) => {
      e.preventDefault();
      if (e.keyCode >= 48 && e.keyCode <= 57) {
        indexActionDispatcher.downNumberKey(e.keyCode - 48);
        return;
      }

      if (isArrowKeyCode(e.keyCode)) {
        indexActionDispatcher.downArrowKey(keyCodeToDirection(e.keyCode));
      }
    },

    onMoveBoundingBox: _.debounce((boxID: number, dx: Pixel, dy: Pixel) => {
      if (globalState.selectedAsset === null || globalState.currentWorkSpace === null) {
        return;
      }
      boxActionDispatcher.move({
        workSpaceName: globalState.currentWorkSpace.name,
        assetID: globalState.selectedAsset.id,
        boxID,
        dx, dy,
      })
    }, 50, {maxWait: 150}),
    onScaleBoundingBox: _.debounce((boxID: number, dx: Pixel, dy: Pixel) => {
      if (globalState.selectedAsset === null || globalState.currentWorkSpace === null) {
        return;
      }
      boxActionDispatcher.scale({
        workSpaceName: globalState.currentWorkSpace.name,
        assetID: globalState.selectedAsset.id,
        boxID,
        dx, dy,
      })
    }, 50, {maxWait: 150}),
    onDeleteBoundingBox: (boxID:number) => {
      if (globalState.selectedAsset === null || globalState.currentWorkSpace === null) {
        return;
      }
      if (globalState.currentWorkSpace === undefined) {
        // tslint:disable-next-line:no-console
        console.warn('workspace is null but bounding box is deleted. id:', boxID);
      }
      boxActionDispatcher.deleteRequest({
        assetID: globalState.selectedAsset.id,
        boxID,
        workSpaceName: globalState.currentWorkSpace?.name
      });
    },
  };
}

interface GlobalState {
  assets: Asset[]
  assignedTagIds: number[]
  selectedTagId?: number
  selectedAsset: AssetWithIndex | null
  selectedAssetUrl?: string
  currentWorkSpace: WorkSpace | null
  tags: Tag[]
  imagePaths: string[]
  isLoadingWorkSpace: boolean
  isScanningDirectories: boolean
  selectedAssetIndex: number
  imageDrawerHeight: number
}

const selector = (state: State): GlobalState => {
  const boxes = state.global.selectedAsset?.boundingBoxes ?? [];
  const assignedTagIds = boxes.map((box) => box.tag.id);
  const selectedAssetIndex = state.global.selectedAsset ?
    findAssetIndexById(state.global.assets, state.global.selectedAsset.id) :
    -1;
  return {
    assets: state.global.assets,
    assignedTagIds: uniq(assignedTagIds),
    selectedAsset: state.global.selectedAsset,
    selectedAssetUrl: state.global.selectedAsset === null ?
      undefined : assetPathToUrl(state.global.selectedAsset?.path),
    currentWorkSpace: state.global.currentWorkSpace,
    tags: state.global.tags,
    imagePaths: state.global.assets.map((a) => assetPathToUrl(a.path)),
    isLoadingWorkSpace: state.global.isLoadingWorkSpaces,
    isScanningDirectories: state.indexPage.scanning,
    selectedTagId: state.global.selectedTagId,
    selectedAssetIndex,
    imageDrawerHeight: state.global.windowHeight,
  };
};

interface LocalState {
  editTagId: number | null
}

const generateInitialLocalState = (): LocalState => {
  return {
    editTagId: null,
  }
};

export default function Test() {
  const classes = useStyles();
  const [localState, setLocalState] = useState(generateInitialLocalState());
  const globalState = useSelector(selector)
  const handlers = useHandlers(localState, setLocalState, globalState);
  const virtualizedAssetProps = useVirtualizedAsset();
  const handleClickAddDirectoryButton = () => {
    if (globalState.currentWorkSpace === null) {
      // tslint:disable-next-line:no-console
      console.warn('workspace is not selected, but AddDirectoryButton is clicked')
      return
    }
    handlers.clickAddDirectoryButton(globalState.currentWorkSpace)
  }

  return (
    <div className={classes.root} onKeyDown={handlers.keyDown} tabIndex={0}>
      <ImageListDrawer
        {...virtualizedAssetProps}
        imagePaths={globalState.imagePaths}
        onClickImage={handlers.clickImage}
        selectedIndex={globalState.selectedAssetIndex}
        height={globalState.imageDrawerHeight}
      />
      <main className={classes.content}>
        <Toolbar/>
        {globalState.selectedAssetUrl === undefined || globalState.selectedAsset === null ? null : <ImagePreview
          src={globalState.selectedAssetUrl}
          asset={globalState.selectedAsset}
          onMoveBoundingBox={handlers.onMoveBoundingBox}
          onScaleBoundingBox={handlers.onScaleBoundingBox}
          onDeleteBoundingBox={handlers.onDeleteBoundingBox}
        />}
        <Typography paragraph={true}>
          Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt
          ut labore et dolore magna aliqua. Rhoncus dolor purus non enim praesent elementum
          facilisis leo vel. Risus at ultrices mi tempus imperdiet. Semper risus in hendrerit
          gravida rutrum quisque non tellus. Convallis convallis tellus id interdum velit laoreet id
          donec ultrices. Odio morbi quis commodo odio aenean sed adipiscing. Amet nisl suscipit
          adipiscing bibendum est ultricies integer quis. Cursus euismod quis viverra nibh cras.
          Metus vulputate eu scelerisque felis imperdiet proin fermentum leo. Mauris commodo quis
          imperdiet massa tincidunt. Cras tincidunt lobortis feugiat vivamus at augue. At augue eget
          arcu dictum varius duis at consectetur lorem. Velit sed ullamcorper morbi tincidunt. Lorem
          donec massa sapien faucibus et molestie ac.
        </Typography>
        <Button variant="outlined" color="primary"
                disabled={globalState.isScanningDirectories || globalState.isLoadingWorkSpace}
                onClick={handleClickAddDirectoryButton}>
          {globalState.isScanningDirectories ? 'Scanning...' : 'Add Directory'}
        </Button>
      </main>
      <TagListDrawer
        tags={globalState.tags}
        editTagId={localState.editTagId ?? undefined}
        selectedTagId={globalState.selectedTagId}
        assignedTagIds={globalState.assignedTagIds}
        onClickItem={handlers.clickTag}
        onClickAddButton={handlers.clickAddTagButton}
        onClickEditButton={handlers.clickEditTagButton}
        onRename={handlers.renameTag}
        onUpdate={handlers.updateTags}
      />
    </div>
  );
}

export async function getServerSideProps() {
  resetServerContext()
  return {props: {}}
}
