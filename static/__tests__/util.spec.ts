import {immutableSplice} from "../src/util";

describe('immutableSplice', () => {
  it('can delete elements', () => {
    const months = ['Jan', 'Feb', 'March', 'April'];
    expect(immutableSplice(months, 0, 1)).toEqual( ['Feb', 'March', 'April']);
    expect(immutableSplice(months, months.length-1, 1)).toEqual( ['Jan', 'Feb', 'March']);
  })

  it('can add elements', () => {
    const months = ['Feb', 'March'];
    expect(immutableSplice(['Feb', 'March'], 0, 0, 'Jan')).toEqual( ['Jan' ,'Feb', 'March']);
    expect(immutableSplice(['Jan', 'Feb'], months.length, 0, 'March')).toEqual( ['Jan', 'Feb', 'March']);
  })

  // See https://developer.mozilla.org/ja/docs/Web/JavaScript/Reference/Global_Objects/Array/splice
  it('behave like Array.splice', () => {
    const months = ['Jan', 'March', 'April', 'June'];
    const newMonths = immutableSplice(months, 1, 0, 'Feb');
    expect(newMonths).toEqual( ["Jan", "Feb", "March", "April", "June"]);
    const newNewMonths = immutableSplice(newMonths, 4, 1, 'May');
    expect(newNewMonths).toEqual(  ["Jan", "Feb", "March", "April", "May"]);
    expect(months).toEqual( ['Jan', 'March', 'April', 'June'] );
  })
})