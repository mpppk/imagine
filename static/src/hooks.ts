import { useMemo } from 'react'
import { useDispatch } from 'react-redux'
import { ActionCreatorsMapObject, bindActionCreators } from 'redux';

export function useActions<T extends ActionCreatorsMapObject>(actions: T, deps?: any[]) {
  const dispatch = useDispatch()
  return useMemo(
    () => {
      return bindActionCreators(actions, dispatch)
    },
    deps ? [dispatch, ...deps] : [dispatch]
  )
}