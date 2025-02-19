import React from 'react';
import { Cell, Legend, Pie, PieChart, ResponsiveContainer } from 'recharts';

export type PieChartData = {
  name: string;
  value: number;
  fill: string;
  filterFunc?: (filter: string) => void;
  filterString?: string;
};

const LegendWithValues = ({ payload }: any) => (
  <ul style={{"marginBottom": "10px"}}>
    {payload.map((p: any) => {
      return (
        <li style={{ color: p.color }} key={p.value}>
          <strong className='legendFilterLink' onClick={() => p.payload.filterFunc(p.payload.filterString)}>▪ {p.value}</strong>: {p.payload.value}
        </li>
      );
    })}
  </ul>
);

const LegendWithZeroValue = ({ payload }: any) => (
  <ul style={{"marginBottom": "10px"}}>
    {payload.map((p: any) => {
      return (
        <li style={{ color: p.color }} key={p.value}>
          none
        </li>
      );
    })}
  </ul>
);

export const PieData = ({data} : {data: PieChartData[]}) => {
  return (
    <ResponsiveContainer width="100%" height="100%">
      <PieChart>
        <Pie data={data} cx="50%" cy="50%" innerRadius="88%" outerRadius="95%" fill="#8884d8" dataKey="value">
          {data.map((entry: any, index: any) => (
            <Cell key={`cell-${index}`} fill={entry.fill} onClick={() => entry.filterFunc(entry.filterString)} className='legendFilterLink'/>
          ))}
          {data.length === 0 && <Cell fill="gray" />}
        </Pie>
        {(data.length === 1 && data[0].name === 'nil') && <Legend iconSize={5} layout="vertical" verticalAlign="bottom" content={LegendWithZeroValue} />}
        {(data.length >= 1) && data[0].name !== 'nil' && <Legend iconSize={5} layout="vertical" verticalAlign="bottom" content={LegendWithValues} />}
      </PieChart>
    </ResponsiveContainer>
  );
};
