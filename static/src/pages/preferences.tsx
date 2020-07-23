import {NextPage} from 'next';
import React, {useState} from 'react';
import {resetServerContext} from "react-beautiful-dnd";
import {TagList} from "../components/TagList";
import {Tag} from "../models/models";

// fake data generator
const generateTags = (count: number) =>
  Array.from({length: count}, (_, k) => k).map(k => ({
    id: k,
    name: `item-${k}`,
  } as Tag));

// tslint:disable-next-line variable-name
export const Preferences: NextPage = () => {
  const [tags, setTags] = useState(generateTags(5));
  const [maxId, setMaxId] = useState(6);
  const [editTagId, setEditTagId] = useState(null as number | null);
  const handleChangeTags = (newTags: Tag[]) => {
    setTags(newTags)
  }

  const handleClickAddItemButton = () => {
    setTags([...tags, {id: maxId, name: ''}])
    setEditTagId(maxId);
    setMaxId(maxId + 1);
  }

  const handleClickEditButton = (tag: Tag) => {
    setEditTagId(tag.id);
  }

  const handleRename = ()  => {
    setEditTagId(null)
  }

  return (
    <TagList
      tags={tags}
      editTagId={editTagId ?? undefined}
      onUpdate={handleChangeTags}
      onClickAddButton={handleClickAddItemButton}
      onClickEditButton={handleClickEditButton}
      onRename={handleRename}
    />
  );
};

export async function getServerSideProps() {
  resetServerContext()
  return {props: {}}
}

export default Preferences;
