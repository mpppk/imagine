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
  const handleChangeTags = (newTags: Tag[]) => {
    setTags(newTags)
  }

  const handleClickAddItemButton = (newTags: Tag[]) => {
    setTags(newTags)
  }

  return (
    <TagList tags={tags} onChange={handleChangeTags} onClickAddButton={handleClickAddItemButton}/>
  );
};



export async function getServerSideProps() {
  resetServerContext()
  return {props: {}}
}

export default Preferences;
