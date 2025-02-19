export function ClientSideSort<T>(d: T[], sortColumn: string, sortType: 'desc' | 'asc' | undefined): T[] {
  return [...d].sort((a, b) => {
    let x: any = getColumnValue(a, sortColumn);
    let y: any = getColumnValue(b, sortColumn);

    const dateRegex = /^([0-9]{4}-[0-9]{2}-[0-9]{2}.*)$/;
    if (dateRegex.test(x) && dateRegex.test(y)) {
      x = new Date(x);
      y = new Date(y);
      if (sortType === 'asc') {
        return x - y;
      } else {
        return y - x;
      }
    } else if (typeof x === 'string' && typeof y === 'string') {
      x = x.charCodeAt(0);
      y = y.charCodeAt(0);
      if (sortType === 'asc') {
        return x - y;
      } else {
        return y - x;
      }
    } else {
      if (sortType === 'asc') {
        return x - y;
      } else {
        return y - x;
      }
    }
  });
}

const getColumnValue = (object: any, sortColumn: string): any => {
  let x: any;
  Object.entries(object).forEach((entry) => {
    if (entry[0] === sortColumn) {
      x = entry[1];
    }
  });
  return x;
};
