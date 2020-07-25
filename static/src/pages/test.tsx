import {Button, Typography} from "@material-ui/core";
import {createStyles, makeStyles, Theme} from '@material-ui/core/styles';
import Toolbar from '@material-ui/core/Toolbar';
import React, {useState} from 'react';
import {resetServerContext} from "react-beautiful-dnd";
import {useSelector} from "react-redux";
import {indexActionCreators} from "../actions";
import {ImageListDrawer} from "../components/ImageListDrawer";
import {TagListDrawer} from "../components/TagListDrawer";
import {useActions} from "../hooks";
import {Tag, WorkSpace} from "../models/models";
import {State} from "../reducers/reducer";
import {immutableSplice} from "../util";

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

const useHandlers = (localState: LocalState, setLocalState: (s: LocalState) => void, _globalState: GlobalState) => {
  const actionCreators = useActions(indexActionCreators);
  return {
    addDirectoryButton: (ws: WorkSpace) => {
      actionCreators.clickAddDirectoryButton({workSpaceName: ws.name});
    },
    clickAddTagButton: () => {
      setLocalState({
        ...localState,
        editTagId: localState.maxId,
        maxId: localState.maxId+1,
        tags: [{id: localState.maxId, name: ''}, ...localState.tags],
      });
    },
    clickEditTagButton: (tag: Tag) => {
      setLocalState({...localState, editTagId: tag.id});
    },
    renameTag: (tag: Tag) => {
      const targetTagIndex = localState.tags.findIndex((t) => t.id === tag.id);
      if (targetTagIndex === -1) {
        // tslint:disable-next-line:no-console
        console.warn('unknown tag ID is provided', tag);
      }
      setLocalState({
        ...localState,
        editTagId: null,
        tags: immutableSplice(localState.tags, targetTagIndex, 1, tag),
      });
    },
  };
}

interface GlobalState {
  currentWorkSpace: WorkSpace | null
  imagePaths: string[]
  isLoadingWorkSpace: boolean
  isScanningDirectories: boolean
}

const selector = (state: State): GlobalState => ({
  currentWorkSpace: state.global.currentWorkSpace,
  imagePaths: state.indexPage.imagePaths,
  isLoadingWorkSpace: state.global.isLoadingWorkSpaces,
  isScanningDirectories: state.indexPage.scanning,
})

interface LocalState {
  tags: Tag[]
  maxId: number
  editTagId: number | null
}

// fake data generator
const generateTags = (count: number) =>
  Array.from({length: count}, (_, k) => k).map(k => ({
    id: k,
    name: `item-${k}`,
  } as Tag));

const generateInitialLocalState = (tagNum: number): LocalState => {
  return {
    editTagId: null,
    maxId: tagNum+1,
    tags: generateTags(tagNum),
  }
};

export default function Test() {
  const classes = useStyles();
  const [localState, setLocalState] = useState(generateInitialLocalState(4));
  const globalState = useSelector(selector)
  const handlers = useHandlers(localState, setLocalState, globalState);
  const handleClickAddDirectoryButton = () => {
    if (globalState.currentWorkSpace === null) {
      // tslint:disable-next-line:no-console
      console.warn('workspace is not selected, but AddDirectoryButton is clicked')
      return
    }
    handlers.addDirectoryButton(globalState.currentWorkSpace)
  }

  return (
    <div className={classes.root}>
      <ImageListDrawer imagePaths={globalState.imagePaths}/>
      <main className={classes.content}>
        <Toolbar/>
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
        <Typography paragraph={true}>
          Consequat mauris nunc congue nisi vitae suscipit. Fringilla est ullamcorper eget nulla
          facilisi etiam dignissim diam. Pulvinar elementum integer enim neque volutpat ac
          tincidunt. Ornare suspendisse sed nisi lacus sed viverra tellus. Purus sit amet volutpat
          consequat mauris. Elementum eu facilisis sed odio morbi. Euismod lacinia at quis risus sed
          vulputate odio. Morbi tincidunt ornare massa eget egestas purus viverra accumsan in. In
          hendrerit gravida rutrum quisque non tellus orci ac. Pellentesque nec nam aliquam sem et
          tortor. Habitant morbi tristique senectus et. Adipiscing elit duis tristique sollicitudin
          nibh sit. Ornare aenean euismod elementum nisi quis eleifend. Commodo viverra maecenas
          accumsan lacus vel facilisis. Nulla posuere sollicitudin aliquam ultrices sagittis orci a.
        </Typography>
        <Button variant="outlined" color="primary"
                disabled={globalState.isScanningDirectories || globalState.isLoadingWorkSpace}
                onClick={handleClickAddDirectoryButton}>
          {globalState.isScanningDirectories ? 'Scanning...' : 'Add Directory'}
        </Button>
      </main>
      <TagListDrawer
        tags={localState.tags}
        editTagId={localState.editTagId ?? undefined}
        onClickAddButton={handlers.clickAddTagButton}
        onClickEditButton={handlers.clickEditTagButton}
        onRename={handlers.renameTag}
      />
    </div>
  );
}

export async function getServerSideProps() {
  resetServerContext()
  return {props: {}}
}

